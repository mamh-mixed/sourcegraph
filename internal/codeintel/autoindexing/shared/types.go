package shared

import (
	"time"

	"github.com/sourcegraph/sourcegraph/internal/codeintel/types"
)

type IndexJob struct {
	Indexer string
}

type SourcedCommits struct {
	RepositoryID   int
	RepositoryName string
	Commits        []string
}

// Index is a subset of the lsif_indexes table and stores both processed and unprocessed
// records.
// type Index struct {
// 	ID                 int                 `json:"id"`
// 	Commit             string              `json:"commit"`
// 	QueuedAt           time.Time           `json:"queuedAt"`
// 	State              string              `json:"state"`
// 	FailureMessage     *string             `json:"failureMessage"`
// 	StartedAt          *time.Time          `json:"startedAt"`
// 	FinishedAt         *time.Time          `json:"finishedAt"`
// 	ProcessAfter       *time.Time          `json:"processAfter"`
// 	NumResets          int                 `json:"numResets"`
// 	NumFailures        int                 `json:"numFailures"`
// 	RepositoryID       int                 `json:"repositoryId"`
// 	LocalSteps         []string            `json:"local_steps"`
// 	RepositoryName     string              `json:"repositoryName"`
// 	DockerSteps        []DockerStep        `json:"docker_steps"`
// 	Root               string              `json:"root"`
// 	Indexer            string              `json:"indexer"`
// 	IndexerArgs        []string            `json:"indexer_args"` // TODO - convert this to `IndexCommand string`
// 	Outfile            string              `json:"outfile"`
// 	ExecutionLogs      []ExecutionLogEntry `json:"execution_logs"`
// 	Rank               *int                `json:"placeInQueue"`
// 	AssociatedUploadID *int                `json:"associatedUpload"`
// }

// type DockerStep struct {
// 	Root     string   `json:"root"`
// 	Image    string   `json:"image"`
// 	Commands []string `json:"commands"`
// }

// func (s *DockerStep) Scan(value any) error {
// 	b, ok := value.([]byte)
// 	if !ok {
// 		return errors.Errorf("value is not []byte: %T", value)
// 	}

// 	return json.Unmarshal(b, &s)
// }

// func (s DockerStep) Value() (driver.Value, error) {
// 	return json.Marshal(s)
// }

// // ExecutionLogEntry represents a command run by the executor.
// type ExecutionLogEntry struct {
// 	Key        string    `json:"key"`
// 	Command    []string  `json:"command"`
// 	StartTime  time.Time `json:"startTime"`
// 	ExitCode   *int      `json:"exitCode,omitempty"`
// 	Out        string    `json:"out,omitempty"`
// 	DurationMs *int      `json:"durationMs,omitempty"`
// }

// func (e *ExecutionLogEntry) Scan(value any) error {
// 	b, ok := value.([]byte)
// 	if !ok {
// 		return errors.Errorf("value is not []byte: %T", value)
// 	}

// 	return json.Unmarshal(b, &e)
// }

// func (e ExecutionLogEntry) Value() (driver.Value, error) {
// 	return json.Marshal(e)
// }

// func ExecutionLogEntries(raw []workerutil.ExecutionLogEntry) (entries []ExecutionLogEntry) {
// 	for _, entry := range raw {
// 		entries = append(entries, ExecutionLogEntry(entry))
// 	}

// 	return entries
// }

// IndexConfiguration stores the index configuration for a repository.
type IndexConfiguration struct {
	ID           int    `json:"id"`
	RepositoryID int    `json:"repository_id"`
	Data         []byte `json:"data"`
}

type GetIndexesOptions struct {
	RepositoryID int
	State        string
	Term         string
	Limit        int
	Offset       int
}

type IndexesWithRepositoryNamespace struct {
	Root    string
	Indexer string
	Indexes []types.Index
}

// UploadLocation is a path and range pair from within a particular upload. The target commit
// denotes the target commit for which the location was set (the originally requested commit).
type UploadLocation struct {
	Dump         Dump
	Path         string
	TargetCommit string
	TargetRange  Range
}

// Dump is a subset of the lsif_uploads table (queried via the lsif_dumps_with_repository_name view)
// and stores only processed records.
type Dump struct {
	ID                int        `json:"id"`
	Commit            string     `json:"commit"`
	Root              string     `json:"root"`
	VisibleAtTip      bool       `json:"visibleAtTip"`
	UploadedAt        time.Time  `json:"uploadedAt"`
	State             string     `json:"state"`
	FailureMessage    *string    `json:"failureMessage"`
	StartedAt         *time.Time `json:"startedAt"`
	FinishedAt        *time.Time `json:"finishedAt"`
	ProcessAfter      *time.Time `json:"processAfter"`
	NumResets         int        `json:"numResets"`
	NumFailures       int        `json:"numFailures"`
	RepositoryID      int        `json:"repositoryId"`
	RepositoryName    string     `json:"repositoryName"`
	Indexer           string     `json:"indexer"`
	IndexerVersion    string     `json:"indexerVersion"`
	AssociatedIndexID *int       `json:"associatedIndex"`
}

// Range is an inclusive bounds within a file.
type Range struct {
	Start Position
	End   Position
}

// Position is a unique position within a file.
type Position struct {
	Line      int
	Character int
}
