package codeintel

import (
	"context"

	"github.com/sourcegraph/sourcegraph/cmd/worker/job"
	"github.com/sourcegraph/sourcegraph/cmd/worker/shared/init/codeintel"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/autoindexing/background/scheduler"
	policies "github.com/sourcegraph/sourcegraph/internal/codeintel/policies/enterprise"
	"github.com/sourcegraph/sourcegraph/internal/env"
	"github.com/sourcegraph/sourcegraph/internal/goroutine"
	"github.com/sourcegraph/sourcegraph/lib/log"
)

type autoindexingScheduler struct{}

func NewAutoindexingSchedulerJob() job.Job {
	return &autoindexingScheduler{}
}

func (j *autoindexingScheduler) Description() string {
	return ""
}

func (j *autoindexingScheduler) Config() []env.Config {
	return []env.Config{
		scheduler.ConfigInst,
	}
}

func (j *autoindexingScheduler) Routines(ctx context.Context, logger log.Logger) ([]goroutine.BackgroundRoutine, error) {
	dbStore, err := codeintel.InitDBStore()
	if err != nil {
		return nil, err
	}

	gitserverClient, err := codeintel.InitGitserverClient()
	if err != nil {
		return nil, err
	}

	policyMatcher := policies.NewMatcher(gitserverClient, policies.IndexingExtractor, false, true)

	return []goroutine.BackgroundRoutine{
		// TODO - index enqueuer should be the service
		scheduler.NewScheduler(dbStore, policyMatcher, nil),
	}, nil
}
