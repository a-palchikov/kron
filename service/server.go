package kron

import (
	"service/api"

	"golang.org/x/net/context"
)

type (
	server struct{}

	/*
		// Node in a cluster
		Node interface {
			// Schedule schedules a process (on a remote node).
			// name names the new cron job.
			// command defines the command to execute, optionally prefixed with a node name,
			// and a list of arguments.
			// cron is a scheduled time as a cron expression.
			// tags name the node tags this job is to run on.
			Schedule(name string, command string, cron string, tags []string) (Job, error)
			// Serf (https://www.serfdom.io/) for cluster membership/failure detection/orchestration.
			Jobs() []Job
		}

			Job interface {
				Stop() error
				Pause() error
				Status() (JobStatus, error)
			}

			JobStatus struct {
				successCount int
				errorCount   int
				tags         []string
				// lastRun	time.Time
			}

			Process struct {
				Stdin  io.Reader
				Stdout io.Reader
				Stderr io.Reader

				// TODO: API
			}
	*/

)

func (s *server) Schedule(ctx context.Context, job *api.Job) (*api.ScheduleResponse, error) {
	// TODO: iterate all available nodes matching those with job.tags
	// dispatch job to node(s) given job's disposition.
	return &api.ScheduleResponse{}, nil
}
