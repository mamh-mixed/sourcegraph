package cleanup

import (
	"time"

	"github.com/sourcegraph/sourcegraph/internal/codeintel/autoindexing/shared"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/shared/types"
	"github.com/sourcegraph/sourcegraph/internal/goroutine"
	"github.com/sourcegraph/sourcegraph/internal/workerutil/dbworker"
)

type AutoIndexingServiceBackgroundJobs interface {
	NewJanitor(
		interval time.Duration,
		minimumTimeSinceLastCheck time.Duration,
		commitResolverBatchSize int,
		commitResolverMaximumCommitLag time.Duration,
	) goroutine.BackgroundRoutine

	NewIndexResetter(interval time.Duration) *dbworker.Resetter[types.Index]
	NewDependencyIndexResetter(interval time.Duration) *dbworker.Resetter[shared.DependencyIndexingJob]
}
