package scheduler

import (
	"testing"

	pb "github.com/a-palchikov/kron/proto/servicepb"
)

// FIXME: use fake time for these tests

func TestGeneratesValidJobId(t *testing.T) {
	scheduler, err := New()
	if err != nil {
		t.Fatalf("cannot create scheduler")
	}
	go scheduler.Run()

	job := &pb.Job{When: "0 0 * * * * 2014"}

	id := scheduler.Submit(job)
	if id == 0 {
		t.Fatalf("expected non-zero id, but got: %d", id)
	}
}

func TestSubmitsJob(t *testing.T) {
	scheduler, err := New()
	if err != nil {
		t.Fatalf("cannot create scheduler")
	}
	go scheduler.Run()

	// time in future
	job := &pb.Job{When: "0 0 * * * * 2016"}

	scheduler.Submit(job)
	if len(scheduler.jobs) == 0 {
		t.Fatalf("expected a job in queue")
	}
}

func TestRemovesExpiredJobs(t *testing.T) {
	scheduler, err := New()
	if err != nil {
		t.Fatalf("cannot create scheduler")
	}
	go scheduler.Run()

	job := &pb.Job{When: "0 0 * * * * 2014"}

	id := scheduler.Submit(job)

	select {
	case jobId := <-scheduler.Expired:
		if id != jobId {
			t.Fatalf("expected job with id %d, but got %d", id, jobId)
		}
	}
}
