package router

import (
	"bytes"
	"log"
	"testing"

	zmq "github.com/pebbe/zmq4"
)

var _ = log.Printf

func TestSingleSubscription(t *testing.T) {
	testMessage := []byte("whoami")
	mux, publisher := createRouterAndPublisher(t)

	go func() {
		publisher.Send("A", zmq.SNDMORE)
		publisher.SendBytes(testMessage, 0)
	}()

	ch := make(chan []byte)
	if err := mux.Add("A", ch); err != nil {
		t.Fatalf("cannot add subscription")
	}
	go mux.Run()

	msg := <-ch
	if !bytes.Equal(msg, testMessage) {
		t.Fatalf("expected `%s` but got `%s`", testMessage, msg)
	}
}

func TestDoesnotReceiveInvalidSubscription(t *testing.T) {
	testMessage := []byte("whoami")
	mux, publisher := createRouterAndPublisher(t)

	go func() {
		publisher.Send("A", zmq.SNDMORE)
		publisher.SendBytes(testMessage, 0)
	}()

	ch := make(chan []byte)
	if err := mux.Add("B", ch); err != nil {
		t.Fatalf("cannot add subscription")
	}
	go mux.Run()

	// FIXME: needs a better way to ensure that no message is received
	select {
	case msg := <-ch:
		t.Fatalf("expected no subscription but got `%s`", msg)
	default:
	}
}

func TestMultipleSubscriptions(t *testing.T) {
	messages := []struct {
		topic   string
		message []byte
		sink    chan []byte
	}{
		{"A", []byte("whoami"), make(chan []byte)},
		{"B", []byte("helloworld"), make(chan []byte)},
	}
	mux, publisher := createRouterAndPublisher(t)

	go func() {
		for _, message := range messages {
			publisher.Send(message.topic, zmq.SNDMORE)
			publisher.SendBytes(message.message, 0)
		}
	}()

	for _, message := range messages {
		if err := mux.Add(message.topic, message.sink); err != nil {
			t.Fatalf("cannot add subscription for %s", message.topic)
		}
	}
	go mux.Run()

	var numReceived int

	for numReceived < 2 {
		select {
		case msg := <-messages[0].sink:
			expected := messages[0].message
			if !bytes.Equal(msg, expected) {
				t.Fatalf("expected `%s` but got `%s`", expected, msg)
			}
		case msg := <-messages[1].sink:
			expected := messages[1].message
			if !bytes.Equal(msg, expected) {
				t.Fatalf("expected `%s` but got `%s`", expected, msg)
			}
		}
		numReceived++
	}
}

func createRouterAndPublisher(t *testing.T) (*Router, *zmq.Socket) {
	const addr = "ipc:///tmp/test.ipc"
	publisher, err := zmq.NewSocket(zmq.PUB)
	if err != nil {
		t.Fatalf("cannot create zmq socket")
	}
	if err = publisher.Bind(addr); err != nil {
		t.Fatalf("cannot publish to %s: %v", addr, err)
	}

	mux, err := New(addr)
	if err != nil {
		t.Fatalf("cannot create router")
	}

	return mux, publisher
}
