// Code generated by go-mockgen 1.1.2; DO NOT EDIT.

package dbmock

import (
	"context"
	"sync"

	database "github.com/sourcegraph/sourcegraph/internal/database"
	basestore "github.com/sourcegraph/sourcegraph/internal/database/basestore"
)

// MockConfStore is a mock implementation of the ConfStore interface (from
// the package github.com/sourcegraph/sourcegraph/internal/database) used
// for unit testing.
type MockConfStore struct {
	// DoneFunc is an instance of a mock function object controlling the
	// behavior of the method Done.
	DoneFunc *ConfStoreDoneFunc
	// HandleFunc is an instance of a mock function object controlling the
	// behavior of the method Handle.
	HandleFunc *ConfStoreHandleFunc
	// SiteCreateIfUpToDateFunc is an instance of a mock function object
	// controlling the behavior of the method SiteCreateIfUpToDate.
	SiteCreateIfUpToDateFunc *ConfStoreSiteCreateIfUpToDateFunc
	// SiteGetLatestFunc is an instance of a mock function object
	// controlling the behavior of the method SiteGetLatest.
	SiteGetLatestFunc *ConfStoreSiteGetLatestFunc
	// TransactFunc is an instance of a mock function object controlling the
	// behavior of the method Transact.
	TransactFunc *ConfStoreTransactFunc
}

// NewMockConfStore creates a new mock of the ConfStore interface. All
// methods return zero values for all results, unless overwritten.
func NewMockConfStore() *MockConfStore {
	return &MockConfStore{
		DoneFunc: &ConfStoreDoneFunc{
			defaultHook: func(error) error {
				return nil
			},
		},
		HandleFunc: &ConfStoreHandleFunc{
			defaultHook: func() *basestore.TransactableHandle {
				return nil
			},
		},
		SiteCreateIfUpToDateFunc: &ConfStoreSiteCreateIfUpToDateFunc{
			defaultHook: func(context.Context, *int32, string) (*database.SiteConfig, error) {
				return nil, nil
			},
		},
		SiteGetLatestFunc: &ConfStoreSiteGetLatestFunc{
			defaultHook: func(context.Context) (*database.SiteConfig, error) {
				return nil, nil
			},
		},
		TransactFunc: &ConfStoreTransactFunc{
			defaultHook: func(context.Context) (database.ConfStore, error) {
				return nil, nil
			},
		},
	}
}

// NewStrictMockConfStore creates a new mock of the ConfStore interface. All
// methods panic on invocation, unless overwritten.
func NewStrictMockConfStore() *MockConfStore {
	return &MockConfStore{
		DoneFunc: &ConfStoreDoneFunc{
			defaultHook: func(error) error {
				panic("unexpected invocation of MockConfStore.Done")
			},
		},
		HandleFunc: &ConfStoreHandleFunc{
			defaultHook: func() *basestore.TransactableHandle {
				panic("unexpected invocation of MockConfStore.Handle")
			},
		},
		SiteCreateIfUpToDateFunc: &ConfStoreSiteCreateIfUpToDateFunc{
			defaultHook: func(context.Context, *int32, string) (*database.SiteConfig, error) {
				panic("unexpected invocation of MockConfStore.SiteCreateIfUpToDate")
			},
		},
		SiteGetLatestFunc: &ConfStoreSiteGetLatestFunc{
			defaultHook: func(context.Context) (*database.SiteConfig, error) {
				panic("unexpected invocation of MockConfStore.SiteGetLatest")
			},
		},
		TransactFunc: &ConfStoreTransactFunc{
			defaultHook: func(context.Context) (database.ConfStore, error) {
				panic("unexpected invocation of MockConfStore.Transact")
			},
		},
	}
}

