package codeintel

import (
	"context"

	"github.com/sourcegraph/sourcegraph/enterprise/cmd/frontend/internal/executorqueue/handler"
	apiclient "github.com/sourcegraph/sourcegraph/enterprise/internal/executor"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/autoindexing"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/shared/types"
	"github.com/sourcegraph/sourcegraph/internal/observation"
)

func QueueOptions(autoIndexingSvc *autoindexing.Service, accessToken func() string, observationContext *observation.Context) handler.QueueOptions[types.Index] {
	recordTransformer := func(ctx context.Context, record types.Index) (apiclient.Job, error) {
		return transformRecord(record, accessToken())
	}

	return handler.QueueOptions[types.Index]{
		Name:              "codeintel",
		Store:             autoindexing.GetWorkerutilStore(autoIndexingSvc),
		RecordTransformer: recordTransformer,
	}
}
