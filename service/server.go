package service

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"

	pb "github.com/a-palchikov/kron/proto/servicepb"
	"google.golang.org/grpc"

	"golang.org/x/net/context"
)

type (
	JobId uint64

	Server struct {
		Store
		Cluster
		config  *Config
		id      uint64
		mu      *sync.Mutex
		pending map[uint64][]uint64 // jobId->[]nodeId
		// scheduler.Scheduler	// cron management
	}
)

func New(config *Config, store Store, cluster Cluster) (*Server, error) {
	server := &Server{
		config:  config,
		Store:   store,
		Cluster: cluster,
	}

	host, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	var id NodeId

	if id, err = cluster.AddNode(host); err != nil {
		return nil, err
	}
	server.id = uint64(id)

	// go server.Scheduler.Loop(server.Store)
	// go s.loop()

	return server, nil
}

// Serve spawns serve loops for both API/feedback services.
func (s *Server) Serve() error {
	apiServer, apiListener, err := bindApi(s)
	if err != nil {
		return err
	}

	feedbackServer, feedbackListener, err := bindFeedback(s)
	if err != nil {
		return err
	}

	// FIXME: report serve errors from feedback service
	go feedbackServer.Serve(feedbackListener)

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
	s.mu.Lock()
	defer s.mu.Unlock()
	s.addPendingJob(req.Job)
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
	s.mu.Lock()
	defer s.mu.Unlock()
	s.removePendingJob(req.Job)
	return &pb.JobStoppedResponse{}, nil
}

// TODO: server loop
func (s *Server) loop() error {
	return nil
}

// Callers must have s.mu locked.
func (s *Server) addPendingJob(job *pb.JobRef) {
	var nodes []uint64
	var ok bool

	if nodes, ok = s.pending[job.JobId]; !ok {
		s.pending[job.JobId] = make([]uint64, 1)
	}
	nodes = append(nodes, job.NodeId)
}

// Callers must have s.mu locked.
func (s *Server) removePendingJob(job *pb.JobRef) {
	var ok bool

	if _, ok = s.pending[job.JobId]; !ok {
		return
	}
	delete(s.pending, job.JobId)
	// TODO: upon job completion, push new schedule to store to sync with all nodes
}
