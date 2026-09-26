package broker

import (
	"fmt"
	"log"
	"time"

	"go-gerbang/config"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

const (
	UsernameNats = "go-gerbang"
	PasswordNats = "G0-gerb@ng-2026"
)

func StartingNatsServer() (*server.Server, error) {
	natsServer, err := server.NewServer(&server.Options{
		ServerName: "Go-Gerbang-Broker",
		Host:       "0.0.0.0",
		Port:       9001,
		// For Low VPS
		MaxConn:       50,
		MaxPayload:    128 * 1024,      // 128 KB
		MaxPending:    2 * 1024 * 1024, // 2 MB
		WriteDeadline: 5 * time.Second,
		// JetStream
		JetStream:          true,
		JetStreamMaxMemory: 32 * 1024 * 1024, // 32 MB
		StoreDir:           "./data",
		// Auth
		Username: UsernameNats,
		Password: PasswordNats,
		// MQTT
		// MQTT: server.MQTTOpts{
		// 	Host: "0.0.0.0",
		// 	Port: 9002,
		// 	// Username: UsernameNats,
		// 	// Password: PasswordNats,
		// },
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create NATS server: %w", err)
	}

	go natsServer.Start()

	if !natsServer.ReadyForConnections(10 * time.Second) {
		natsServer.Shutdown()
		return nil, fmt.Errorf("NATS server failed to start")
	}

	fmt.Printf("[INFO] NATS server running :9001\n")

	return natsServer, nil
}

var NatsClient *nats.Conn

func StartingNatsClient() {
	serverURL := config.Config("NATS_SERVER_URL")

	var err error
	NatsClient, err = nats.Connect(serverURL, nats.UserInfo(UsernameNats, PasswordNats))
	if err != nil {
		log.Printf("Error connecting to NATS server: %v", err)
	}

	NatsClient.Opts.MaxReconnect = 3
	NatsClient.Opts.ReconnectWait = 1 * time.Second
}
