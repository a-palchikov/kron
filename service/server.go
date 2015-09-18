package service

import (
	"fmt"
	"io"
	"net"
	"sync"

	pb "github.com/a-palchikov/kron/proto/servicepb"
	"github.com/a-palchikov/kron/scheduler"
	"google.golang.org/grpc"

	"golang.org/x/net/context"
)

type (
	JobId uint64

	Server struct {
		Cluster
		Store
		*Worker
		mu                  *sync.Mutex
		pending             map[uint64][]uint64 // jobId->[]nodeId
		scheduler.Scheduler                     // cron management
	}
)

func New(config *Config, store Store, cluster Cluster, tags []string) (*Server, error) {
	// FIXME: create only parts of the Worker that the server needs as well
	// the mode-sensitive parts are to be inited on-demand with the mode (master<->worker) change
	worker, err := NewWorker(config, store, cluster, tags)
	if err != nil {
		return nil, err
	}

	server := &Server{
		Worker:  worker,
		Store:   store,
		Cluster: cluster,
		pending: make(map[uint64][]uint64),
	}

	go server.Scheduler.Run()
	go server.loop()

	return server, nil
}

// Serve spawns serve loops for both API/feedback services.
func (s *Server) Serve() error {
	apiServer, apiListener, err := bindApi(s)
	if err != nil {
		return err
	}
	defer apiListener.Close()

	feedbackServer, feedbackListener, err := bindFeedback(s)
	if err != nil {
		return err
	}

	// FIXME: report serve errors from feedback service
	go func() {
		defer feedbackListener.Close()
		feedbackServer.Serve(feedbackListener)
	}()

	return apiServer.Serve(apiListener)
}

// bindApi creates an API server.
func bindApi(s *Server) (*grpc.Server, net.Listener, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.ApiPort))
	if err != nil {
		return nil, nil, err
	}

	server := grpc.NewServer()

	pb.RegisterSchedulerServiceServer(server, s)
	return server, listener, nil
}

// bindFeedback creates a feedback server.
func bindFeedback(s *Server) (*grpc.Server, net.Listener, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.FeedbackPort))
	if err != nil {
		return nil, nil, err
	}

	server := grpc.NewServer()

	pb.RegisterFeedbackServiceServer(server, s)
	return server, listener, nil
}

// SchedulerService
func (s *Server) Schedule(ctx context.Context, req *pb.ScheduleRequest) (*pb.ScheduleResponse, error) {
	// jobId := s.Scheduler.Submit(req.Job)
	return &pb.ScheduleResponse{}, nil
}

func (s *Server) GetJob(ctx context.Context, req *pb.GetJobRequest) (*pb.GetJobResponse, error) {
	return &pb.GetJobResponse{}, nil
}

// FeedbackService
func (s *Server) JobStarted(ctx context.Context, req *pb.JobStartedRequest) (*pb.JobStartedResponse, error) {
	return &pb.JobStartedResponse{}, nil
}

func (s *Server) JobProgress(stream pb.FeedbackService_JobProgressServer) error {
	var err error
	// var step *pb.JobProgressStep

	for {
		_, err = stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.JobProgressResponse{})
		}
		if err != nil {
			return err
		}
		// TODO: process output from job
	}
}

func (s *Server) JobStopped(ctx context.Context, req *pb.JobStoppedRequest) (*pb.JobStoppedResponse, error) {
	return &pb.JobStoppedResponse{}, nil
}

func (s *Server) loop() error {
	var id scheduler.JobId

	for {
		select {
		case id = <-s.Scheduler.Next:
			job := s.jobForId(id)
			s.Store.SetJob(job)
		case id = <-s.Scheduler.Expired:
			s.removeJob(id)
			s.Store.SetSchedule(s.jobs)
		}
	}
}

func (s *Server) jobForId(id scheduler.JobId) *pb.Job {
	for _, job := range s.jobs {
		if job.Id == uint64(id) {
			return job
		}
	}
	return nil
}

func (s *Server) removeJob(jobId scheduler.JobId) {
	var i int
	var job *pb.Job

	for i, job = range s.jobs {
		if job.Id == uint64(jobId) {
			break
		}
	}
	if i < len(s.jobs) {
		s.jobs = append(s.jobs[:i], s.jobs[i+1:]...)
	}
}
