package luasandbox

import (
	"context"
	"fmt"
	"strings"
	"testing"

	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"

	"github.com/sourcegraph/sourcegraph/internal/observation"
)

func TestSandboxHasNoIO(t *testing.T) {
	ctx := context.Background()

	sandbox, err := newService(&observation.TestContext).CreateSandbox(ctx, CreateOptions{})
	if err != nil {
		t.Fatalf("unexpected error creating sandbox: %s", err)
	}
	defer sandbox.Close()

	t.Run("default", func(t *testing.T) {
		script := `
			io.open('service_test.go', 'rb')
		`
		if _, err := sandbox.RunScript(ctx, RunOptions{}, script); err == nil {
			t.Fatalf("expected error running script")
		} else if !strings.Contains(err.Error(), "attempt to index a non-table object(nil) with key 'open'") {
			t.Fatalf("unexpected error running script: %s", err)
		}
	})

	t.Run("module", func(t *testing.T) {
		script := `
			local io = require("io")
			io.open('service_test.go', 'rb')
		`
		if _, err := sandbox.RunScript(ctx, RunOptions{}, script); err == nil {
			t.Fatalf("expected error running script")
		} else if !strings.Contains(err.Error(), "module io not found") {
			t.Fatalf("unexpected error running script: %s", err)
		}
	})
}

func TestRunWithBasicModule(t *testing.T) {
	var stashedValue lua.LValue

	api := map[string]lua.LGFunction{
		"add": func(state *lua.LState) int {
			a := state.CheckNumber(1)
			b := state.CheckNumber(2)
			state.Push(a + b)

			return 1
		},
		"stash": func(state *lua.LState) int {
			stashedValue = state.CheckAny(1)
			return 1
		},
	}

	testModule := func(state *lua.LState) int {
		t := state.NewTable()
		state.SetFuncs(t, api)
		state.Push(t)
		return 1
	}

	ctx := context.Background()

	sandbox, err := newService(&observation.TestContext).CreateSandbox(ctx, CreateOptions{
		Modules: map[string]lua.LGFunction{
			"testmod": testModule,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error creating sandbox: %s", err)
	}
	defer sandbox.Close()

	script := `
		local testmod = require("testmod")
		testmod.stash(testmod.add(3, testmod.add(6, 9)))
		return testmod.add(38, 4)
	`
	retValue, err := sandbox.RunScript(ctx, RunOptions{}, script)
	if err != nil {
		t.Fatalf("unexpected error running script: %s", err)
	}
	if lua.LVAsNumber(retValue) != 42 {
		t.Errorf("unexpected return value. want=%d have=%v", 42, retValue)
	}
	if lua.LVAsNumber(stashedValue) != 18 {
		t.Errorf("unexpected stashed value. want=%d have=%d", 18, stashedValue)
	}
}

func TestPlayground(t *testing.T) {
	type WX struct {
		Wrapped *lua.LFunction
	}

	recognizers := map[string]lua.LValue{}
	api := map[string]lua.LGFunction{
		"create": func(state *lua.LState) int {
			t := state.CheckTable(1)
			var lx *lua.LFunction
			t.ForEach(func(l1, l2 lua.LValue) {
				switch lua.LVAsString(l1) {
				case "patterns":
					(l2.(*lua.LTable)).ForEach(func(l1, l2 lua.LValue) {
						fmt.Printf("Pattern: %s\n", lua.LVAsString(l2))
					})
				case "generate":
					if l2.Type() != lua.LTFunction {
						panic("Bad type")
					}
					lx = l2.(*lua.LFunction)
				default:
					panic("NOPE")
				}
			})
			state.Push(luar.New(state, WX{Wrapped: lx}))
			return 1
		},
		"register": func(state *lua.LState) int {
			state.CheckTable(1).ForEach(func(l1, l2 lua.LValue) {
				recognizers[lua.LVAsString(l1)] = l2
			})
			return 1
		},
	}

	testModule := func(state *lua.LState) int {
		t := state.NewTable()
		state.SetFuncs(t, api)
		state.Push(t)
		return 1
	}

	ctx := context.Background()

	sandbox, err := newService(&observation.TestContext).CreateSandbox(ctx, CreateOptions{
		Modules: map[string]lua.LGFunction{
			"testmod": testModule,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error creating sandbox: %s", err)
	}
	defer sandbox.Close()

	script := `
		local testmod = require("testmod")
		local goext_recognizer = testmod.create {
			patterns = {".go"},
			generate = function(paths, r)
				print('LOL', r)
				coroutine.yield {foo="bonk1"}
				coroutine.yield ({bar="bonk2"}, {baz="bonk3"})
				return {bonk="bonk4"}
			end,
		}

		testmod.register({
			["sg.go"] = goext_recognizer,
			["sg.typescript"] = false,
		})

		return {foo="bar"}
	`
	if _, err := sandbox.RunScript(ctx, RunOptions{}, script); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	for k, v := range recognizers {
		fmt.Printf("> %s\n", k)
		if ux, ok := v.(*lua.LUserData); ok {
			fmt.Printf("Is user data\n")
			if vx, ok := ux.Value.(WX); ok {
				fmt.Printf("Is WX\n")
				values, err := sandbox.CallGenerator(ctx, RunOptions{}, vx.Wrapped, [3]string{"foo", "bar", "baz"}, 250, false)
				if err != nil {
					t.Fatalf("oopsie: %s", err)
				}
				fmt.Printf("> VALUES: %v\n", values)
			}
		}
	}
}
