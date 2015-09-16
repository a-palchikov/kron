package main

import (
	"flag"
	"fmt"

	"github.com/a-palchikov/kron/service"
	pb "github.com/a-palchikov/kron/service/proto"
	"github.com/pebbe/zmq4"
)

var (
	isMaster           = flag.Bool("master", true, "force master node")
	infrastructureAddr = flag.String("infrastructure", ":5555", "address of persistency/membership layer")
	a                  = flag.String("", ":5555", "")
)

func main() {
	flag.Parse()

	config := Config{
		master: isMaster,
	}
	infrastructure := createInfrastructure()
	// TODO:
	//  * init pubsub endpoint (server: publisher, nodes: subscribers)
	//  *
	server, err := service.New(&config, &infrastructure)
	// server.Serve()
}

type infrastructure struct {
	socket *zmq.Socket
}

func createInfrastructure(addr string) (service.Infrastructure, error) {
	// TODO
	var err error
	context, err := zmq.Context()
	if err != nil {
		return nil, err
	}
	socket, err := context.NewSocket(zmq.PULL)
	if err != nil {
		return nil, err
	}
	err = socket.Bind(fmt.Sprintf("tcp://%s", addr))
	if err != nil {
		return nil, err
	}
	result := &infrastructure{socket: socket}

	go result.loop()
	return result, nil
}

func (i *infrastructure) SetSchedule([]*pb.Job) error {
	return nil
}

func (i *infrastructure) SetJob(*pb.Job) error {
	return nil
}

func (i *infrastructure) AddNode(host string) (service.NodeId, error) {
	return 0, nil
}

func (i *infrastructure) loop() {
	for {
		message, _ := i.socket.Recv(0)

		// TODO: decode message
	}
}
