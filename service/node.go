package service

import (
	"os"

	"google.golang.org/grpc"

	"github.com/a-palchikov/kron/proc"
	pb "github.com/a-palchikov/kron/proto/servicepb"
)

// Worker node specific functionality
type Worker struct {
	config *Config
	id     uint64
	jobs   []*pb.Job                // redundant job schedule
	queue  chan *pb.ScheduleRequest // job execution queue
	tags   map[string]struct{}
	// Feedback client
	conn   *grpc.ClientConn
	client pb.FeedbackServiceClient
}

func NewWorker(config *Config, store Store, cluster Cluster, tags []string) (*Worker, error) {
	w := &Worker{
		queue: make(chan *pb.ScheduleRequest),
		tags:  make(map[string]struct{}),
	}

	host, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	var id NodeId

	if id, err = cluster.AddNode(host); err != nil {
		return nil, err
	}
	w.id = uint64(id)

	for _, tag := range tags {
		w.tags[tag] = struct{}{}
	}

	store.SetObserver(w)

	// go w.loop()
	return w, nil
}

func (w *Worker) initFeedbackClient() error {
	conn, err := grpc.Dial(w.config.FeedbackAddr)
	if err != nil {
		return err
	}
	w.client = pb.NewFeedbackServiceClient(conn)
	return nil
}

func (w *Worker) closeFeedbackClient() error {
	err := w.conn.Close()
	w.client = nil
	return err
}

// Observer
func (w *Worker) JobScheduleUpdate(jobs []*pb.Job) {
	w.jobs = jobs
}

func (w *Worker) JobRequestUpdate(req *pb.ScheduleRequest) {
	if w.canExecute(req) {
		// schedule the job for execution
		w.queue <- req
	}
}

// canExecute decides if the worker can execute the specified job.
// The decision is based on matching against the set of tags the worker has assigned.
// Further refinements can include checking against machine configuration, etc.
func (w *Worker) canExecute(req *pb.ScheduleRequest) bool {
	var ok bool

	for _, tag := range req.Tags {
		if _, ok = w.tags[tag]; ok {
			return true
		}
	}
	return false
}

func (w *Worker) loop() {
	var req *pb.ScheduleRequest
	var cmd *proc.Cmd

	for {
		select {
		case req = <-w.queue:
			// FIXME: no blocking here
			cmd = proc.Command(req.Job.CommandLine, req.Quota)
			cmd.Run()
			// TODO: communicate job progress with the master node using feedback client
		}
	}
}
