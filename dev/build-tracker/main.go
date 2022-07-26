package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	_ "log"
	"net/http"
	"os"
	"time"

	"github.com/sourcegraph/log"
)

var ErrRequestBody = fmt.Errorf("failed to read body from request")
var ErrJSONUnmarshall = fmt.Errorf("failed to unmarshal body")
var ErrInvalidToken = fmt.Errorf("buildkite token is invalid")
var ErrInvalidHeader = fmt.Errorf("Header of request is invalid")
var ErrUnwantedEvent = fmt.Errorf("Unwanted event received")

const DEFAULT_CHANNEL = "#william-buildchecker-webhook-test"

type BuildTrackingServer struct {
	logger       log.Logger
	store        *BuildStore
	bkToken      string
	notifyClient *NotificationClient
}

type token struct {
	Buildkite string
	Slack     string
	Github    string
}

func strp(v *string) string {
	if v == nil {
		return ""
	}

	return *v
}

func intp(v *int, defaultValue int) int {
	if v == nil {
		return defaultValue
	}

	return *v
}

func mustEnvVar(name string) string {
	value, exists := os.LookupEnv(name)
	if !exists {
		panic(fmt.Sprintf("%s not found in environment", name))
	}

	return value
}

func MustGetTokenFromEnv() token {
	return token{
		Buildkite: mustEnvVar("BUILDKITE_WEBHOOK_TOKEN"),
		Slack:     mustEnvVar("SLACK_TOKEN"),
		Github:    mustEnvVar("GITHUB_TOKEN"),
	}
}

func NewBuildTrackingServer(logger log.Logger, tokens token, channel string) (*BuildTrackingServer, error) {
	serverLog := logger.Scoped("server", "Build Tracking Server which tracks events received from Buildkite and sends notifications on failures")
	return &BuildTrackingServer{
		logger:       serverLog,
		store:        NewBuildStore(serverLog),
		bkToken:      tokens.Buildkite,
		notifyClient: NewNotificationClient(serverLog, tokens.Slack, tokens.Github, channel),
	}, nil
}

func (s *BuildTrackingServer) processBuildkiteRequest(req *http.Request, token string) (*BuildEvent, error) {
	h, ok := req.Header["X-Buildkite-Token"]
	if !ok || len(h) == 0 {
		return nil, ErrInvalidToken
	} else if h[0] != token {
		return nil, ErrInvalidToken
	}

	h, ok = req.Header["X-Buildkite-Event"]
	if !ok || len(h) == 0 {
		return nil, ErrInvalidHeader
	}

	eventName := h[0]
	s.logger.Debug("receied event", log.String("eventName", eventName))

	var event BuildEvent
	err := readBody(s.logger, req, &event)
	if errors.Is(err, ErrRequestBody) {
		s.logger.Error("failed to read body of request", log.Error(err))
		return nil, ErrRequestBody
	} else if errors.Is(err, ErrJSONUnmarshall) {
		s.logger.Error("failed to unmarshall body of request", log.Error(err))
		return nil, ErrJSONUnmarshall
	}

	return &event, nil
}

// handleEvent handles an event received from the http listener. A event is valid when:
// - Has the correct headers from Buildkite
// - On of the following events
//   * job.finished
//   * build.finished
// - Has valid JSON
// Note that if we received an unwanted event ie. the event is not "job.finished" or "build.finished" we respond with a 200 OK regardless.
// Once all the conditions are met, the event is processed in a go routine with `processEvent`
func (s *BuildTrackingServer) handleEvent(w http.ResponseWriter, req *http.Request) {
	event, err := s.processBuildkiteRequest(req, s.bkToken)

	switch err {
	case ErrInvalidToken:
		w.WriteHeader(http.StatusUnauthorized)
		return
	case ErrInvalidHeader:
		w.WriteHeader(http.StatusBadRequest)
		return
	case ErrUnwantedEvent:
		w.WriteHeader(http.StatusOK)
		return
	case ErrRequestBody:
		w.WriteHeader(http.StatusBadRequest)
		return
	case ErrJSONUnmarshall:
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	go s.processEvent(event)
	w.WriteHeader(http.StatusOK)
}

func (s *BuildTrackingServer) handleHealthz(w http.ResponseWriter, req *http.Request) {
	// do our super exhaustive check
	w.WriteHeader(http.StatusOK)
}

func readBody[T any](logger log.Logger, req *http.Request, target T) error {
	data, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		logger.Error("failed to read request body", log.Error(err))
		return ErrRequestBody
	}

	logger.Debug("event data", log.String("data", string(data)))
	err = json.Unmarshal(data, &target)
	if err != nil {
		logger.Error("failed to unmarshall request body", log.Error(err))
		return ErrJSONUnmarshall
	}

	return nil
}

