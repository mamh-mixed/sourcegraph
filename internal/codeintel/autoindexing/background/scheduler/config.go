package scheduler

import (
	"time"

	"github.com/sourcegraph/sourcegraph/internal/env"
)

type config struct {
	env.BaseConfig

	Interval                       time.Duration
	RepositoryMinimumCheckInterval time.Duration
	RepositoryBatchSize            int
	PolicyBatchSize                int
}

var ConfigInst = &config{}

func (c *config) Load() {
	c.Interval = c.GetInterval("CODEINTEL_AUTOINDEXING_SCHEDULER_INTERVAL", "1s", "How frequently to run the autoindexer scheduling routine.")
	c.RepositoryMinimumCheckInterval = c.GetInterval("CODEINTEL_AUTOINDEXING_SCHEDULER_REPOSITORY_MINIMUM_CHECK_INTERVAL", "10m", "How frequently to re-index the same repository.")
	c.RepositoryBatchSize = c.GetInt("CODEINTEL_AUTOINDEXING_SCHEDULER_REPOSITORY_BATCH_SIZE", "100", "How many repositories to index at a time.")
	c.PolicyBatchSize = c.GetInt("CODEINTEL_AUTOINDEXING_SCHEDULER_POLICY_BATCH_SIZE", "100", "How many policies to load at a time.")
}
