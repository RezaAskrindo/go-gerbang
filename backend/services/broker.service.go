package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go-gerbang/broker"
	"go-gerbang/config"
	"go-gerbang/handlers"
	"go-gerbang/models"

	"github.com/gofiber/fiber/v3"
	"github.com/nats-io/nats.go"
)

func PublishService(c fiber.Ctx) error {
	publish := c.Query("p")

	if publish == "" {
		return handlers.UnprocessableEntityErrorResponse(c, fmt.Errorf("need base p query"))
	}

	subject := ""
	payload := map[string]any{}

	switch publish {
	case "wa":
		accID := config.Config("WA_ACCOUNT_ID")
		subject = "whatsapp.send"
		message := fmt.Sprintf("Testing With Event %s", time.Now().Format(time.RFC3339))
		payload = map[string]interface{}{
			"account_id": accID,
			"to":         "6285724416179",
			"message":    message,
		}
	case "mail":
		accID := config.Config("EMAIL_ACCOUNT_ID")
		subject = "email.send"
		message := fmt.Sprintf("<p>Testing With Event %s</p>", time.Now().Format(time.RFC3339))
		payload = map[string]interface{}{
			"account_id": accID,
			"recipients": []string{"rezaoda@gmail.com"},
			"subject":    "Testing Email",
			"body_html":  message,
		}
	}

	PublishEvent(subject, payload)

	return c.Status(fiber.StatusCreated).SendString("")

	// messageBytes, err := json.Marshal(payload)
	// if err != nil {
	// 	return handlers.InternalServerErrorResponse(c, err)
	// }
	// if err := broker.NatsClient.Publish(subject, []byte(messageBytes)); err != nil {
	// 	return handlers.InternalServerErrorResponse(c, err)
	// }
	// return handlers.SuccessResponse(c, true, "on publish", nil, nil)
}

// Not use anymore?
// func SubscribeService(c fiber.Ctx) error {
// 	subject := "send_mail"
// 	// broker.NatsClient.Subscribe(subject, func(msg *nats.Msg) {
// 	// 	fmt.Printf("Received message on %s: %s\n", subject, string(msg.Data))
// 	// })

// 	return handlers.SuccessResponse(c, true, "on subscribe", subject, nil)
// }

func PublishEvent(subject string, rawData interface{}) {
	data, err := json.Marshal(rawData)
	if err != nil {
		log.Printf("Error marshaling data: %v", err)
	}
	if err := broker.NatsClient.Publish(subject, []byte(data)); err != nil {
		log.Printf("Error publishing to subject %s: %v", subject, err)
	}
	// PUT ON DATABASE
	// log.Printf("Published to %s: %+v\n", subject, rawData)
}

func SubscribeEvent() {
	handlers := map[string]nats.MsgHandler{
		"logger":           handlers.HandleLogger,
		"users.find_by_id": handleFindById,
	}

	for subject, handler := range handlers {
		if _, err := broker.NatsClient.Subscribe(subject, handler); err != nil {
			log.Println("Error on subcribe NATS:", err)
		}
	}

	// log.Println("Listening for subcribe events...")
}

type natsResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    any    `json:"data,omitempty"`
}

func respond(msg *nats.Msg, success bool, code int, message string, data any) {
	raw, err := json.Marshal(natsResponse{success, message, code, data})
	if err != nil {
		return
	}
	_ = msg.Respond(raw)
}

func handleFindById(msg *nats.Msg) {
	var req struct {
		UserId string `json:"user_id"`
	}

	if err := json.Unmarshal(msg.Data, &req); err != nil || req.UserId == "" {
		respond(msg, false, 422, "need userId params", nil)
		return
	}

	user := new(models.User)
	if err := models.FindUserById(user, req.UserId); err != nil {
		respond(msg, false, 404, err.Error(), nil)
		return
	}

	respond(msg, true, 200, "success to get users", user)
}
