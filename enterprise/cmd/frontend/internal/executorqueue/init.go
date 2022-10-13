package executorqueue

import (
	"context"

	"github.com/sourcegraph/sourcegraph/internal/conf/conftypes"
	"github.com/sourcegraph/sourcegraph/internal/database"
	metricsstore "github.com/sourcegraph/sourcegraph/internal/metrics/store"
	"github.com/sourcegraph/sourcegraph/internal/observation"
	executorDB "github.com/sourcegraph/sourcegraph/internal/services/executors/store/db"

	"github.com/sourcegraph/sourcegraph/cmd/frontend/enterprise"
	"github.com/sourcegraph/sourcegraph/enterprise/cmd/frontend/internal/executorqueue/handler"
	"github.com/sourcegraph/sourcegraph/enterprise/cmd/frontend/internal/executorqueue/queues/batches"
	codeintelqueue "github.com/sourcegraph/sourcegraph/enterprise/cmd/frontend/internal/executorqueue/queues/codeintel"
)

// Init initializes the executor endpoints required for use with the executor service.
func Init(
	ctx context.Context,
	db database.DB,
	conf conftypes.UnifiedWatchable,
	enterpriseServices *enterprise.Services,
	observationContext *observation.Context,
) error {
	autoIndexingService := enterpriseServices.CodeIntelAutoIndexingService
	codeintelUploadHandler := enterpriseServices.NewCodeIntelUploadHandler(false)
	batchesWorkspaceFileGetHandler := enterpriseServices.BatchesChangesFileGetHandler
	batchesWorkspaceFileExistsHandler := enterpriseServices.BatchesChangesFileGetHandler
	accessToken := func() string { return conf.SiteConfig().ExecutorsAccessToken }

	metricsStore := metricsstore.NewDistributedStore("executors:")
	executorStore := executorDB.New(db)

	// Register queues. If this set changes, be sure to also update the list of valid
	// queue names in ./metrics/queue_allocation.go, and register a metrics exporter
	// in the worker.
	//
	// Note: In order register a new queue type please change the validate() check code in enterprise/cmd/executor/config.go
	codeintelHandler := handler.NewHandler(executorStore, metricsStore, codeintelqueue.QueueOptions(autoIndexingService, accessToken, observationContext))
	batchesHandler := handler.NewHandler(executorStore, metricsStore, batches.QueueOptions(db, accessToken, observationContext))
	queueOptions := map[string]handler.PubHandler{
		codeintelHandler.Name: codeintelHandler,
		batchesHandler.Name:   batchesHandler,
	}

	queueHandler, err := newExecutorQueueHandler(
		db,
		queueOptions,
		accessToken,
		codeintelUploadHandler,
		batchesWorkspaceFileGetHandler,
		batchesWorkspaceFileExistsHandler,
	)
	if err != nil {
		return err
	}

	enterpriseServices.NewExecutorProxyHandler = queueHandler
	return nil
}
