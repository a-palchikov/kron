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
	apiPort, subscriptionPort int
	verbose                   bool
)

func init() {
	flag.IntVar(&apiPort, "api", 5555, "api server")
	flag.IntVar(&subscriptionPort, "store", 5556, "subscription service")
	flag.BoolVar(&verbose, "verbose", true, "vebosity")
}

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
	if err = createSubscriptionServer(srv); err != nil {
		log.Fatalln(err)
	}
	err = createAndRunApiServer(srv)
	srv.quit <- struct{}{}
	<-srv.done
	log.Fatalln(err)
}

func createAndRunApiServer(s *server) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", apiPort))
	if err != nil {
		return errors.New(fmt.Sprintf("cannot listen on %d: %v", apiPort, err))
	}
	if verbose {
		log.Printf("api server on %d", apiPort)
	}
	apiServer := grpc.NewServer()
	pb.RegisterStoreServiceServer(apiServer, s)
	return apiServer.Serve(listener)
}

func createSubscriptionServer(s *server) error {
	ctx, err := zmq.NewContext()
	if err != nil {
		return err
	}
	socket, err := ctx.NewSocket(zmq.PUB)
	socket.Bind(fmt.Sprintf("tcp://*:%d", subscriptionPort))
	if err != nil {
		return err
	}
	if verbose {
		log.Printf("subscription service on %d", subscriptionPort)
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

	var err error

	for {
		select {
		case schedule := <-s.schedules:
			// FIXME: report error
			if err = send(socket, schedule, service.ScheduleUpdateNotification); err != nil {
				log.Printf("cannot publish schedule=%v: %v", schedule, err)
			}
		case job := <-s.jobs:
			// FIXME: report error
			if err = send(socket, job, service.JobUpdateNotification); err != nil {
				log.Printf("cannot publish job=%v: %v", job, err)
			}
		case <-s.quit:
			s.done <- struct{}{}
			return
		}
	}
}

func send(socket *zmq.Socket, value interface{}, topic string) error {
	var b []byte
	var err error

	if b, err = json.Marshal(value); err != nil {
		return err
	}
	socket.Send(topic, zmq.SNDMORE)
	socket.SendBytes(b, 0)
	return nil
}
