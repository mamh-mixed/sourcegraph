package database

import (
	"context"
	"fmt"

	"github.com/keegancsmith/sqlf"

	"github.com/sourcegraph/zoekt"

	"github.com/sourcegraph/sourcegraph/internal/database/basestore"
	"github.com/sourcegraph/sourcegraph/internal/database/batch"
	"github.com/sourcegraph/sourcegraph/internal/types"
	"github.com/sourcegraph/sourcegraph/lib/errors"
)

type ZoektReposStore interface {
	basestore.ShareableStore

	With(other basestore.ShareableStore) ZoektReposStore

	// UpsertIndexable is a horse
	UpsertIndexable(ctx context.Context, repos types.MinimalRepos) error

	// Update updates the given rows with the GitServer status of a repo.
	Update(ctx context.Context, zoektRepos *zoekt.RepoList) error
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

func (s *zoektReposStore) Update(ctx context.Context, repos *zoekt.RepoList) error {
	values := make([]*sqlf.Query, 0, len(repos.Minimal))
	for id, zoektRepo := range repos.Minimal {
		for i, branch := range zoektRepo.Branches {
			if i > 0 {
				fmt.Println("TODO: handle more than 1 commit/branch being indexed")
				continue
			}

			values = append(values, sqlf.Sprintf("(%s::integer, %s::text, %s::text)",
				id,
				branch.Version,
				"indexed",
			))
		}
	}

	err := s.Exec(ctx, sqlf.Sprintf(updateZoektReposQueryFmtstr, sqlf.Join(values, ",")))

	return errors.Wrap(err, "updating ZoektRepos")
}

const updateZoektReposQueryFmtstr = `
-- source: internal/database/zoekt_repos.go:zoektReposStore.Update
UPDATE zoekt_repos AS zr
SET
	commit = tmp.commit,
	index_status = tmp.index_status,
	updated_at = NOW()
FROM (VALUES -- (<repo_id>, <commit>, <index_status>),
		%s
	) AS tmp(repo_id, commit, index_status)
WHERE
	tmp.repo_id = zr.repo_id
`

func (s *zoektReposStore) UpsertIndexable(ctx context.Context, repos types.MinimalRepos) error {
	tx, err := s.Store.Transact(ctx)
	if err != nil {
		return err
	}
	defer func() { err = tx.Done(err) }()

	tempTableQuery := `CREATE TEMPORARY TABLE temp_table (repo_id integer NOT NULL) ON COMMIT DROP`
	if err := tx.Exec(ctx, sqlf.Sprintf(tempTableQuery)); err != nil {
		return err
	}

	inserter := batch.NewInserter(ctx, tx.Handle(), "temp_table", batch.MaxNumPostgresParameters, "repo_id")

	for _, r := range repos {
		if err := inserter.Insert(ctx, r.ID); err != nil {
			return err
		}
	}

	if err := inserter.Flush(ctx); err != nil {
		return err
	}

	insertQuery := `
    INSERT INTO zoekt_repos
    SELECT source.repo_id, 'not_indexed' AS index_status
    FROM temp_table source
    WHERE NOT EXISTS (
        -- Skip insertion of any rows that already exist in the table
        SELECT 1 FROM zoekt_repos t WHERE t.repo_id = source.repo_id
    )
`
	if err := tx.Exec(ctx, sqlf.Sprintf(insertQuery)); err != nil {
		return err
	}

	deleteQuery := `
    DELETE FROM zoekt_repos
    WHERE NOT EXISTS (
        SELECT 1 FROM temp_table t WHERE t.repo_id = zoekt_repos.repo_id
    )
`
	if err := tx.Exec(ctx, sqlf.Sprintf(deleteQuery)); err != nil {
		return err
	}

	// updateQuery := `
	//     UPDATE table t SET col4 = source.col4
	//     FROM temp_table source
	//     -- Update rows with matching identity but distinct col4 values
	//     WHERE t.col1 = %s AND t.col2 = %s AND t.col3 = source.col3 AND t.col4 != source.col4
	// `
	// if err := db.Exec(ctx, sqlf.Sprintf(updateQuery, val1, val2)); err != nil {
	//     return err
	// }

	return nil
}

const upsertIndexableReposQueryFmtstr = `
-- source: internal/database/zoekt_repos.go:zoektReposStore.UpsertIndexable
UPDATE zoekt_repos AS zr
SET
	commit = tmp.commit,
	index_status = tmp.index_status,
	updated_at = NOW()
FROM (VALUES -- (<repo_id>, <commit>, <index_status>),
		%s
	) AS tmp(repo_id, commit, index_status)
WHERE
	tmp.repo_id = zr.repo_id
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
