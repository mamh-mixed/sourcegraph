package autoindexing

import (
	"context"
	"os"

	"github.com/inconshreveable/log15"

	"github.com/sourcegraph/sourcegraph/internal/api"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/stores/dbstore"
	"github.com/sourcegraph/sourcegraph/internal/observation"
	"github.com/sourcegraph/sourcegraph/lib/codeintel/autoindex/config"
	"github.com/sourcegraph/sourcegraph/lib/errors"
)

// queueIndexForRepositoryAndCommit determines a set of index jobs to enqueue for the given repository and commit.
//
// If the force flag is false, then the presence of an upload or index record for this given repository and commit
// will cause this method to no-op. Note that this is NOT a guarantee that there will never be any duplicate records
// when the flag is false.
func (s *Service) QueueIndexForRepositoryAndCommit(ctx context.Context, repositoryID int, commit, configuration string, force bool, trace observation.TraceLogger) ([]dbstore.Index, error) {
	if !force {
		isQueued, err := s.dbStore.IsQueued(ctx, repositoryID, commit)
		if err != nil {
			return nil, errors.Wrap(err, "dbstore.IsQueued")
		}
		if isQueued {
			return nil, nil
		}
	}

	indexes, err := s.getIndexRecords(ctx, repositoryID, commit, configuration)
	if err != nil {
		return nil, err
	}
	if len(indexes) == 0 {
		return nil, nil
	}
	// trace.Log(log.Int("numIndexes", len(indexes)))

	return s.dbStore.InsertIndexes(ctx, indexes)
}

var overrideScript = os.Getenv("SRC_CODEINTEL_INFERENCE_OVERRIDE_SCRIPT")

// InferIndexJobsFromRepositoryStructure collects the result of InferIndexJobs over all registered recognizers.
func (s *Service) InferIndexJobsFromRepositoryStructure(ctx context.Context, repositoryID int, commit string) ([]config.IndexJob, error) {
	// if err := s.gitserverLimiter.Wait(ctx); err != nil {
	// 	return nil, err
	// }

	repoName, err := s.dbStore.RepoName(ctx, repositoryID)
	if err != nil {
		return nil, err
	}

	indexes, err := s.inferenceService.InferIndexJobs(ctx, api.RepoName(repoName), commit, overrideScript)
	if err != nil {
		return nil, err
	}

	if len(indexes) > s.maximumIndexJobsPerInferredConfiguration {
		log15.Info("Too many inferred roots. Scheduling no index jobs for repository.", "repository_id", repositoryID)
		return nil, nil
	}

	return indexes, nil
}

// inferIndexJobsFromRepositoryStructure collects the result of  InferIndexJobHints over all registered recognizers.
func (s *Service) InferIndexJobHintsFromRepositoryStructure(ctx context.Context, repositoryID int, commit string) ([]config.IndexJobHint, error) {
	// if err := s.gitserverLimiter.Wait(ctx); err != nil {
	// 	return nil, err
	// }

	repoName, err := s.dbStore.RepoName(ctx, repositoryID)
	if err != nil {
		return nil, err
	}

	indexes, err := s.inferenceService.InferIndexJobHints(ctx, api.RepoName(repoName), commit, overrideScript)
	if err != nil {
		return nil, err
	}

	return indexes, nil
}
