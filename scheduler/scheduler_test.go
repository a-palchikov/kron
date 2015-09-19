package scheduler

import (
	"testing"
	"time"

	pb "github.com/a-palchikov/kron/proto/servicepb"
)

var (
	now  = parseTime("2000-Jan-01 00:00:00")
	fake = &fakeClock{now}
)

func TestSubmitsJob(t *testing.T) {
	scheduler := createScheduler(t)
	job := &pb.Job{Id: 1, When: "0 0 * * * * 2001"}

	scheduler.Submit(job)
	<-scheduler.Next
}

func TestRemovesExpiredJobs(t *testing.T) {
	scheduler := createScheduler(t)
	job := &pb.Job{Id: 1, When: "0 0 * * * * 1999"}
	scheduler.Submit(job)

	select {
	case expiredId := <-scheduler.Expired:
		if JobId(job.Id) != expiredId {
			t.Fatalf("expected job with id %d, but got %d", job.Id, expiredId)
		}
	}
}

func TestSchedulesJobsInProperOrder(t *testing.T) {
	scheduler := createScheduler(t)
	jobs := []*pb.Job{
		&pb.Job{Id: 1, When: "3 0 1 1 1 * 2000"},
		&pb.Job{Id: 2, When: "2 0 1 1 1 * 2000"},
		&pb.Job{Id: 3, When: "1 0 1 1 1 * 2000"},
	}
	when := parseTime("2000-Jan-01 01:00:00")
	expected := []JobId{3, 2, 1}
	var actual []JobId

	go func() {
		for {
			select {
			case <-scheduler.Expired:
			case id := <-scheduler.Next:
				actual = append(actual, id)
			}
		}
	}()

	fake.rewind(when)
	scheduler.SubmitMultiple(jobs)
	// FIXME: spend some time
	time.Sleep(100 * time.Millisecond)

	assertEqual(t, expected, actual)
}

func createScheduler(t *testing.T) *Scheduler {
	scheduler, err := NewWithClock(fake)
	if err != nil {
		t.Fatalf("cannot create scheduler")
	}
	go scheduler.Run()
	return scheduler
}

func parseTime(value string) time.Time {
	const shortForm = "2006-Jan-01 15:04:05"
	t, err := time.ParseInLocation(shortForm, value, time.Local)
	if err != nil {
		panic(err)
	}
	return t
}

func assertEqual(t *testing.T, a, b []JobId) {
	if len(a) != len(b) {
		t.Fatalf("a, b differ in len: %d != %d", len(a), len(b))
	}
	for i, id := range b {
		if a[i] != id {
			t.Fatalf("a,b differ at %d: %v != %v", i, a[i], id)
		}
	}
}