// notifyIfFailed sends a notification over slack if the provided build has failed. If the build is successful not notifcation is sent
func (s *BuildTrackingServer) notifyIfFailed(build *Build) error {
	if build.hasFailed() {
		s.logger.Info("detected failed build - sending notification", log.Int("buildNumber", intp(build.Number, 0)))
		return s.notifyClient.sendNotification(build)
	}

	s.logger.Info("build has not failed", log.Int("buildNumber", intp(build.Number, 0)))
	return nil
}

func (s *BuildTrackingServer) startOldBuildCleaner(every, window time.Duration) func() {
	ticker := time.NewTicker(every)
	done := make(chan interface{})

	// We could technically remove  the builds immediately after we've sent a notification for or it, or the build has passed.
	// But we keep builds a little longer and prediodically clean them out so that we can in future allow possibly querying
	// of builds and other use cases, like retrying a build etc.
	go func() {
		for {
			select {
			case <-ticker.C:
				{
					oldBuilds := make([]int, 0)
					now := time.Now()
					for _, b := range s.store.AllFinishedBuilds() {
						finishedAt := *b.FinishedAt
						delta := now.Sub(finishedAt.Time)
						if delta >= window {
							s.logger.Debug("build past age window", log.Int("buildNumber", *b.Number), log.Time("FinishedAt", finishedAt.Time), log.Duration("window", window))
							oldBuilds = append(oldBuilds, *b.Number)
						}
					}
					s.logger.Info("deleting old builds", log.Int("oldBuildCount", len(oldBuilds)))
					s.store.DelByBuildNumber(oldBuilds...)
				}
			case <-done:
				{
					ticker.Stop()
					return
				}
			}
		}
	}()

	return func() { done <- nil }
}

// processEvent processes a BuildEvent received from Buildkite. If the event is for a `build.finished` event we get the
// full build which includes all recorded jobs for the build and send a notification.
// processEvent delegates the decision to actually send a notifcation
func (s *BuildTrackingServer) processEvent(event *BuildEvent) {
	s.logger.Info("processing event", log.String("eventName", event.Name), log.Int("buildNumber", event.buildNumber()), log.String("jobName", event.jobName()))
	s.store.Add(event)
	if event.isBuildFinished() {
		build := s.store.GetByBuildNumber(event.buildNumber())
		if err := s.notifyIfFailed(build); err != nil {
			s.logger.Error("failed to send notification for build", log.Int("buildNumber", event.buildNumber()), log.Error(err))
		}
	}
}

// Server starts the http server and listens for buildkite build events to be sent on the route "/buildkite"
func (s *BuildTrackingServer) Serve() error {
	http.HandleFunc("/buildkite", s.handleEvent)
	http.HandleFunc("/healthz", s.handleHealthz)
	s.logger.Info("listening on :8080")
	return http.ListenAndServe(":8080", nil)
}

func main() {
	sync := log.Init(log.Resource{
		Name:      "BuildTracker",
		Namespace: "CI",
	})
	defer sync.Sync()

	logger := log.Scoped("BuildTracker", "main entrypoint for Build Tracking Server")

	channel := os.Getenv("SLACK_CHANNEL")
	if channel == "" {
		channel = DEFAULT_CHANNEL
	}

	server, err := NewBuildTrackingServer(logger, MustGetTokenFromEnv(), channel)
	if err != nil {
		logger.Fatal("failed to create BuildTracking server", log.Error(err))
	}

	stopFn := server.startOldBuildCleaner(5*time.Minute, 24*time.Hour)
	defer stopFn()
	if err := server.Serve(); err != nil {
		logger.Fatal("server exited with error", log.Error(err))
	}
}
