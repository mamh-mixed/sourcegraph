package database

import (
	"context"
	"fmt"
	"time"

	"github.com/keegancsmith/sqlf"

	"github.com/sourcegraph/zoekt"

	"github.com/sourcegraph/sourcegraph/internal/api"
	"github.com/sourcegraph/sourcegraph/internal/database/basestore"
	"github.com/sourcegraph/sourcegraph/internal/database/batch"
	"github.com/sourcegraph/sourcegraph/internal/database/dbutil"
	"github.com/sourcegraph/sourcegraph/internal/types"
	"github.com/sourcegraph/sourcegraph/lib/errors"
)

type ZoektReposStore interface {
	basestore.ShareableStore

	With(other basestore.ShareableStore) ZoektReposStore

	// UpsertIndexable is a horse
	UpsertIndexable(ctx context.Context, repos types.MinimalRepos, indexed map[uint32]*zoekt.MinimalRepoListEntry) error

	// Update updates the given rows with the GitServer status of a repo.
	GetStatistics(ctx context.Context) (ZoektRepoStatistics, error)
}

var _ ZoektReposStore = (*zoektReposStore)(nil)

// zoektReposStore is responsible for data stored in the gitserver_repos table.
type zoektReposStore struct {
	*basestore.Store
}

// ZoektRepossWith instantiates and returns a new zoektReposStore using
// the other store handle.
func ZoektReposWith(other basestore.ShareableStore) ZoektReposStore {
	return &zoektReposStore{Store: basestore.NewWithHandle(other.Handle())}
}

func (s *zoektReposStore) With(other basestore.ShareableStore) ZoektReposStore {
	return &zoektReposStore{Store: s.Store.With(other)}
}

func (s *zoektReposStore) Transact(ctx context.Context) (ZoektReposStore, error) {
	txBase, err := s.Store.Transact(ctx)
	return &zoektReposStore{Store: txBase}, err
}

type ZoektRepo struct {
	RepoID      api.RepoID
	Commit      api.CommitID
	IndexStatus string

	UpdatedAt time.Time
	CreatedAt time.Time
}

func (s *zoektReposStore) GetZoektRepo(ctx context.Context, repo api.RepoID) (*ZoektRepo, error) {
	return scanZoektRepo(s.QueryRow(ctx, sqlf.Sprintf(getZoektRepoQueryFmtstr, repo)))
}

func scanZoektRepo(sc dbutil.Scanner) (*ZoektRepo, error) {
	var zr ZoektRepo
	return &zr, sc.Scan(
		&zr.RepoID,
		&zr.Commit,
		&zr.IndexStatus,
		&zr.UpdatedAt,
		&zr.CreatedAt,
	)
}

const getZoektRepoQueryFmtstr = `
-- source: internal/database/zoekt_repos.go:zoektReposStore.GetZoektRepo
SELECT
	repo_id,
	commit,
	index_status,
	updated_at,
	created_at
FROM zoekt_repos zr
JOIN repo ON repo.id = zr.repo_id
WHERE
	repo.deleted_at is NULL
AND
	repo.blocked IS NULL
AND
	zr.repo_id = %s
;
`

func (s *zoektReposStore) UpsertIndexable(ctx context.Context, repos types.MinimalRepos, indexed map[uint32]*zoekt.MinimalRepoListEntry) error {
	tx, err := s.Store.Transact(ctx)
	if err != nil {
		return err
	}
	defer func() { err = tx.Done(err) }()

	tempTableQuery := `CREATE TEMPORARY TABLE temp_table (
		repo_id integer NOT NULL,
		index_status text NOT NULL,
		commit text
	) ON COMMIT DROP`
	if err := tx.Exec(ctx, sqlf.Sprintf(tempTableQuery)); err != nil {
		return err
	}

	inserter := batch.NewInserter(ctx, tx.Handle(), "temp_table", batch.MaxNumPostgresParameters, "repo_id", "index_status", "commit")

	for _, r := range repos {
		indexStatus := "not_indexed"
		commit := ""

		indexedEntry, ok := indexed[uint32(r.ID)]
		if ok {
			indexStatus = "indexed"
			for i, branch := range indexedEntry.Branches {
				if i != 0 {
					fmt.Printf("TODO: only persisting one branch, ignoring: %+v\n", branch)
					continue
				}
				commit = branch.Version
			}
		}

		commitColumn := func() *string {
			if commit == "" {
				return nil
			}
			return &commit
		}

		if err := inserter.Insert(ctx, r.ID, indexStatus, commitColumn()); err != nil {
			return err
		}
	}

	if err := inserter.Flush(ctx); err != nil {
		return err
	}

	insertQuery := `
    INSERT INTO zoekt_repos (repo_id, index_status, commit)
    SELECT source.repo_id, source.index_status, source.commit
    FROM temp_table source
    WHERE NOT EXISTS (
        -- Skip insertion of any rows that already exist in the table
        SELECT 1 FROM zoekt_repos t WHERE t.repo_id = source.repo_id
    )`
	if err := tx.Exec(ctx, sqlf.Sprintf(insertQuery)); err != nil {
		return errors.Wrap(err, "inserting zoekt_repos failed")
	}

	deleteQuery := `
    DELETE FROM zoekt_repos
    WHERE NOT EXISTS (
        SELECT 1 FROM temp_table t WHERE t.repo_id = zoekt_repos.repo_id
    )`
	if err := tx.Exec(ctx, sqlf.Sprintf(deleteQuery)); err != nil {
		return errors.Wrap(err, "deleting zoekt_repos failed")
	}

	updateQuery := `
    UPDATE zoekt_repos t
	SET
		index_status = source.index_status,
		commit       = source.commit,
		updated_at   = now()
    FROM temp_table source
    WHERE t.repo_id = source.repo_id AND (t.index_status != source.index_status OR t.commit != source.commit)
	`
	if err := tx.Exec(ctx, sqlf.Sprintf(updateQuery)); err != nil {
		return errors.Wrap(err, "updating zoekt repos failed")
	}

	return nil
}

type ZoektRepoStatistics struct {
	Total      int
	Indexed    int
	NotIndexed int
}

func (s *zoektReposStore) GetStatistics(ctx context.Context) (ZoektRepoStatistics, error) {
	var zrs ZoektRepoStatistics
	row := s.QueryRow(ctx, sqlf.Sprintf(getZoektRepoStatisticsQueryFmtstr))
	err := row.Scan(&zrs.Total, &zrs.Indexed, &zrs.NotIndexed)
	if err != nil {
		return zrs, err
	}
	return zrs, nil
}

const getZoektRepoStatisticsQueryFmtstr = `
-- source: internal/database/zoekt_repos.go:zoektReposStore.GetStatistics
SELECT
	COUNT(*) AS total,
	COUNT(*) FILTER(WHERE index_status = 'indexed') AS indexed,
	COUNT(*) FILTER(WHERE index_status = 'not_indexed') AS not_indexed
FROM zoekt_repos zr
JOIN repo ON repo.id = zr.repo_id
WHERE
	repo.deleted_at is NULL
AND
	repo.blocked IS NULL
;
`
