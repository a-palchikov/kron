package main

import (
	"flag"
	"fmt"
	"log"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	pb "github.com/a-palchikov/kron/proto/servicepb"
	storepb "github.com/a-palchikov/kron/proto/storepb"
	"github.com/a-palchikov/kron/service"
	zmq "github.com/pebbe/zmq4"
)

var (
	isMaster              = flag.Bool("master", true, "force master node")
	apiPort               = flag.Int("api", 5557, "api server")
	feedbackPort          = flag.Int("feedback", 5558, "feedback server")
	storeApiAddr          = flag.String("storeApi", ":5555", "store api server")
	storeSubscriptionAddr = flag.String("storeSubscription", ":5556", "store subscriptions service")
)

func main() {
	flag.Parse()

	config := service.Config{
		Master:       *isMaster,
		ApiPort:      *apiPort,
		FeedbackPort: *feedbackPort,
	}
	store, err := connectToStore(*storeApiAddr, *storeSubscriptionAddr)
	if err != nil {
		log.Fatalf("cannot connect to store: %v", err)
	}
	server, err := service.New(&config, store, store)
	if err != nil {
		log.Fatalf("cannot create service: %v", err)
	}
	log.Fatalln(server.Serve())
}

type store struct {
	socket *zmq.Socket
	conn   *grpc.ClientConn
	client storepb.StoreServiceClient
}

// connectToStore connects to the store service using the following two endpoints:
//  * API: push notifications from master to store
//  * subscription service: notifications from store to nodes
func connectToStore(apiAddr string, subscriptionAddr string) (*store, error) {
	context, err := zmq.NewContext()

	if err != nil {
		return nil, err
	}
	socket, err := context.NewSocket(zmq.SUB)
	if err != nil {
		return nil, err
	}
	err = socket.Connect(fmt.Sprintf("tcp://%s", subscriptionAddr))
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(apiAddr)
	if err != nil {
		return nil, err
	}
	client := storepb.NewStoreServiceClient(conn)

	result := &store{socket: socket, conn: conn, client: client}

	go result.loop()
	return result, nil
}

func (s *store) Close() error {
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

func (s *store) loop() {
	// TODO: handle subscription notifications
}

// TODO: push notifications from store to server
