package gitserver

import (
	"context"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sourcegraph/sourcegraph/internal/actor"
	"github.com/sourcegraph/sourcegraph/internal/api"
	"github.com/sourcegraph/sourcegraph/internal/authz"
	"github.com/sourcegraph/sourcegraph/internal/database"
	"github.com/sourcegraph/sourcegraph/internal/gitserver/gitdomain"
	"github.com/sourcegraph/sourcegraph/lib/errors"
)

func TestParseShortLog(t *testing.T) {
	tests := []struct {
		name    string
		input   string // in the format of `git shortlog -sne`
		want    []*gitdomain.PersonCount
		wantErr error
	}{
		{
			name: "basic",
			input: `
  1125	Jane Doe <jane@sourcegraph.com>
   390	Bot Of Doom <bot@doombot.com>
`,
			want: []*gitdomain.PersonCount{
				{
					Name:  "Jane Doe",
					Email: "jane@sourcegraph.com",
					Count: 1125,
				},
				{
					Name:  "Bot Of Doom",
					Email: "bot@doombot.com",
					Count: 390,
				},
			},
		},
		{
			name: "commonly malformed (email address as name)",
			input: `  1125	jane@sourcegraph.com <jane@sourcegraph.com>
   390	Bot Of Doom <bot@doombot.com>
`,
			want: []*gitdomain.PersonCount{
				{
					Name:  "jane@sourcegraph.com",
					Email: "jane@sourcegraph.com",
					Count: 1125,
				},
				{
					Name:  "Bot Of Doom",
					Email: "bot@doombot.com",
					Count: 390,
				},
			},
		},
	}
	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			got, gotErr := parseShortLog([]byte(tst.input))
			if (gotErr == nil) != (tst.wantErr == nil) {
				t.Fatalf("gotErr %+v wantErr %+v", gotErr, tst.wantErr)
			}
			if !reflect.DeepEqual(got, tst.want) {
				t.Logf("got %q", got)
				t.Fatalf("want %q", tst.want)
			}
		})
	}
}

func TestDiff(t *testing.T) {
	ctx := context.Background()
	db := database.NewMockDB()

	t.Run("invalid bases", func(t *testing.T) {
		for _, input := range []string{
			"",
			"-foo",
			".foo",
		} {
			t.Run("invalid base: "+input, func(t *testing.T) {
				i, err := NewClient(db).Diff(ctx, DiffOptions{Base: input})
				if i != nil {
					t.Errorf("unexpected non-nil iterator: %+v", i)
				}
				if err == nil {
					t.Error("unexpected nil error")
				}
			})
		}
	})

	t.Run("rangeSpec calculation", func(t *testing.T) {
		for _, tc := range []struct {
			opts DiffOptions
			want string
		}{
			{opts: DiffOptions{Base: "foo", Head: "bar"}, want: "foo...bar"},
		} {
			t.Run("rangeSpec: "+tc.want, func(t *testing.T) {
				c := NewClient(db)
				Mocks.ExecReader = func(args []string) (reader io.ReadCloser, err error) {
					// The range spec is the sixth argument.
					if args[5] != tc.want {
						t.Errorf("unexpected rangeSpec: have: %s; want: %s", args[5], tc.want)
					}
					return nil, nil
				}
				t.Cleanup(ResetMocks)
				_, _ = c.Diff(ctx, tc.opts)
			})
		}
	})

	t.Run("ExecReader error", func(t *testing.T) {
		c := NewClient(db)
		Mocks.ExecReader = func(args []string) (reader io.ReadCloser, err error) {
			return nil, errors.New("ExecReader error")
		}
		t.Cleanup(ResetMocks)

		i, err := c.Diff(ctx, DiffOptions{Base: "foo", Head: "bar"})
		if i != nil {
			t.Errorf("unexpected non-nil iterator: %+v", i)
		}
		if err == nil {
			t.Error("unexpected nil error")
		}
	})

	t.Run("success", func(t *testing.T) {
		const testDiffFiles = 3
		const testDiff = `diff --git INSTALL.md INSTALL.md
index e5af166..d44c3fc 100644
--- INSTALL.md
+++ INSTALL.md
@@ -3,10 +3,10 @@
 Line 1
 Line 2
 Line 3
-Line 4
+This is cool: Line 4
 Line 5
 Line 6
-Line 7
-Line 8
+Another Line 7
+Foobar Line 8
 Line 9
 Line 10
diff --git JOKES.md JOKES.md
index ea80abf..1b86505 100644
--- JOKES.md
+++ JOKES.md
@@ -4,10 +4,10 @@ Joke #1
 Joke #2
 Joke #3
 Joke #4
-Joke #5
+This is not funny: Joke #5
 Joke #6
-Joke #7
+This one is good: Joke #7
 Joke #8
-Joke #9
+Waffle: Joke #9
 Joke #10
 Joke #11
diff --git README.md README.md
index 9bd8209..d2acfa9 100644
--- README.md
+++ README.md
@@ -1,12 +1,13 @@
 # README

-Line 1
+Foobar Line 1
 Line 2
 Line 3
 Line 4
 Line 5
-Line 6
+Barfoo Line 6
 Line 7
 Line 8
 Line 9
 Line 10
+Another line
`

		var testDiffFileNames = []string{
			"INSTALL.md",
			"JOKES.md",
			"README.md",
		}

		c := NewClient(db)
		Mocks.ExecReader = func(args []string) (reader io.ReadCloser, err error) {
			return io.NopCloser(strings.NewReader(testDiff)), nil
		}
		t.Cleanup(ResetMocks)

		i, err := c.Diff(ctx, DiffOptions{Base: "foo", Head: "bar"})
		if i == nil {
			t.Error("unexpected nil iterator")
		}
		if err != nil {
			t.Errorf("unexpected non-nil error: %+v", err)
		}
		defer i.Close()

		count := 0
		for {
			diff, err := i.Next()
			if err == io.EOF {
				break
			} else if err != nil {
				t.Errorf("unexpected iteration error: %+v", err)
			}

			if diff.OrigName != testDiffFileNames[count] {
				t.Errorf("unexpected diff file name: have: %s; want: %s", diff.OrigName, testDiffFileNames[count])
			}
			count++
		}
		if count != testDiffFiles {
			t.Errorf("unexpected diff count: have %d; want %d", count, testDiffFiles)
		}
	})
}

