package scheduler

import (
	// "log"
	"sort"
	"time"

	pb "github.com/a-palchikov/kron/proto/servicepb"
	"github.com/gorhill/cronexpr"
)

type (
	JobId uint64

	job struct {
		id   JobId
		expr *cronexpr.Expression
		next time.Time
	}

	Scheduler struct {
		// Jobs are scheduled on this channel
		Next chan JobId
		// Expired jobs are reported here
		Expired chan JobId
		queue   chan *job
		jobs    []*job
		done    chan struct{}
	}
)

func New() (*Scheduler, error) {
	s := &Scheduler{
		Next:    make(chan JobId),
		Expired: make(chan JobId),
		queue:   make(chan *job),
		done:    make(chan struct{}),
	}

	return s, nil
}

func (s *Scheduler) Submit(job *pb.Job) JobId {
	id := nextJobId()
	s.queue <- newJob(job, id)
	return id
}

func newJob(j *pb.Job, id JobId) *job {
	job := &job{
		id:   id,
		expr: cronexpr.MustParse(j.When),
	}
	return job
}

func (s *Scheduler) Close() {
	s.done <- struct{}{}
	<-s.done
}

func (s *Scheduler) Run() {
	nextTime := time.Now()

	for {
		var tick <-chan time.Time
		var newJob *job
		var scheduled []*job

		sort.Sort(byTime{s.jobs})

		s.removeExpired(nextTime)

		if len(s.jobs) > 0 {
			tickDuration := s.jobs[0].next.Sub(nextTime)
			tick = time.After(tickDuration)
			scheduled = append(scheduled, s.jobs[0])
			// aggregate tasks with the same time in future into scheduled
			for _, job := range s.jobs[1:] {
				if job.next.Sub(nextTime) == tickDuration {
					scheduled = append(scheduled, job)
				}
			}
		}

		select {
		case newJob = <-s.queue:
			nextTime = time.Now()
			s.jobs = append(s.jobs, newJob)
			newJob.next = newJob.expr.Next(nextTime)
			// log.Printf("new job %v for %v", newJob, newJob.next)
		case <-tick:
			nextTime = time.Now()
			for _, job := range scheduled {
				s.Next <- job.id
				job.next = job.expr.Next(nextTime)
			}
		case <-s.done:
			s.done <- struct{}{}
			return
		}
	}
}

func (s *Scheduler) removeExpired(now time.Time) {
	numExpired := s.numExpired(now)
	expired := s.jobs[:numExpired]
	s.jobs = s.jobs[numExpired:]
	for _, job := range expired {
		s.Expired <- job.id
	}
}

// numExpired returns the number of jobs to remove from the head
func (s *Scheduler) numExpired(now time.Time) int {
	var pos int
	var job *job

	for _, job = range s.jobs {
		if !job.expr.Next(now).IsZero() {
			break
		}
		pos++
	}
	return pos
}

func nextJobId() JobId {
	return JobId(time.Now().UnixNano() / 100)
}

type byTime struct {
	jobs []*job
}

func (c byTime) Len() int           { return len(c.jobs) }
func (c byTime) Less(i, j int) bool { return c.jobs[i].next.Before(c.jobs[j].next) }
func (c byTime) Swap(i, j int)      { c.jobs[i], c.jobs[j] = c.jobs[j], c.jobs[i] }
