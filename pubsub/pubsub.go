// package pubsub implements a basic publish/subscribe over zeromq.
package pubsub

import (
	"github.com/pebbe/zmq4"
)

func Serve(receiverPort, senderPort int) error {
	context, err := zmq.NewContext()
	if err != nil {
		return err
	}

	receiver, err := context.NewSocket(zmq.PULL)
	if err != nil {
		return err
	}
	defer receiver.Close()
	if err = receiver.Bind(fmt.Sprintf("tcp://*:%d", receiverPort)); err != nil {
		return err
	}

	sender, err := context.NewSocket(zmq.PUB)
	if err != nil {
		return err
	}
	defer sender.Close()
	if err = sender.Bind(fmt.Sprintf("tcp://*:%d", senderPort)); err != nil {
		return err
	}

	var message string
	for {
		if message, err = receiver.Recv(); err != nil {
			break
		}
		if err = sender.SendMessage(message); err != nil {
			break
		}
	}

	return err
}