// NewMockConfStoreFrom creates a new mock of the MockConfStore interface.
// All methods delegate to the given implementation, unless overwritten.
func NewMockConfStoreFrom(i database.ConfStore) *MockConfStore {
	return &MockConfStore{
		DoneFunc: &ConfStoreDoneFunc{
			defaultHook: i.Done,
		},
		HandleFunc: &ConfStoreHandleFunc{
			defaultHook: i.Handle,
		},
		SiteCreateIfUpToDateFunc: &ConfStoreSiteCreateIfUpToDateFunc{
			defaultHook: i.SiteCreateIfUpToDate,
		},
		SiteGetLatestFunc: &ConfStoreSiteGetLatestFunc{
			defaultHook: i.SiteGetLatest,
		},
		TransactFunc: &ConfStoreTransactFunc{
			defaultHook: i.Transact,
		},
	}
}

// ConfStoreDoneFunc describes the behavior when the Done method of the
// parent MockConfStore instance is invoked.
type ConfStoreDoneFunc struct {
	defaultHook func(error) error
	hooks       []func(error) error
	history     []ConfStoreDoneFuncCall
	mutex       sync.Mutex
}

// Done delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockConfStore) Done(v0 error) error {
	r0 := m.DoneFunc.nextHook()(v0)
	m.DoneFunc.appendCall(ConfStoreDoneFuncCall{v0, r0})
	return r0
}

