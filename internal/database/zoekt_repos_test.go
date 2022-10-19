package database

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sourcegraph/log/logtest"
	"github.com/sourcegraph/zoekt"
	"github.com/stretchr/testify/assert"

	"github.com/sourcegraph/sourcegraph/internal/api"
	"github.com/sourcegraph/sourcegraph/internal/database/basestore"
	"github.com/sourcegraph/sourcegraph/internal/database/dbtest"
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

	assertZoektRepos(t, ctx, s, map[api.RepoID]*ZoektRepo{
		repo1.ID: {RepoID: repo1.ID, IndexStatus: "not_indexed", Commit: ""},
		repo2.ID: {RepoID: repo2.ID, IndexStatus: "not_indexed", Commit: ""},
		repo3.ID: {RepoID: repo3.ID, IndexStatus: "not_indexed", Commit: ""},
	})
}

func TestZoektRepos_UpdateIndexStatuses(t *testing.T) {
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

	// No repo is indexed
	assertZoektRepoStatistics(t, ctx, s, ZoektRepoStatistics{
		Total:      3,
		NotIndexed: 3,
	})

	assertZoektRepos(t, ctx, s, map[api.RepoID]*ZoektRepo{
		repos[0].ID: {RepoID: repos[0].ID, IndexStatus: "not_indexed"},
		repos[1].ID: {RepoID: repos[1].ID, IndexStatus: "not_indexed"},
		repos[2].ID: {RepoID: repos[2].ID, IndexStatus: "not_indexed"},
	})

	// 1/3 repo is indexed
	indexed := map[uint32]*zoekt.MinimalRepoListEntry{
		uint32(repos[0].ID): {Branches: []zoekt.RepositoryBranch{{Name: "main", Version: "d34db33f"}}},
	}

	if err := s.UpdateIndexStatuses(ctx, indexed); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	assertZoektRepoStatistics(t, ctx, s, ZoektRepoStatistics{
		Total:      3,
		Indexed:    1,
		NotIndexed: 2,
	})

	assertZoektRepos(t, ctx, s, map[api.RepoID]*ZoektRepo{
		repos[0].ID: {RepoID: repos[0].ID, IndexStatus: "indexed", Commit: "d34db33f"},
		repos[1].ID: {RepoID: repos[1].ID, IndexStatus: "not_indexed"},
		repos[2].ID: {RepoID: repos[2].ID, IndexStatus: "not_indexed"},
	})

	// Index all repositories
	indexed = map[uint32]*zoekt.MinimalRepoListEntry{
		// different commit
		uint32(repos[0].ID): {Branches: []zoekt.RepositoryBranch{{Name: "main", Version: "f00b4r"}}},
		// new
		uint32(repos[1].ID): {Branches: []zoekt.RepositoryBranch{{Name: "main-2", Version: "b4rf00"}}},
		// new
		uint32(repos[2].ID): {Branches: []zoekt.RepositoryBranch{{Name: "main", Version: "d00d00"}}},
	}

	if err := s.UpdateIndexStatuses(ctx, indexed); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	assertZoektRepoStatistics(t, ctx, s, ZoektRepoStatistics{
		Total:   3,
		Indexed: 3,
	})

	assertZoektRepos(t, ctx, s, map[api.RepoID]*ZoektRepo{
		repos[0].ID: {RepoID: repos[0].ID, IndexStatus: "indexed", Commit: "f00b4r"},
		repos[1].ID: {RepoID: repos[1].ID, IndexStatus: "indexed", Commit: "b4rf00"},
		repos[2].ID: {RepoID: repos[2].ID, IndexStatus: "indexed", Commit: "d00d00"},
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

func assertZoektRepos(t *testing.T, ctx context.Context, s *zoektReposStore, want map[api.RepoID]*ZoektRepo) {
	t.Helper()

	for repoID, w := range want {
		have, err := s.GetZoektRepo(ctx, repoID)
		if err != nil {
			t.Fatalf("unexpected error from GetZoektRepo: %s", err)
		}

		assert.NotZero(t, have.UpdatedAt)
		assert.NotZero(t, have.CreatedAt)

		w.UpdatedAt = have.UpdatedAt
		w.CreatedAt = have.CreatedAt

		if diff := cmp.Diff(have, w); diff != "" {
			t.Errorf("ZoektRepo for repo %d differs: %s", repoID, diff)
		}
	}
}
