package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"golang.org/x/net/context"

	api "github.com/a-palchikov/kron/proto/servicepb"
	pb "github.com/a-palchikov/kron/proto/storepb"
	zmq "github.com/pebbe/zmq4"
)

var (
	serverPort = flag.Int("server", 5555, "gRPC server endpoint")
	storePort  = flag.Int("store", 5556, "store service")
)

type server struct {
	schedules chan []*api.Job
	jobs      chan *api.Job
	quit      chan struct{}
	done      chan struct{}
}

func main() {
	flag.Parse()

	srv := &server{
		schedules: make(chan []*api.Job),
		jobs:      make(chan *api.Job),
		done:      make(chan struct{}),
		quit:      make(chan struct{}),
	}

	var err error
	if err = createStoreServer(srv); err != nil {
		log.Fatalln(err)
	}
	err = createAndRunApiServer(srv)
	srv.quit <- struct{}{}
	<-srv.done
	log.Fatalln(err)
}

func createAndRunApiServer(s *server) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", serverPort))
	if err != nil {
		return errors.New(fmt.Sprintf("cannot listen on %d: %v", serverPort, err))
	}
	apiServer := grpc.NewServer()
	pb.RegisterStoreServiceServer(apiServer, s)
	return apiServer.Serve(listener)
}

func createStoreServer(s *server) error {
	ctx, err := zmq.NewContext()
	if err != nil {
		return err
	}
	socket, err := ctx.NewSocket(zmq.PUB)
	socket.Bind(fmt.Sprintf("tcp://*:%d", storePort))
	if err != nil {
		return err
	}
	go s.loop(socket)
	return nil
}

func (s *server) SetSchedule(ctx context.Context, req *pb.SetScheduleRequest) (*pb.SetScheduleResponse, error) {
	s.schedules <- req.Jobs
	return &pb.SetScheduleResponse{}, nil
}

func (s *server) ScheduleJob(ctx context.Context, req *pb.ScheduleJobRequest) (*pb.ScheduleJobResponse, error) {
	s.jobs <- req.Job
	return &pb.ScheduleJobResponse{}, nil
}

func (s *server) loop(socket *zmq.Socket) {
	defer socket.Close()

	var message []byte

	for {
		select {
		case schedule := <-s.schedules:
			// FIXME: report error
			message, _ = encode(schedule)
			socket.SendBytes(message, 0)
		case job := <-s.jobs:
			// FIXME: report error
			message, _ = encode(job)
			socket.SendBytes(message, 0)
		case <-s.quit:
			s.done <- struct{}{}
			return
		}
	}
}

func encode(value interface{}) ([]byte, error) {
	if b, err := json.Marshal(value); err != nil {
		return nil, err
	} else {
		return b, nil
	}
}