// SetDefaultHook sets function that is called when the Done method of the
// parent MockConfStore instance is invoked and the hook queue is empty.
func (f *ConfStoreDoneFunc) SetDefaultHook(hook func(error) error) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Done method of the parent MockConfStore instance invokes the hook at the
// front of the queue and discards it. After the queue is empty, the default
// hook function is invoked for any future action.
func (f *ConfStoreDoneFunc) PushHook(hook func(error) error) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *ConfStoreDoneFunc) SetDefaultReturn(r0 error) {
	f.SetDefaultHook(func(error) error {
		return r0
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *ConfStoreDoneFunc) PushReturn(r0 error) {
	f.PushHook(func(error) error {
		return r0
	})
}

func (f *ConfStoreDoneFunc) nextHook() func(error) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ConfStoreDoneFunc) appendCall(r0 ConfStoreDoneFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ConfStoreDoneFuncCall objects describing
// the invocations of this function.
func (f *ConfStoreDoneFunc) History() []ConfStoreDoneFuncCall {
	f.mutex.Lock()
	history := make([]ConfStoreDoneFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ConfStoreDoneFuncCall is an object that describes an invocation of method
// Done on an instance of MockConfStore.
type ConfStoreDoneFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 error
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ConfStoreDoneFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ConfStoreDoneFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}

// ConfStoreHandleFunc describes the behavior when the Handle method of the
// parent MockConfStore instance is invoked.
type ConfStoreHandleFunc struct {
	defaultHook func() *basestore.TransactableHandle
	hooks       []func() *basestore.TransactableHandle
	history     []ConfStoreHandleFuncCall
	mutex       sync.Mutex
}

// Handle delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockConfStore) Handle() *basestore.TransactableHandle {
	r0 := m.HandleFunc.nextHook()()
	m.HandleFunc.appendCall(ConfStoreHandleFuncCall{r0})
	return r0
}

// SetDefaultHook sets function that is called when the Handle method of the
// parent MockConfStore instance is invoked and the hook queue is empty.
func (f *ConfStoreHandleFunc) SetDefaultHook(hook func() *basestore.TransactableHandle) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Handle method of the parent MockConfStore instance invokes the hook at
// the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *ConfStoreHandleFunc) PushHook(hook func() *basestore.TransactableHandle) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *ConfStoreHandleFunc) SetDefaultReturn(r0 *basestore.TransactableHandle) {
	f.SetDefaultHook(func() *basestore.TransactableHandle {
		return r0
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *ConfStoreHandleFunc) PushReturn(r0 *basestore.TransactableHandle) {
	f.PushHook(func() *basestore.TransactableHandle {
		return r0
	})
}

func (f *ConfStoreHandleFunc) nextHook() func() *basestore.TransactableHandle {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ConfStoreHandleFunc) appendCall(r0 ConfStoreHandleFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ConfStoreHandleFuncCall objects describing
// the invocations of this function.
func (f *ConfStoreHandleFunc) History() []ConfStoreHandleFuncCall {
	f.mutex.Lock()
	history := make([]ConfStoreHandleFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ConfStoreHandleFuncCall is an object that describes an invocation of
// method Handle on an instance of MockConfStore.
type ConfStoreHandleFuncCall struct {
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 *basestore.TransactableHandle
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ConfStoreHandleFuncCall) Args() []interface{} {
	return []interface{}{}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ConfStoreHandleFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}

// ConfStoreSiteCreateIfUpToDateFunc describes the behavior when the
// SiteCreateIfUpToDate method of the parent MockConfStore instance is
// invoked.
type ConfStoreSiteCreateIfUpToDateFunc struct {
	defaultHook func(context.Context, *int32, string) (*database.SiteConfig, error)
	hooks       []func(context.Context, *int32, string) (*database.SiteConfig, error)
	history     []ConfStoreSiteCreateIfUpToDateFuncCall
	mutex       sync.Mutex
}

// SiteCreateIfUpToDate delegates to the next hook function in the queue and
// stores the parameter and result values of this invocation.
func (m *MockConfStore) SiteCreateIfUpToDate(v0 context.Context, v1 *int32, v2 string) (*database.SiteConfig, error) {
	r0, r1 := m.SiteCreateIfUpToDateFunc.nextHook()(v0, v1, v2)
	m.SiteCreateIfUpToDateFunc.appendCall(ConfStoreSiteCreateIfUpToDateFuncCall{v0, v1, v2, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the SiteCreateIfUpToDate
// method of the parent MockConfStore instance is invoked and the hook queue
// is empty.
func (f *ConfStoreSiteCreateIfUpToDateFunc) SetDefaultHook(hook func(context.Context, *int32, string) (*database.SiteConfig, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// SiteCreateIfUpToDate method of the parent MockConfStore instance invokes
// the hook at the front of the queue and discards it. After the queue is
// empty, the default hook function is invoked for any future action.
func (f *ConfStoreSiteCreateIfUpToDateFunc) PushHook(hook func(context.Context, *int32, string) (*database.SiteConfig, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *ConfStoreSiteCreateIfUpToDateFunc) SetDefaultReturn(r0 *database.SiteConfig, r1 error) {
	f.SetDefaultHook(func(context.Context, *int32, string) (*database.SiteConfig, error) {
		return r0, r1
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *ConfStoreSiteCreateIfUpToDateFunc) PushReturn(r0 *database.SiteConfig, r1 error) {
	f.PushHook(func(context.Context, *int32, string) (*database.SiteConfig, error) {
		return r0, r1
	})
}

func (f *ConfStoreSiteCreateIfUpToDateFunc) nextHook() func(context.Context, *int32, string) (*database.SiteConfig, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ConfStoreSiteCreateIfUpToDateFunc) appendCall(r0 ConfStoreSiteCreateIfUpToDateFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ConfStoreSiteCreateIfUpToDateFuncCall
// objects describing the invocations of this function.
func (f *ConfStoreSiteCreateIfUpToDateFunc) History() []ConfStoreSiteCreateIfUpToDateFuncCall {
	f.mutex.Lock()
	history := make([]ConfStoreSiteCreateIfUpToDateFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ConfStoreSiteCreateIfUpToDateFuncCall is an object that describes an
// invocation of method SiteCreateIfUpToDate on an instance of
// MockConfStore.
type ConfStoreSiteCreateIfUpToDateFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 *int32
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 string
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 *database.SiteConfig
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ConfStoreSiteCreateIfUpToDateFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ConfStoreSiteCreateIfUpToDateFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// ConfStoreSiteGetLatestFunc describes the behavior when the SiteGetLatest
// method of the parent MockConfStore instance is invoked.
type ConfStoreSiteGetLatestFunc struct {
	defaultHook func(context.Context) (*database.SiteConfig, error)
	hooks       []func(context.Context) (*database.SiteConfig, error)
	history     []ConfStoreSiteGetLatestFuncCall
	mutex       sync.Mutex
}

// SiteGetLatest delegates to the next hook function in the queue and stores
// the parameter and result values of this invocation.
func (m *MockConfStore) SiteGetLatest(v0 context.Context) (*database.SiteConfig, error) {
	r0, r1 := m.SiteGetLatestFunc.nextHook()(v0)
	m.SiteGetLatestFunc.appendCall(ConfStoreSiteGetLatestFuncCall{v0, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the SiteGetLatest method
// of the parent MockConfStore instance is invoked and the hook queue is
// empty.
func (f *ConfStoreSiteGetLatestFunc) SetDefaultHook(hook func(context.Context) (*database.SiteConfig, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// SiteGetLatest method of the parent MockConfStore instance invokes the
// hook at the front of the queue and discards it. After the queue is empty,
// the default hook function is invoked for any future action.
func (f *ConfStoreSiteGetLatestFunc) PushHook(hook func(context.Context) (*database.SiteConfig, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *ConfStoreSiteGetLatestFunc) SetDefaultReturn(r0 *database.SiteConfig, r1 error) {
	f.SetDefaultHook(func(context.Context) (*database.SiteConfig, error) {
		return r0, r1
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *ConfStoreSiteGetLatestFunc) PushReturn(r0 *database.SiteConfig, r1 error) {
	f.PushHook(func(context.Context) (*database.SiteConfig, error) {
		return r0, r1
	})
}

func (f *ConfStoreSiteGetLatestFunc) nextHook() func(context.Context) (*database.SiteConfig, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ConfStoreSiteGetLatestFunc) appendCall(r0 ConfStoreSiteGetLatestFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ConfStoreSiteGetLatestFuncCall objects
// describing the invocations of this function.
func (f *ConfStoreSiteGetLatestFunc) History() []ConfStoreSiteGetLatestFuncCall {
	f.mutex.Lock()
	history := make([]ConfStoreSiteGetLatestFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ConfStoreSiteGetLatestFuncCall is an object that describes an invocation
// of method SiteGetLatest on an instance of MockConfStore.
type ConfStoreSiteGetLatestFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 *database.SiteConfig
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ConfStoreSiteGetLatestFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ConfStoreSiteGetLatestFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// ConfStoreTransactFunc describes the behavior when the Transact method of
// the parent MockConfStore instance is invoked.
type ConfStoreTransactFunc struct {
	defaultHook func(context.Context) (database.ConfStore, error)
	hooks       []func(context.Context) (database.ConfStore, error)
	history     []ConfStoreTransactFuncCall
	mutex       sync.Mutex
}

// Transact delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockConfStore) Transact(v0 context.Context) (database.ConfStore, error) {
	r0, r1 := m.TransactFunc.nextHook()(v0)
	m.TransactFunc.appendCall(ConfStoreTransactFuncCall{v0, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the Transact method of
// the parent MockConfStore instance is invoked and the hook queue is empty.
func (f *ConfStoreTransactFunc) SetDefaultHook(hook func(context.Context) (database.ConfStore, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Transact method of the parent MockConfStore instance invokes the hook at
// the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *ConfStoreTransactFunc) PushHook(hook func(context.Context) (database.ConfStore, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *ConfStoreTransactFunc) SetDefaultReturn(r0 database.ConfStore, r1 error) {
	f.SetDefaultHook(func(context.Context) (database.ConfStore, error) {
		return r0, r1
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *ConfStoreTransactFunc) PushReturn(r0 database.ConfStore, r1 error) {
	f.PushHook(func(context.Context) (database.ConfStore, error) {
		return r0, r1
	})
}

func (f *ConfStoreTransactFunc) nextHook() func(context.Context) (database.ConfStore, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ConfStoreTransactFunc) appendCall(r0 ConfStoreTransactFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ConfStoreTransactFuncCall objects
// describing the invocations of this function.
func (f *ConfStoreTransactFunc) History() []ConfStoreTransactFuncCall {
	f.mutex.Lock()
	history := make([]ConfStoreTransactFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ConfStoreTransactFuncCall is an object that describes an invocation of
// method Transact on an instance of MockConfStore.
type ConfStoreTransactFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 database.ConfStore
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ConfStoreTransactFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ConfStoreTransactFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}
