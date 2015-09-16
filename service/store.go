package service

import pb "github.com/a-palchikov/kron/proto/servicepb"

type Store interface {
	// SetSchedule publishes given schedule so that nodes stay in sync
	SetSchedule([]*pb.Job) error
	// SetJob publishes scheduled job on the store topic so that matching nodes
	// can pick it up
	SetJob(*pb.Job) error
	// TODO: expose feedback service address?
	// If the master is re-elected (after the failure of the old master), feedback service
	// needs to be available on the new master node.
}
