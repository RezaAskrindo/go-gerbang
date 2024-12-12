package broker

import (
	"fmt"
	"log"

	"go-gerbang/config"

	// "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

// func StartingNatsServer() (*server.Server, error) {
// 	natsServer, err := server.NewServer(&server.Options{
// 		Host: "127.0.0.1",
// 		Port: 4222, // Default NATS port
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create NATS server: %w", err)
// 	}

// 	go natsServer.Start()

// 	if !natsServer.ReadyForConnections(10 * time.Second) {
// 		return nil, fmt.Errorf("NATS server failed to start")
// 	}

// 	fmt.Printf("Embedded NATS server started on %s\n", natsServer.ClientURL())

// 	return natsServer, nil
// }

var NatsClient *nats.Conn

func StartingNatsClient() {
	serverURL := config.Config("NATS_SERVER_URL")

	var err error
	NatsClient, err = nats.Connect(serverURL)
	if err != nil {
		log.Printf("Error connecting to NATS server: %v", err)
	}

	fmt.Println("Connected to NATS server at:", serverURL)

	// defer NatsClient.Close()
}