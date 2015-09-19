package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"

	"golang.org/x/net/context"

	pb "github.com/a-palchikov/kron/proto/storepb"
	zmq "github.com/pebbe/zmq4"
)

var (
	apiPort  = flag.Int("bind-port", 5555, "API port")
	pushPort = flag.Int("push-port", 5556, "PUSH port")
)

type (
	server struct {
		socket  *zmq.Socket
		done    chan struct{}
		mu      *sync.Mutex // guards values
		values  map[string][]byte
		updates chan keyUpdate
	}
	keyUpdate struct {
		key     string
		payload []byte
	}
)

func main() {
	flag.Parse()

	var err error
	s := &server{
		done:    make(chan struct{}),
		values:  make(map[string][]byte),
		mu:      &sync.Mutex{},
		updates: make(chan keyUpdate),
	}

	if err = createPushServer(s); err != nil {
		log.Fatalln(err)
	}
	err = createAndRunApiServer(s)
	s.Close()
	log.Fatalln(err)
}

func createAndRunApiServer(impl *server) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *apiPort))
	if err != nil {
		return errors.New(fmt.Sprintf("cannot listen on %d: %v", *apiPort, err))
	}
	server := grpc.NewServer()

	log.Printf("API server on %d", *apiPort)
	pb.RegisterStoreServiceServer(server, impl)
	return server.Serve(listener)
}

func createPushServer(impl *server) error {
	ctx, err := zmq.NewContext()
	if err != nil {
		return err
	}
	socket, err := ctx.NewSocket(zmq.PUB)
	if err != nil {
		return err
	}
	if err = socket.Bind(fmt.Sprintf("tcp://*:%d", *pushPort)); err != nil {
		return err
	}
	impl.socket = socket
	log.Printf("PUSH service on %d", *pushPort)
	go impl.loop()
	return nil
}

func (s *server) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[req.Key] = req.Value
	s.updates <- keyUpdate{req.Key, req.Value}
	return &pb.SetResponse{}, nil
}

func (s *server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if value, ok := s.values[req.Key]; ok {
		return &pb.GetResponse{Value: value}, nil
	} else {
		return &pb.GetResponse{Error: pb.GetResponse_NotFound}, nil
	}
}

func (s *server) Close() {
	s.done <- struct{}{}
	<-s.done
	s.socket.Close()
}

func (s *server) loop() {
	var err error

	for {
		select {
		case update := <-s.updates:
			// FIXME: report error
			if err = send(s.socket, update); err != nil {
				log.Printf("cannot push update for key=%s", update.key)
			}
		case <-s.done:
			s.done <- struct{}{}
			return
		}
	}
}

func send(socket *zmq.Socket, update keyUpdate) error {
	var err error

	if _, err = socket.Send(update.key, zmq.SNDMORE); err != nil {
		return err
	}
	if _, err = socket.SendBytes(update.payload, 0); err != nil {
		return err
	}
	return nil
}