func TestDiffPath(t *testing.T) {
	testDiff := `
diff --git a/foo.md b/foo.md
index 51a59ef1c..493090958 100644
--- a/foo.md
+++ b/foo.md
@@ -1 +1 @@
-this is my file content
+this is my file contnent
`
	db := database.NewMockDB()
	client := NewClient(db)
	t.Run("basic", func(t *testing.T) {
		Mocks.ExecReader = func(args []string) (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader(testDiff)), nil
		}
		ctx := context.Background()
		checker := authz.NewMockSubRepoPermissionChecker()
		ctx = actor.WithActor(ctx, &actor.Actor{
			UID: 1,
		})
		hunks, err := client.DiffPath(ctx, "", "sourceCommit", "", "file", checker)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if len(hunks) != 1 {
			t.Errorf("unexpected hunks returned: %d", len(hunks))
		}
	})
	t.Run("with sub-repo permissions enabled", func(t *testing.T) {
		Mocks.ExecReader = func(args []string) (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader(testDiff)), nil
		}
		ctx := context.Background()
		checker := authz.NewMockSubRepoPermissionChecker()
		ctx = actor.WithActor(ctx, &actor.Actor{
			UID: 1,
		})
		fileName := "foo"
		checker.EnabledFunc.SetDefaultHook(func() bool {
			return true
		})
		// User doesn't have access to this file
		checker.PermissionsFunc.SetDefaultHook(func(ctx context.Context, i int32, content authz.RepoContent) (authz.Perms, error) {
			if content.Path == fileName {
				return authz.None, nil
			}
			return authz.Read, nil
		})
		hunks, err := client.DiffPath(ctx, "", "sourceCommit", "", fileName, checker)
		if !reflect.DeepEqual(err, os.ErrNotExist) {
			t.Errorf("unexpected error: %s", err)
		}
		if hunks != nil {
			t.Errorf("expected DiffPath to return no results, got %v", hunks)
		}
	})
}

