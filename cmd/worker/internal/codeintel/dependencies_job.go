package codeintel

import (
	"context"

	"github.com/sourcegraph/sourcegraph/cmd/worker/job"
	"github.com/sourcegraph/sourcegraph/cmd/worker/shared/init/codeintel"
	workerdb "github.com/sourcegraph/sourcegraph/cmd/worker/shared/init/db"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/autoindexing"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/dependencies/background/cratesyncer"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/dependencies/background/indexer"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/dependencies/background/resolver"
	livedependencies "github.com/sourcegraph/sourcegraph/internal/codeintel/dependencies/live"
	policies "github.com/sourcegraph/sourcegraph/internal/codeintel/policies/enterprise"
	"github.com/sourcegraph/sourcegraph/internal/database"
	"github.com/sourcegraph/sourcegraph/internal/env"
	"github.com/sourcegraph/sourcegraph/internal/goroutine"
	"github.com/sourcegraph/sourcegraph/lib/log"
)

type dependenciesJob struct{}

func NewDependenciesJob() job.Job {
	return &dependenciesJob{}
}

func (j *dependenciesJob) Description() string {
	return ""
}

func (j *dependenciesJob) Config() []env.Config {
	return []env.Config{
		indexer.ConfigInst,
		resolver.ConfigInst,
	}
}

func (j *dependenciesJob) Routines(ctx context.Context, logger log.Logger) ([]goroutine.BackgroundRoutine, error) {
	db, err := workerdb.Init()
	if err != nil {
		return nil, err
	}

	dbStore, err := codeintel.InitDBStore()
	if err != nil {
		return nil, err
	}

	gitserverClient, err := codeintel.InitGitserverClient()
	if err != nil {
		return nil, err
	}

	// gitSvc := live.NewGitService(database.NewDB(db))

	policyMatcher := policies.NewMatcher(gitserverClient, policies.IndexingExtractor, false, true)

	repoUpdaterClient := codeintel.InitRepoUpdaterClient()
	enqueuerDBStoreShim := &autoindexing.DBStoreShim{Store: dbStore}
	indexEnqueuer := autoindexing.GetService(database.NewDB(db), enqueuerDBStoreShim, gitserverClient, repoUpdaterClient)

	return []goroutine.BackgroundRoutine{
		indexer.NewIndexer(database.NewDB(db), livedependencies.NewSyncer(), dbStore, policyMatcher),
		resolver.NewResolver(database.NewDB(db), livedependencies.NewSyncer()),
		cratesyncer.NewCratesSyncer(database.NewDB(db), indexEnqueuer),
	}, nil
}
