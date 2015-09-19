// Package router implements routing of ZMQ subscriptions to Go channels
package router

import (
	"time"

	zmq "github.com/pebbe/zmq4"
)

type (
	Router struct {
		ctx       *zmq.Context
		receivers map[*zmq.Socket]*receiver
		poller    *zmq.Poller
		addr      string
		done      bool
	}

	receiver struct {
		sink chan<- []byte
	}
)

func New(addr string) (*Router, error) {
	ctx, err := zmq.NewContext()
	if err != nil {
		return nil, err
	}
	router := &Router{
		ctx:       ctx,
		receivers: make(map[*zmq.Socket]*receiver),
		poller:    zmq.NewPoller(),
		addr:      addr,
	}

	return router, nil
}

// Add adds a new sink channel for messages of topic `topic`
func (r *Router) Add(topic string, sink chan<- []byte) error {
	var err error
	var socket *zmq.Socket

	socket, err = r.ctx.NewSocket(zmq.SUB)
	if err != nil {
		return err
	}
	if err = socket.Connect(r.addr); err != nil {
		return err
	}
	if err = socket.SetSubscribe(topic); err != nil {
		return err
	}
	r.poller.Add(socket, zmq.POLLIN)
	r.receivers[socket] = &receiver{sink: sink}
	return nil
}

// Close terminates the router loop
func (s *Router) Close() {
	s.done = true
}

func (s *Router) Run() error {
	var err error
	var items []zmq.Polled
	var payload []byte
	const interval = 1 * time.Second

	for !s.done {
		items, err = s.poller.Poll(interval)
		if err != nil {
			return err
		}
		for _, item := range items {
			if receiver, ok := s.receivers[item.Socket]; ok {
				// receive the name of the topic
				if _, err = item.Socket.RecvBytes(0); err != nil {
					return err
				}
				if payload, err = item.Socket.RecvBytes(0); err != nil {
					return err
				}
				receiver.sink <- payload
			}
		}
	}
	return nil
}