func TestRepository_BlameFile(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	gitCommands := []string{
		"echo line1 > f",
		"git add f",
		"GIT_COMMITTER_NAME=a GIT_COMMITTER_EMAIL=a@a.com GIT_COMMITTER_DATE=2006-01-02T15:04:05Z git commit -m foo --author='a <a@a.com>' --date 2006-01-02T15:04:05Z",
		"echo line2 >> f",
		"git add f",
		"GIT_COMMITTER_NAME=a GIT_COMMITTER_EMAIL=a@a.com GIT_COMMITTER_DATE=2006-01-02T15:04:05Z git commit -m foo --author='a <a@a.com>' --date 2006-01-02T15:04:05Z",
		"git mv f f2",
		"echo line3 >> f2",
		"git add f2",
		"GIT_COMMITTER_NAME=a GIT_COMMITTER_EMAIL=a@a.com GIT_COMMITTER_DATE=2006-01-02T15:04:05Z git commit -m foo --author='a <a@a.com>' --date 2006-01-02T15:04:05Z",
	}
	gitWantHunks := []*Hunk{
		{
			StartLine: 1, EndLine: 2, StartByte: 0, EndByte: 6, CommitID: "e6093374dcf5725d8517db0dccbbf69df65dbde0",
			Message: "foo", Author: gitdomain.Signature{Name: "a", Email: "a@a.com", Date: MustParseTime(time.RFC3339, "2006-01-02T15:04:05Z")},
			Filename: "f",
		},
		{
			StartLine: 2, EndLine: 3, StartByte: 6, EndByte: 12, CommitID: "fad406f4fe02c358a09df0d03ec7a36c2c8a20f1",
			Message: "foo", Author: gitdomain.Signature{Name: "a", Email: "a@a.com", Date: MustParseTime(time.RFC3339, "2006-01-02T15:04:05Z")},
			Filename: "f",
		},
		{
			StartLine: 3, EndLine: 4, StartByte: 12, EndByte: 18, CommitID: "311d75a2b414a77f5158a0ed73ec476f5469b286",
			Message: "foo", Author: gitdomain.Signature{Name: "a", Email: "a@a.com", Date: MustParseTime(time.RFC3339, "2006-01-02T15:04:05Z")},
			Filename: "f2",
		},
	}
	tests := map[string]struct {
		repo api.RepoName
		path string
		opt  *BlameOptions

		wantHunks []*Hunk
	}{
		"git cmd": {
			repo: MakeGitRepository(t, gitCommands...),
			path: "f2",
			opt: &BlameOptions{
				NewestCommit: "master",
			},
			wantHunks: gitWantHunks,
		},
	}

	for label, test := range tests {
		newestCommitID, err := ResolveRevision(ctx, database.NewMockDB(), test.repo, string(test.opt.NewestCommit), ResolveRevisionOptions{})
		if err != nil {
			t.Errorf("%s: ResolveRevision(%q) on base: %s", label, test.opt.NewestCommit, err)
			continue
		}

		test.opt.NewestCommit = newestCommitID
		runBlameFileTest(ctx, t, test.repo, test.path, test.opt, nil, label, test.wantHunks)

		checker := authz.NewMockSubRepoPermissionChecker()
		ctx = actor.WithActor(ctx, &actor.Actor{
			UID: 1,
		})
		// Sub-repo permissions
		// Case: user has read access to file, doesn't filter anything
		checker.EnabledFunc.SetDefaultHook(func() bool {
			return true
		})
		checker.PermissionsFunc.SetDefaultHook(func(ctx context.Context, i int32, content authz.RepoContent) (authz.Perms, error) {
			if content.Path == "f2" {
				return authz.Read, nil
			}
			return authz.None, nil
		})
		runBlameFileTest(ctx, t, test.repo, test.path, test.opt, checker, label, test.wantHunks)

		// Sub-repo permissions
		// Case: user doesn't have access to the file, nothing returned.
		checker.PermissionsFunc.SetDefaultHook(func(ctx context.Context, i int32, content authz.RepoContent) (authz.Perms, error) {
			return authz.None, nil
		})
		runBlameFileTest(ctx, t, test.repo, test.path, test.opt, checker, label, nil)
	}
}

func runBlameFileTest(ctx context.Context, t *testing.T, repo api.RepoName, path string, opt *BlameOptions,
	checker authz.SubRepoPermissionChecker, label string, wantHunks []*Hunk) {
	t.Helper()
	hunks, err := BlameFile(ctx, database.NewMockDB(), repo, path, opt, checker)
	if err != nil {
		t.Errorf("%s: BlameFile(%s, %+v): %s", label, path, opt, err)
		return
	}
	if !reflect.DeepEqual(hunks, wantHunks) {
		t.Errorf("%s: hunks != wantHunks\n\nhunks ==========\n%s\n\nwantHunks ==========\n%s", label, AsJSON(hunks), AsJSON(wantHunks))
	}
}

func TestIsAbsoluteRevision(t *testing.T) {
	yes := []string{"8cb03d28ad1c6a875f357c5d862237577b06e57c", "20697a062454c29d84e3f006b22eb029d730cd00"}
	no := []string{"ref: refs/heads/appsinfra/SHEP-20-review", "master", "HEAD", "refs/heads/master", "20697a062454c29d84e3f006b22eb029d730cd0", "20697a062454c29d84e3f006b22eb029d730cd000", "  20697a062454c29d84e3f006b22eb029d730cd00  ", "20697a062454c29d84e3f006b22eb029d730cd0 "}
	for _, s := range yes {
		if !IsAbsoluteRevision(s) {
			t.Errorf("%q should be an absolute revision", s)
		}
	}
	for _, s := range no {
		if IsAbsoluteRevision(s) {
			t.Errorf("%q should not be an absolute revision", s)
		}
	}
}

