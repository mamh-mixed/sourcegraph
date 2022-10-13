package zoektrepos

import (
	"context"
	"time"

	"github.com/sourcegraph/log"

	"github.com/sourcegraph/sourcegraph/cmd/frontend/backend"
	"github.com/sourcegraph/sourcegraph/cmd/worker/job"
	workerdb "github.com/sourcegraph/sourcegraph/cmd/worker/shared/init/db"
	"github.com/sourcegraph/sourcegraph/internal/actor"
	"github.com/sourcegraph/sourcegraph/internal/conf"
	"github.com/sourcegraph/sourcegraph/internal/database"
	"github.com/sourcegraph/sourcegraph/internal/env"
	"github.com/sourcegraph/sourcegraph/internal/gitserver"
	"github.com/sourcegraph/sourcegraph/internal/goroutine"
	"github.com/sourcegraph/sourcegraph/internal/search"
)

type updater struct{}

var _ job.Job = &updater{}

func NewUpdater() job.Job {
	return &updater{}
}

func (j *updater) Description() string {
	return ""
}

func (j *updater) Config() []env.Config {
	return nil
}

func (j *updater) Routines(startupCtx context.Context, logger log.Logger) ([]goroutine.BackgroundRoutine, error) {
	db, err := workerdb.InitDBWithLogger(logger)
	if err != nil {
		return nil, err
	}

	gitserverclient := gitserver.NewClient(db)

	return []goroutine.BackgroundRoutine{
		// TODO: Fix the time here
		goroutine.NewPeriodicGoroutine(context.Background(), 3*time.Second, &handler{
			db:              db,
			logger:          logger,
			gitserverClient: gitserverclient,
		}),
	}, nil
}

type handler struct {
	db              database.DB
	logger          log.Logger
	gitserverClient gitserver.Client
}

var _ goroutine.Handler = &handler{}
var _ goroutine.ErrorHandler = &handler{}

func (h *handler) Handle(ctx context.Context) error {
	if !conf.SearchIndexEnabled() {
		return nil
	}

	// IMPORTANTE!!
	ctx = actor.WithInternalActor(ctx)

	indexable, err := backend.NewRepos(h.logger, h.db, h.gitserverClient).ListIndexable(ctx)
	if err != nil {
		return err
	}

	indexed, err := search.ListAllIndexed(ctx)
	if err != nil {
		return err
	}

	return h.db.ZoektRepos().UpsertIndexable(ctx, indexable, indexed.Minimal)
}

func (h *handler) HandleError(err error) {
	h.logger.Error("error updating zoekt repos", log.Error(err))
}
