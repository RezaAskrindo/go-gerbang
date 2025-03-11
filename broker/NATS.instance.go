package broker

import (
	"fmt"
	"log"
	"time"

	"go-gerbang/config"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

func StartingNatsServer() (*server.Server, error) {
	natsServer, err := server.NewServer(&server.Options{
		Host: "127.0.0.1",
		// Port: 4222, // Default NATS port
		Port: 9001,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create NATS server: %w", err)
	}

	go natsServer.Start()

	if !natsServer.ReadyForConnections(10 * time.Second) {
		return nil, fmt.Errorf("NATS server failed to start")
	}

	fmt.Printf("Embedded NATS server started on %s\n", natsServer.ClientURL())

	return natsServer, nil
}

var NatsClient *nats.Conn

// StartingNatsClient initializes a connection to the NATS server using the URL from the configuration.
func StartingNatsClient() {
	serverURL := config.Config("NATS_SERVER_URL")

	var err error
	NatsClient, err = nats.Connect(serverURL)
	if err != nil {
		log.Printf("Error connecting to NATS server: %v", err)
	}

	log.Printf("Connected to NATS server at:%s", serverURL)
}
