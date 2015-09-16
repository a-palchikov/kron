package kron

import (
	pb "service/proto"

	"golang.org/x/net/context"
)

type (
	Server struct {
		Infrastructure
		config *Config
		id     NodeId
		//scheduler.Scheduler	// actual cron management
	}

	/*
		Node interface {
			IsActive() bool
			// FIXME: more queries
		}
	*/

	NodeId uint64

	Infrastructure interface {
		// Persist current schedule
		// This is done every time a job is added/removed to keep all nodes in sync
		SetSchedule([]*pb.Job) error
		// ScheduleJob publishes scheduled job on the schedule topic so that matching nodes
		// can pick it up
		ScheduleJob(*pb.Job) error
		// FIXME: this is a bit hacky as there are responsibilities of a cluster membership
		// framework
		// Add a new node to cluster
		AddNode(host string) (NodeId, error)
		// Nodes() ([]Node, error)
	}
)

func New(config *Config, infrastructure Infrastructure) (*Server, error) {
	server := Server{
		config:         config,
		Infrastructure: infrastructure,
	}

	if !config.master { // perform for server too?
		host, err := os.Hostname()
		if err != nil {
			return nil, err
		}

		if id, err := infrastructure.AddNode(host); err != nil {
			return nil, err
		}
	}

	// go server.Scheduler.Loop(server.Infrastructure)	// Infrastructure arg?

	return &server, nil
}

func (s *Server) Stop() error {
	// server.Scheduler.Stop()
	return nil
}

func (s *server) Schedule(ctx context.Context, req *pb.ScheduleRequest) (*pb.ScheduleResponse, error) {
	// TODO
	return &pb.ScheduleResponse{}, nil
}
