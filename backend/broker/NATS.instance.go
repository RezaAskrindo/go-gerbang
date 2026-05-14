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
		Host: "0.0.0.0",
		Port: 9001,
	})
	if err != nil {
		log.Printf("failed to create NATS server: %v", err)
		return nil, fmt.Errorf("failed to create NATS server: %w", err)
	}

	go natsServer.Start()

	if !natsServer.ReadyForConnections(10 * time.Second) {
		log.Printf("NATS server failed to start")
		return nil, fmt.Errorf("NATS server failed to start")
	}

	fmt.Printf("✅ NATS server running :9001\n")

	return natsServer, nil
}

var NatsClient *nats.Conn

func StartingNatsClient() {
	serverURL := config.Config("NATS_SERVER_URL")

	var err error
	NatsClient, err = nats.Connect(serverURL)
	if err != nil {
		log.Printf("Error connecting to NATS server: %v", err)
	}

	fmt.Printf("✅ NATS client connected at:%s\n", serverURL)
}
