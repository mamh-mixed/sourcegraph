package scheduler

import (
	"context"

	"github.com/sourcegraph/sourcegraph/internal/goroutine"
)

func NewScheduler(dbStore DBStore, policyMatcher PolicyMatcher, indexEnqueuer IndexEnqueuer) goroutine.BackgroundRoutine {
	return goroutine.NewPeriodicGoroutine(context.Background(), ConfigInst.Interval, &scheduler{
		dbStore:       dbStore,
		policyMatcher: policyMatcher,
		indexEnqueuer: indexEnqueuer,
	})
}
