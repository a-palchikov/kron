package service

import pb "github.com/a-palchikov/kron/proto/servicepb"

const (
	ScheduleUpdateNotification = "SetSchedule"
	JobUpdateNotification      = "SetJob"
)

type Observer interface {
	// ScheduleUpdate is invoked by Store when the master node pushes schedule update
	JobScheduleUpdate([]*pb.Job)
	// JobUpdate is invoked by Store when a new job is scheduled
	// The worker nodes will use it to decide if they can execute the job
	JobRequestUpdate(*pb.ScheduleRequest)
}

type Store interface {
	// SetSchedule is used by the master node to publish schedules so that all nodes are in sync
	// with regard to the current schedule
	SetSchedule([]*pb.Job) error
	// SetJob is used by the master node to publishe work to worker nodes
	SetJob(*pb.Job) error
	// SetObserver adds a new observer for updates
	SetObserver(Observer)
	Close() error
	// TODO: expose feedback service address?
	// If the master is re-elected (after the failure of the old master), feedback service
	// needs to be available on the new master node.
}
