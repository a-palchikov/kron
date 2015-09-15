package main

import (
	"flag"

	"github.com/a-palchikov/kron/pubsub"
)

var receiver = flag.Int("receiver", ":5056", "message receiver port")
var sender = flag.Int("sender", ":5057", "message sender port")

func main() {
	flag.Parse()

	pubsub.Serve(receiver, sender)
}
