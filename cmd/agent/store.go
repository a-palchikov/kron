package main

import (
	"fmt"
	"log"

	"github.com/a-palchikov/kron/cmd/agent/internal/router"
	pb "github.com/a-palchikov/kron/proto/servicepb"
	storepb "github.com/a-palchikov/kron/proto/storepb"
	"github.com/a-palchikov/kron/service"
	zmq "github.com/pebbe/zmq4"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const updateEndpoint = "ipc://updates"

// store implements service.Store/service.Cluster
type store struct {
	// socket   *zmq.Socket
	conn     *grpc.ClientConn
	client   storepb.StoreServiceClient
	observer service.Observer
	// channels to receive subscriptions from Store
	schedules chan []byte
	jobs      chan []byte
	mux       *router.Router
	// masterChanges chan []byte
	done chan struct{}
}

// connectToStore connects to the Store service using the following two endpoints:
//  * API: push notifications from master to store
//  * subscription service: notifications from store to nodes
func connectToStore(apiAddr string, subscriptionAddr string) (*store, error) {
	ctx, err := zmq.NewContext()
	if err != nil {
		return nil, err
	}

	// frontend subscribes to all topics on the Store subscriptions service
	frontend, err := ctx.NewSocket(zmq.SUB)
	if err != nil {
		return nil, err
	}
	err = frontend.Connect(fmt.Sprintf("tcp://%s", subscriptionAddr))
	if err != nil {
		return nil, err
	}

	// backend is a local publisher that will help distribute topics locally
	backend, err := ctx.NewSocket(zmq.PUB)
	if err != nil {
		return nil, err
	}
	err = backend.Bind(updateEndpoint)
	if err != nil {
		return nil, err
	}

	if err = zmq.Proxy(frontend, backend, nil); err != nil {
		return nil, err
	}

	mux, err := router.New(updateEndpoint)
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(apiAddr)
	if err != nil {
		return nil, err
	}
	client := storepb.NewStoreServiceClient(conn)

	s := &store{
		conn:      conn,
		client:    client,
		schedules: make(chan []byte),
		jobs:      make(chan []byte),
		mux:       mux,
		done:      make(chan struct{}),
	}

	if err = mux.Add(service.ScheduleUpdateNotification, s.schedules); err != nil {
		return nil, err
	}
	if err = mux.Add(service.JobUpdateNotification, s.jobs); err != nil {
		return nil, err
	}

	go s.loop()
	go mux.Run()
	return s, nil
}

func (s *store) Close() error {
	s.done <- struct{}{}
	<-s.done
	return s.conn.Close()
}

func (s *store) SetSchedule(jobs []*pb.Job) error {
	var err error

	_, err = s.client.SetSchedule(context.Background(), &storepb.SetScheduleRequest{
		Jobs: jobs,
	})
	return err
}

func (s *store) SetJob(job *pb.Job) error {
	var err error

	_, err = s.client.ScheduleJob(context.Background(), &storepb.ScheduleJobRequest{
		Job: job,
	})
	return err
}

func (s *store) AddNode(host string) (service.NodeId, error) {
	return 0, nil
}

func (s *store) Nodes() ([]service.Node, error) {
	return nil, nil
}

func (s *store) SetObserver(observer service.Observer) {
	s.observer = observer
}

func (s *store) loop() {
	var err error
	var payload []byte

	for {
		select {
		case payload = <-s.schedules:
			var schedule []*pb.Job

			if err = decode(payload, &schedule); err != nil {
				log.Printf("cannot decode schedule: %v", err)
			} else {
				s.observer.JobScheduleUpdate(schedule)
			}
		case payload = <-s.jobs:
			var req pb.ScheduleRequest

			if err = decode(payload, &req); err != nil {
				log.Printf("cannot decode job: %v", err)
			} else {
				s.observer.JobRequestUpdate(&req)
			}
			/*
				case payload = <-s.masterChanges:
					change, _ := decodeMasterChange(payload)
					s.notifier.ScheduleUpdate(change)
			*/
		case <-s.done:
			s.mux.Close()
			s.done <- struct{}{}
			return
		}
	}
}
