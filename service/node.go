package service

import (
	"os"

	"github.com/a-palchikov/kron/proc"
	pb "github.com/a-palchikov/kron/proto/servicepb"
)

// Worker node specific functionality
type Worker struct {
	Cluster
	Store
	id uint64
	// redundant schedule for this node, if it gets elected to be a new master
	jobs           []*pb.Job
	queue          chan *pb.ScheduleRequest // worker node execution queue
	tags           map[string]struct{}
	feedbackClient *pb.FeedbackServiceClient
}

func NewWorker(store Store, cluster Cluster, tags []string) (*Worker, error) {
	w := &Worker{
		Cluster: cluster,
		Store:   store,
		queue:   make(chan *pb.ScheduleRequest),
		tags:    make(map[string]struct{}),
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

	go w.loop()
	return w, nil
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
		}
	}
}
