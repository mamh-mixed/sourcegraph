package database

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/keegancsmith/sqlf"
	"github.com/sourcegraph/log/logtest"
	"github.com/sourcegraph/zoekt"

	"github.com/sourcegraph/sourcegraph/internal/api"
	"github.com/sourcegraph/sourcegraph/internal/database/basestore"
	"github.com/sourcegraph/sourcegraph/internal/database/dbtest"
	"github.com/sourcegraph/sourcegraph/internal/database/dbutil"
	"github.com/sourcegraph/sourcegraph/internal/types"
)

func TestZoektRepos_GetZoektRepo(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	logger := logtest.Scoped(t)
	db := NewDB(logger, dbtest.NewDB(logger, t))
	ctx := context.Background()
	s := &zoektReposStore{Store: basestore.NewWithHandle(db.Handle())}

	repo1, _ := createTestRepo(ctx, t, db, &createTestRepoPayload{Name: "repo1"})
	repo2, _ := createTestRepo(ctx, t, db, &createTestRepoPayload{Name: "repo2"})
	repo3, _ := createTestRepo(ctx, t, db, &createTestRepoPayload{Name: "repo3"})

	insertZoektRepo := func(r api.RepoID, indexStatus string, commit *string) {
		err := s.Exec(ctx, sqlf.Sprintf("INSERT INTO zoekt_repos (repo_id, index_status, commit) VALUES (%s, %s, %s)", r, indexStatus, dbutil.NullStringColumn(commit)))
		if err != nil {
			t.Fatalf("failed to query repo name: %s", err)
		}
	}
	strPtr := func(str string) *string {
		if str == "" {
			return nil
		}
		return &str
	}

	insertZoektRepo(repo1.ID, "not_indexed", nil)
	insertZoektRepo(repo2.ID, "indexed", strPtr("d34db33f"))
	insertZoektRepo(repo3.ID, "indexed", strPtr("c4f3b4b3"))
}

func TestZoektRepos_UpsertIndexable(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	logger := logtest.Scoped(t)
	db := NewDB(logger, dbtest.NewDB(logger, t))
	ctx := context.Background()
	s := &zoektReposStore{Store: basestore.NewWithHandle(db.Handle())}

	var repos types.MinimalRepos
	for _, name := range []api.RepoName{
		"repo1",
		"repo2",
		"repo3",
	} {
		r, _ := createTestRepo(ctx, t, db, &createTestRepoPayload{Name: name})
		repos = append(repos, types.MinimalRepo{ID: r.ID, Name: r.Name})
	}

	indexed := map[uint32]*zoekt.MinimalRepoListEntry{
		uint32(repos[0].ID): {Branches: []zoekt.RepositoryBranch{{Name: "main", Version: "d34db33f"}}},
	}

	if err := s.UpsertIndexable(ctx, repos, indexed); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	assertZoektRepoStatistics(t, ctx, s, ZoektRepoStatistics{
		Total:      3,
		Indexed:    1,
		NotIndexed: 2,
	})
}

func assertZoektRepoStatistics(t *testing.T, ctx context.Context, s *zoektReposStore, wantZoektStats ZoektRepoStatistics) {
	t.Helper()

	stats, err := s.GetStatistics(ctx)
	if err != nil {
		t.Fatalf("zoektRepoStore.GetStatistics failed: %s", err)
	}

	if diff := cmp.Diff(stats, wantZoektStats); diff != "" {
		t.Errorf("ZoektRepoStatistics differ: %s", diff)
	}
}
