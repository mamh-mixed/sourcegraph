package dependencies

import (
	"time"

	"github.com/sourcegraph/sourcegraph/internal/codeintel/autoindexing/shared"
	"github.com/sourcegraph/sourcegraph/internal/workerutil"
)

type AutoIndexingServiceBackgroundJobs interface {
	NewDependencySyncScheduler(interval time.Duration) *workerutil.Worker[shared.DependencySyncingJob]
	NewDependencyIndexingScheduler(interval time.Duration, numHandlers int) *workerutil.Worker[shared.DependencyIndexingJob]
}
