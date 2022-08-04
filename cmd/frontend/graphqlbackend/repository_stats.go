package graphqlbackend

import (
	"context"

	"github.com/sourcegraph/sourcegraph/cmd/frontend/backend"
	"github.com/sourcegraph/sourcegraph/internal/usagestats"
)

type repositoryStatsResolver struct {
	gitDirBytes       uint64
	indexedLinesCount uint64

	total   int
	cloned  int
	cloning int
}

func (r *repositoryStatsResolver) GitDirBytes() BigInt {
	return BigInt{Int: int64(r.gitDirBytes)}
}

func (r *repositoryStatsResolver) IndexedLinesCount() BigInt {
	return BigInt{Int: int64(r.indexedLinesCount)}
}

func (r *repositoryStatsResolver) Total() int32   { return int32(r.total) }
func (r *repositoryStatsResolver) Cloned() int32  { return int32(r.cloned) }
func (r *repositoryStatsResolver) Cloning() int32 { return int32(r.cloning) }

func (r *schemaResolver) RepositoryStats(ctx context.Context) (*repositoryStatsResolver, error) {
	// 🚨 SECURITY: Only site admins may query repository statistics for the site.
	db := r.db
	if err := backend.CheckCurrentUserIsSiteAdmin(ctx, db); err != nil {
		return nil, err
	}

	stats, err := usagestats.GetRepositories(ctx, db)
	if err != nil {
		return nil, err
	}

	repoStats, err := r.db.RepoStatisticsStore().GetRepoStatistics(ctx)
	if err != nil {
		return nil, err
	}

	var (
		cloned  int
		cloning int
	)
	gitserverStats, err := r.db.RepoStatisticsStore().GetGitserverReposStatistics(ctx)
	if err != nil {
		return nil, err
	}
	for _, gs := range gitserverStats {
		cloned += gs.Cloned
		cloning += gs.Cloning
	}

	return &repositoryStatsResolver{
		gitDirBytes:       stats.GitDirBytes,
		indexedLinesCount: stats.DefaultBranchNewLinesCount + stats.OtherBranchesNewLinesCount,
		total:             repoStats.Total,
		cloned:            cloned,
		cloning:           cloning,
	}, nil
}
