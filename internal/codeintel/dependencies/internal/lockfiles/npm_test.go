package lockfiles

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseYarnLockFile(t *testing.T) {
	tests := []struct {
		lockfile  string
		wantDeps  string
		wantGraph string
	}{
		{
			lockfile: "testdata/parse/yarn.lock/yarn_graph.lock",
			wantDeps: `@types/tinycolor2 1.4.3
ansi-styles 4.3.0
chalk 4.1.2
color-convert 2.0.1
color-name 1.1.4
gradient-string 2.0.1
has-flag 4.0.0
supports-color 7.2.0
tinycolor2 1.4.2
tinygradient 1.1.5
`,
			wantGraph: `npm/gradient-string:
	npm/tinygradient:
		npm/types/tinycolor2
		npm/tinycolor2
	npm/chalk:
		npm/supports-color:
			npm/has-flag
		npm/ansi-styles:
			npm/color-convert:
				npm/color-name
`,
		},
		{
			lockfile: "testdata/parse/yarn.lock/yarn_normal.lock",
			wantDeps: `asap 2.0.6
jquery 3.4.1
promise 8.0.3
`,
			wantGraph: `npm/promise:
	npm/asap
npm/jquery
`,
		},
		{
			lockfile: "testdata/parse/yarn.lock/yarn_lock_thorsten.lock",
			wantDeps: `@babel/code-frame 7.16.7
@babel/helper-validator-identifier 7.16.7
@babel/highlight 7.17.12
@eslint/eslintrc 1.3.0
@humanwhocodes/config-array 0.9.5
@humanwhocodes/object-schema 1.2.1
@types/minimist 1.2.2
@types/normalize-package-data 2.4.1
@types/tinycolor2 1.4.3
acorn-jsx 5.3.2
acorn 8.7.1
ajv 6.12.6
ansi-regex 5.0.1
ansi-styles 3.2.1
ansi-styles 4.3.0
argparse 2.0.1
arrify 1.0.1
balanced-match 1.0.2
brace-expansion 1.1.11
callsites 3.1.0
camelcase-keys 7.0.2
camelcase 6.3.0
chalk-animation 2.0.2
chalk 2.4.2
chalk 4.1.2
color-convert 1.9.3
color-convert 2.0.1
color-name 1.1.3
color-name 1.1.4
concat-map 0.0.1
cross-spawn 7.0.3
debug 4.3.4
decamelize-keys 1.1.0
decamelize 1.2.0
decamelize 5.0.1
deep-is 0.1.4
doctrine 3.0.0
error-ex 1.3.2
escape-string-regexp 1.0.5
escape-string-regexp 4.0.0
eslint-scope 7.1.1
eslint-utils 3.0.0
eslint-visitor-keys 2.1.0
eslint-visitor-keys 3.3.0
eslint 8.18.0
espree 9.3.2
esquery 1.4.0
esrecurse 4.3.0
estraverse 5.3.0
esutils 2.0.3
fast-deep-equal 3.1.3
fast-json-stable-stringify 2.1.0
fast-levenshtein 2.0.6
file-entry-cache 6.0.1
find-up 5.0.0
flat-cache 3.0.4
flatted 3.2.5
fs.realpath 1.0.0
function-bind 1.1.1
functional-red-black-tree 1.0.1
glob-parent 6.0.2
glob 7.2.3
globals 13.15.0
gradient-string 2.0.1
hard-rejection 2.1.0
has-flag 3.0.0
has-flag 4.0.0
has 1.0.3
hosted-git-info 4.1.0
ignore 5.2.0
import-fresh 3.3.0
imurmurhash 0.1.4
indent-string 5.0.0
inflight 1.0.6
inherits 2.0.4
is-arrayish 0.2.1
is-core-module 2.9.0
is-extglob 2.1.1
is-glob 4.0.3
is-plain-obj 1.1.0
isexe 2.0.0
js-tokens 4.0.0
js-yaml 4.1.0
json-parse-even-better-errors 2.3.1
json-schema-traverse 0.4.1
json-stable-stringify-without-jsonify 1.0.1
kind-of 6.0.3
levn 0.4.1
lines-and-columns 1.2.4
locate-path 6.0.0
lodash.merge 4.6.2
lru-cache 6.0.0
map-obj 1.0.1
map-obj 4.3.0
meow 10.1.2
min-indent 1.0.1
minimatch 3.1.2
minimist-options 4.1.0
ms 2.1.2
natural-compare 1.4.0
normalize-package-data 3.0.3
once 1.4.0
optionator 0.9.1
p-limit 3.1.0
p-locate 5.0.0
parent-module 1.0.1
parse-json 5.2.0
path-exists 4.0.0
path-is-absolute 1.0.1
path-key 3.1.1
prelude-ls 1.2.1
punycode 2.1.1
quick-lru 5.1.1
read-pkg-up 8.0.0
read-pkg 6.0.0
redent 4.0.0
regexpp 3.2.0
resolve-from 4.0.0
rimraf 3.0.2
semver 7.3.7
shebang-command 2.0.0
shebang-regex 3.0.0
spdx-correct 3.1.1
spdx-exceptions 2.3.0
spdx-expression-parse 3.0.1
spdx-license-ids 3.0.11
strip-ansi 6.0.1
strip-indent 4.0.0
strip-json-comments 3.1.1
supports-color 5.5.0
supports-color 7.2.0
text-table 0.2.0
tinycolor2 1.4.2
tinygradient 1.1.5
trim-newlines 4.0.2
type-check 0.4.0
type-fest 0.20.2
type-fest 1.4.0
uri-js 4.4.1
v8-compile-cache 2.3.0
validate-npm-package-license 3.0.4
which 2.0.2
word-wrap 1.2.3
wrappy 1.0.2
yallist 4.0.0
yargs-parser 20.2.9
yocto-queue 0.1.0
`,
			wantGraph: `npm/promise:
	npm/asap
npm/jquery
`,
		},
	}

	for _, tt := range tests {
		yarnLockFile, err := os.ReadFile(tt.lockfile)
		if err != nil {
			t.Fatal(err)
		}

		r := strings.NewReader(string(yarnLockFile))

		deps, graph, err := parseYarnLockFile(r)
		if err != nil {
			t.Fatal(err)
		}

		buf := bytes.Buffer{}
		for _, dep := range deps {
			_, err := fmt.Fprintf(&buf, "%s %s\n", dep.PackageSyntax(), dep.PackageVersion())
			if err != nil {
				t.Fatal()
			}
		}
		got := buf.String()

		if d := cmp.Diff(tt.wantDeps, got); d != "" {
			t.Fatalf("wrong deps: +want,-got\n%s", d)
		}

		gotGraph := graph.String()
		if d := cmp.Diff(tt.wantGraph, gotGraph); d != "" {
			t.Fatalf("wrong graph: +want,-got\n%s", d)
		}
	}
}
