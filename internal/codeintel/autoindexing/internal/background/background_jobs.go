package background

import (
	"time"

	"github.com/derision-test/glock"

	"github.com/sourcegraph/log"

	"github.com/sourcegraph/sourcegraph/internal/codeintel/autoindexing/internal/store"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/autoindexing/shared"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/shared/types"
	"github.com/sourcegraph/sourcegraph/internal/database"
	"github.com/sourcegraph/sourcegraph/internal/goroutine"
	"github.com/sourcegraph/sourcegraph/internal/observation"
	"github.com/sourcegraph/sourcegraph/internal/workerutil"
	"github.com/sourcegraph/sourcegraph/internal/workerutil/dbworker"
	dbworkerstore "github.com/sourcegraph/sourcegraph/internal/workerutil/dbworker/store"
)

type BackgroundJob interface {
	NewDependencyIndexingScheduler(pollInterval time.Duration, numHandlers int) *workerutil.Worker[shared.DependencyIndexingJob]
	NewDependencySyncScheduler(pollInterval time.Duration) *workerutil.Worker[shared.DependencySyncingJob]
	NewDependencyIndexResetter(interval time.Duration) *dbworker.Resetter[shared.DependencyIndexingJob]
	NewIndexResetter(interval time.Duration) *dbworker.Resetter[types.Index]
	NewOnDemandScheduler(interval time.Duration, batchSize int) goroutine.BackgroundRoutine
	NewScheduler(interval time.Duration, repositoryProcessDelay time.Duration, repositoryBatchSize int, policyBatchSize int) goroutine.BackgroundRoutine
	NewJanitor(
		interval time.Duration,
		minimumTimeSinceLastCheck time.Duration,
		commitResolverBatchSize int,
		commitResolverMaximumCommitLag time.Duration,
	) goroutine.BackgroundRoutine

	SetService(service AutoIndexingService)
	WorkerutilStore() dbworkerstore.Store[types.Index]
	DependencySyncStore() dbworkerstore.Store[shared.DependencySyncingJob]
	DependencyIndexingStore() dbworkerstore.Store[shared.DependencyIndexingJob]
}

type backgroundJob struct {
	uploadSvc       UploadService
	depsSvc         DependenciesService
	policiesSvc     PoliciesService
	autoindexingSvc AutoIndexingService

	policyMatcher   PolicyMatcher
	repoUpdater     RepoUpdaterClient
	gitserverClient GitserverClient

	store                   store.Store
	repoStore               ReposStore
	workerutilStore         dbworkerstore.Store[types.Index]
	gitserverRepoStore      GitserverRepoStore
	dependencySyncStore     dbworkerstore.Store[shared.DependencySyncingJob]
	externalServiceStore    ExternalServiceStore
	dependencyIndexingStore dbworkerstore.Store[shared.DependencyIndexingJob]

	operations *operations
	clock      glock.Clock
	logger     log.Logger

	metrics                *resetterMetrics
	janitorMetrics         *janitorMetrics
	depencencySyncMetrics  workerutil.WorkerMetrics
	depencencyIndexMetrics workerutil.WorkerMetrics
}

func New(
	db database.DB,
	store store.Store,
	uploadSvc UploadService,
	depsSvc DependenciesService,
	policiesSvc PoliciesService,
	policyMatcher PolicyMatcher,
	gitserverClient GitserverClient,
	repoUpdater RepoUpdaterClient,
	observationContext *observation.Context,
) BackgroundJob {
	repoStore := db.Repos()
	gitserverRepoStore := db.GitserverRepos()
	externalServiceStore := db.ExternalServices()
	workerutilStore := dbworkerstore.NewWithMetrics(db.Handle(), indexWorkerStoreOptions, observationContext)
	dependencySyncStore := dbworkerstore.NewWithMetrics(db.Handle(), dependencySyncingJobWorkerStoreOptions, observationContext)
	dependencyIndexingStore := dbworkerstore.NewWithMetrics(db.Handle(), dependencyIndexingJobWorkerStoreOptions, observationContext)

	return &backgroundJob{
		uploadSvc:   uploadSvc,
		depsSvc:     depsSvc,
		policiesSvc: policiesSvc,

		policyMatcher:   policyMatcher,
		repoUpdater:     repoUpdater,
		gitserverClient: gitserverClient,

		store:                   store,
		repoStore:               repoStore,
		workerutilStore:         workerutilStore,
		gitserverRepoStore:      gitserverRepoStore,
		dependencySyncStore:     dependencySyncStore,
		externalServiceStore:    externalServiceStore,
		dependencyIndexingStore: dependencyIndexingStore,

		operations: newOperations(observationContext),
		clock:      glock.NewRealClock(),
		logger:     observationContext.Logger,

		metrics:                newResetterMetrics(observationContext),
		janitorMetrics:         newJanitorMetrics(observationContext),
		depencencySyncMetrics:  workerutil.NewMetrics(observationContext, "codeintel_dependency_index_processor"),
		depencencyIndexMetrics: workerutil.NewMetrics(observationContext, "codeintel_dependency_index_queueing"),
	}
}

func (b *backgroundJob) SetService(service AutoIndexingService) {
	b.autoindexingSvc = service
}

func (b backgroundJob) WorkerutilStore() dbworkerstore.Store[types.Index] { return b.workerutilStore }
func (b backgroundJob) DependencySyncStore() dbworkerstore.Store[shared.DependencySyncingJob] {
	return b.dependencySyncStore
}

func (b backgroundJob) DependencyIndexingStore() dbworkerstore.Store[shared.DependencyIndexingJob] {
	return b.dependencyIndexingStore
}
