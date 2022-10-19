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
	"github.com/sourcegraph/sourcegraph/lib/errors"
)

type ZoektReposStore interface {
	basestore.ShareableStore

	With(other basestore.ShareableStore) ZoektReposStore

	// UpdateIndexStatuses is a horse
	UpdateIndexStatuses(ctx context.Context, indexed map[uint32]*zoekt.MinimalRepoListEntry) error

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
	var commit string

	err := sc.Scan(
		&zr.RepoID,
		&dbutil.NullString{S: &commit},
		&zr.IndexStatus,
		&zr.UpdatedAt,
		&zr.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	zr.Commit = api.CommitID(commit)

	return &zr, nil
}

const getZoektRepoQueryFmtstr = `
-- source: internal/database/zoekt_repos.go:zoektReposStore.GetZoektRepo
SELECT
	zr.repo_id,
	zr.commit,
	zr.index_status,
	zr.updated_at,
	zr.created_at
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

func (s *zoektReposStore) UpdateIndexStatuses(ctx context.Context, indexed map[uint32]*zoekt.MinimalRepoListEntry) error {
	tx, err := s.Store.Transact(ctx)
	if err != nil {
		return err
	}
	defer func() { err = tx.Done(err) }()

	if err := tx.Exec(ctx, sqlf.Sprintf(updateIndexStatusesCreateTempTableQuery)); err != nil {
		return err
	}

	inserter := batch.NewInserter(ctx, tx.Handle(), "temp_table", batch.MaxNumPostgresParameters, tempTableColumns...)

	for repoID, entry := range indexed {
		commit := ""
		indexStatus := "indexed"

		for i, branch := range entry.Branches {
			if i != 0 {
				fmt.Printf("TODO: only persisting one branch, ignoring: %+v\n", branch)
				continue
			}
			commit = branch.Version
		}

		if err := inserter.Insert(ctx, repoID, indexStatus, dbutil.NullStringColumn(commit)); err != nil {
			return err
		}
	}

	if err := inserter.Flush(ctx); err != nil {
		return err
	}

	if err := tx.Exec(ctx, sqlf.Sprintf(updateIndexStatusesUpdateQuery)); err != nil {
		return errors.Wrap(err, "updating zoekt repos failed")
	}

	return nil
}

var tempTableColumns = []string{
	"repo_id",
	"index_status",
	"commit",
}

const updateIndexStatusesCreateTempTableQuery = `
CREATE TEMPORARY TABLE temp_table (
	repo_id integer NOT NULL,
	index_status text NOT NULL,
	commit text
) ON COMMIT DROP
`

const updateIndexStatusesUpdateQuery = `
UPDATE zoekt_repos t
SET
	index_status = source.index_status,
	commit       = source.commit,
	updated_at   = now()
FROM temp_table source
WHERE
	t.repo_id = source.repo_id
AND
	(t.index_status != source.index_status OR t.commit != source.commit)
`

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