func TestRepository_ResolveBranch(t *testing.T) {
	ClientMocks.LocalGitserver = true
	defer ResetClientMocks()

	gitCommands := []string{
		"GIT_COMMITTER_NAME=a GIT_COMMITTER_EMAIL=a@a.com GIT_COMMITTER_DATE=2006-01-02T15:04:05Z git commit --allow-empty -m foo --author='a <a@a.com>' --date 2006-01-02T15:04:05Z",
	}
	tests := map[string]struct {
		repo         api.RepoName
		branch       string
		wantCommitID api.CommitID
	}{
		"git cmd": {
			repo:         MakeGitRepositoryAndSetReposDir(t, gitCommands...),
			branch:       "master",
			wantCommitID: "ea167fe3d76b1e5fd3ed8ca44cbd2fe3897684f8",
		},
	}

	for label, test := range tests {
		commitID, err := ResolveRevision(context.Background(), database.NewMockDB(), test.repo, test.branch, ResolveRevisionOptions{})
		if err != nil {
			t.Errorf("%s: ResolveRevision: %s", label, err)
			continue
		}

		if commitID != test.wantCommitID {
			t.Errorf("%s: got commitID == %v, want %v", label, commitID, test.wantCommitID)
		}
	}
}

func TestRepository_ResolveBranch_error(t *testing.T) {
	ClientMocks.LocalGitserver = true
	defer ResetClientMocks()

	gitCommands := []string{
		"GIT_COMMITTER_NAME=a GIT_COMMITTER_EMAIL=a@a.com GIT_COMMITTER_DATE=2006-01-02T15:04:05Z git commit --allow-empty -m foo --author='a <a@a.com>' --date 2006-01-02T15:04:05Z",
	}
	tests := map[string]struct {
		repo    api.RepoName
		branch  string
		wantErr func(error) bool
	}{
		"git cmd": {
			repo:    MakeGitRepositoryAndSetReposDir(t, gitCommands...),
			branch:  "doesntexist",
			wantErr: func(err error) bool { return errors.HasType(err, &gitdomain.RevisionNotFoundError{}) },
		},
	}

	for label, test := range tests {
		commitID, err := ResolveRevision(context.Background(), database.NewMockDB(), test.repo, test.branch, ResolveRevisionOptions{})
		if !test.wantErr(err) {
			t.Errorf("%s: ResolveRevision: %s", label, err)
			continue
		}

		if commitID != "" {
			t.Errorf("%s: got commitID == %v, want empty", label, commitID)
		}
	}
}

func TestRepository_ResolveTag(t *testing.T) {
	ClientMocks.LocalGitserver = true
	defer ResetClientMocks()

	gitCommands := []string{
		"GIT_COMMITTER_NAME=a GIT_COMMITTER_EMAIL=a@a.com GIT_COMMITTER_DATE=2006-01-02T15:04:05Z git commit --allow-empty -m foo --author='a <a@a.com>' --date 2006-01-02T15:04:05Z",
		"git tag t",
	}
	tests := map[string]struct {
		repo         api.RepoName
		tag          string
		wantCommitID api.CommitID
	}{
		"git cmd": {
			repo:         MakeGitRepositoryAndSetReposDir(t, gitCommands...),
			tag:          "t",
			wantCommitID: "ea167fe3d76b1e5fd3ed8ca44cbd2fe3897684f8",
		},
	}

	for label, test := range tests {
		commitID, err := ResolveRevision(context.Background(), database.NewMockDB(), test.repo, test.tag, ResolveRevisionOptions{})
		if err != nil {
			t.Errorf("%s: ResolveRevision: %s", label, err)
			continue
		}

		if commitID != test.wantCommitID {
			t.Errorf("%s: got commitID == %v, want %v", label, commitID, test.wantCommitID)
		}
	}
}

func TestRepository_ResolveTag_error(t *testing.T) {
	ClientMocks.LocalGitserver = true
	defer ResetClientMocks()

	gitCommands := []string{
		"GIT_COMMITTER_NAME=a GIT_COMMITTER_EMAIL=a@a.com GIT_COMMITTER_DATE=2006-01-02T15:04:05Z git commit --allow-empty -m foo --author='a <a@a.com>' --date 2006-01-02T15:04:05Z",
	}
	tests := map[string]struct {
		repo    api.RepoName
		tag     string
		wantErr func(error) bool
	}{
		"git cmd": {
			repo:    MakeGitRepositoryAndSetReposDir(t, gitCommands...),
			tag:     "doesntexist",
			wantErr: func(err error) bool { return errors.HasType(err, &gitdomain.RevisionNotFoundError{}) },
		},
	}

	for label, test := range tests {
		commitID, err := ResolveRevision(context.Background(), database.NewMockDB(), test.repo, test.tag, ResolveRevisionOptions{})
		if !test.wantErr(err) {
			t.Errorf("%s: ResolveRevision: %s", label, err)
			continue
		}

		if commitID != "" {
			t.Errorf("%s: got commitID == %v, want empty", label, commitID)
		}
	}
}
