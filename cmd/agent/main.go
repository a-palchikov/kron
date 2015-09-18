package main

import (
	"flag"
	"log"

	"github.com/a-palchikov/kron/service"
)

var (
	master                = flag.Bool("master", true, "master mode")
	apiPort               = flag.Int("api", 5557, "api server")
	feedbackPort          = flag.Int("feedback", 5558, "feedback server")
	storeApiAddr          = flag.String("storeApi", ":5555", "store api server")
	storeSubscriptionAddr = flag.String("storeSubscription", ":5556", "store subscriptions service")
)

func main() {
	flag.Parse()

	config := service.Config{
		Master:       *master,
		ApiPort:      *apiPort,
		FeedbackPort: *feedbackPort,
	}
	store, err := connectToStore(*storeApiAddr, *storeSubscriptionAddr)
	if err != nil {
		log.Fatalf("cannot connect to store: %v", err)
	}
	server, err := service.New(&config, store, store, nil)
	if err != nil {
		log.Fatalf("cannot create service: %v", err)
	}
	log.Fatalln(server.Serve())
}
