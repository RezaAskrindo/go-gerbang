package services

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"go-gerbang/broker"
	"go-gerbang/config"
	"go-gerbang/handlers"
	"go-gerbang/models"

	"github.com/gofiber/fiber/v3"
	"github.com/nats-io/nats.go"
)

const (
	WaPublishEventName    = "whatsapp.send"
	EmailPublishEventName = "email.send"
)

func QueueUserInformation(accountId *string, user *models.User, sendPass bool) bool {
	accID := config.Config("EMAIL_ACCOUNT_ID")

	if accountId != nil {
		accID = *accountId
	}

	subject := EmailPublishEventName

	if accID != "" {
		pass := ""
		if sendPass {
			pass = user.Password
		}

		payload := map[string]interface{}{
			"account_id":    accID,
			"recipients":    []string{user.Email},
			"subject":       "Create Account Success",
			"template_name": "user_created",
			"template_data": map[string]any{
				"Name":     user.FullName,
				"Username": user.Username,
				"Email":    user.Email,
				"Password": pass,
			},
		}

		PublishEvent(subject, payload)
	}

	return true
}

func PublishService(c fiber.Ctx) error {
	publish := c.Query("p")

	if publish == "" {
		return handlers.UnprocessableEntityErrorResponse(c, fmt.Errorf("need base p query, wa or mail"))
	}

	accountID := c.Query("id")
	destination := c.Query("to")

	subject := ""
	payload := map[string]any{}

	switch publish {
	case "wa":
		accID := config.Config("WA_ACCOUNT_ID")
		if accountID != "" {
			accID = accountID
		}
		target := "6285724416179"
		if destination != "" {
			target = destination
		}
		subject = WaPublishEventName
		message := fmt.Sprintf("Testing With Event %s", time.Now().Format(time.RFC3339))
		payload = map[string]interface{}{
			"account_id": accID,
			"to":         target,
			"message":    message,
		}
	case "mail":
		accID := config.Config("EMAIL_ACCOUNT_ID")
		if accountID != "" {
			accID = accountID
		}
		target := []string{"rezaoda@gmail.com"}
		if destination != "" {
			target = []string{destination}
		}
		subject = EmailPublishEventName
		message := fmt.Sprintf("<p>Testing With Event %s</p>", time.Now().Format(time.RFC3339))
		payload = map[string]interface{}{
			"account_id": accID,
			"recipients": target,
			"subject":    "Testing Email",
			"body_html":  message,
		}
	}

	PublishEvent(subject, payload)

	return c.Status(fiber.StatusCreated).SendString("")
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
		"logger":                 handlers.HandleLogger,
		"users.find_by_id":       handleFindById,
		"users.find_by_identity": handleFindByIdentity,
		"users.create":           handleCreateUser,
		"users.update":           handleUpdateUser,
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

func handleFindByIdentity(msg *nats.Msg) {
	var req struct {
		Email          string `json:"email"`
		Username       string `json:"username"`
		PhoneNumber    string `json:"phone_number"`
		IdentityNumber string `json:"identity_number"`
	}

	if err := json.Unmarshal(msg.Data, &req); err != nil ||
		req.Email == "" ||
		req.Username == "" ||
		req.PhoneNumber == "" ||
		req.IdentityNumber == "" {
		respond(msg, false, 422, "need email or username or phone_number or identity number params", nil)
		return
	}

	user := new(models.User)
	if err := models.FindUserByIdentity(user, req.Username, req.Email, req.PhoneNumber, req.IdentityNumber); err != nil {
		respond(msg, false, 404, err.Error(), nil)
		return
	}

	respond(msg, true, 200, "success to get users", user)
}

func handleCreateUser(msg *nats.Msg) {
	var req struct {
		Username       string `json:"username"`
		Email          string `json:"email"`
		Password       string `json:"password"`
		PhoneNumber    string `json:"phone_number"`
		IdentityNumber string `json:"identity_number"`
		Active         bool   `json:"active"`
		SendNotif      bool   `json:"send_notif"`
		SendPass       bool   `json:"send_pass"`
		AccountId      string `json:"account_id"`
	}

	if err := json.Unmarshal(msg.Data, &req); err != nil {
		respond(msg, false, 422, "invalid request body", nil)
		return
	}

	user := &models.User{
		Username:       req.Username,
		Email:          req.Email,
		Password:       req.Password,
		PhoneNumber:    req.PhoneNumber,
		IdentityNumber: req.IdentityNumber,
	}

	// Validate user data
	if err := handlers.ValidateStruct(*user); err != nil {
		respond(msg, false, 422, "error validation user", err)
		return
	}

	// Check if user already exists
	userExist := new(models.User)
	if err := models.FindUserByIdentity(userExist, user.Username, user.Email, user.PhoneNumber, user.IdentityNumber); err == nil {
		respond(msg, false, 409, "Account Already Exist", userExist)
		return
	}

	// Hash password
	user.PasswordHash = handlers.GeneratePasswordHash(user.Password)

	// Set account status if active flag is true
	if req.Active {
		user.StatusAccount = 10
	}

	// Create user in database
	if err := models.CreateUser(user); err.Error != nil {
		respond(msg, false, 409, err.Error.Error(), nil)
		return
	}

	// Send notification if requested
	if req.SendNotif {
		QueueUserInformation(&req.AccountId, user, req.SendPass)
	}

	respond(msg, true, 201, "Success Create User", user)
}

// FORMAT UPDATE USER
//
//	{
//	  "user_id": "12345",
//	  "user": {
//	    "username": "newusername",
//	    "email": "newemail@example.com",
//	    "phone_number": "1234567890"
//	  }
//	}
func handleUpdateUser(msg *nats.Msg) {
	var req struct {
		UserId string      `json:"user_id"`
		User   models.User `json:"user"`
	}

	if err := json.Unmarshal(msg.Data, &req); err != nil {
		respond(msg, false, 422, "invalid request body", nil)
		return
	}

	if req.UserId == "" {
		respond(msg, false, 422, "need userId params", nil)
		return
	}

	if err := handlers.ValidateStruct(req.User); err != nil {
		respond(msg, false, 422, "error validation user", err)
		return
	}

	if err := models.UpdateUser(req.UserId, &req.User).Error; err != nil {
		respond(msg, false, 500, err.Error(), nil)
		return
	}

	respond(msg, true, 200, "success to update user", nil)
}

func jsErr(c fiber.Ctx, status int, err error) error {
	return c.Status(status).JSON(fiber.Map{"error": err.Error()})
}

// POST /jetstream/streams/:stream/purge?subject=whatsapp.send
// Deletes all messages in the stream, or only one subject if given.
func PurgeStreamHandler(c fiber.Ctx) error {
	js, err := broker.NatsClient.JetStream()
	if err != nil {
		return jsErr(c, fiber.StatusServiceUnavailable, err)
	}

	var opts *nats.StreamPurgeRequest
	if subj := c.Query("subject"); subj != "" {
		opts = &nats.StreamPurgeRequest{Subject: subj}
	}

	if opts != nil {
		err = js.PurgeStream(c.Params("stream"), opts)
	} else {
		err = js.PurgeStream(c.Params("stream"))
	}
	if err != nil {
		return jsErr(c, fiber.StatusInternalServerError, err)
	}
	return c.JSON(fiber.Map{"status": "purged"})
}

// POST /jetstream/streams/:stream/messages/:seq/retry
// Republishes a stuck message (one that hit MaxDeliver) as a fresh message.
// Add ?delete=true to remove the old copy afterwards.
func RetryMessageHandler(c fiber.Ctx) error {
	js, err := broker.NatsClient.JetStream()
	if err != nil {
		return jsErr(c, fiber.StatusServiceUnavailable, err)
	}

	stream := c.Params("stream")
	seq, err := strconv.ParseUint(c.Params("seq"), 10, 64)
	if err != nil {
		return jsErr(c, fiber.StatusBadRequest, err)
	}

	raw, err := js.GetMsg(stream, seq)
	if err != nil {
		return jsErr(c, fiber.StatusNotFound, err)
	}

	ack, err := js.PublishMsg(&nats.Msg{
		Subject: raw.Subject,
		Data:    raw.Data,
		Header:  raw.Header,
	})
	if err != nil {
		return jsErr(c, fiber.StatusInternalServerError, err)
	}

	if c.Query("delete") == "true" {
		_ = js.DeleteMsg(stream, seq)
	}
	return c.JSON(fiber.Map{"status": "requeued", "new_seq": ack.Sequence})
}

// DELETE /jetstream/streams/:stream/messages/:seq
func DeleteMessageHandler(c fiber.Ctx) error {
	js, err := broker.NatsClient.JetStream()
	if err != nil {
		return jsErr(c, fiber.StatusServiceUnavailable, err)
	}

	seq, err := strconv.ParseUint(c.Params("seq"), 10, 64)
	if err != nil {
		return jsErr(c, fiber.StatusBadRequest, err)
	}
	if err := js.DeleteMsg(c.Params("stream"), seq); err != nil {
		return jsErr(c, fiber.StatusInternalServerError, err)
	}
	return c.JSON(fiber.Map{"status": "deleted"})
}

// POST /jetstream/streams/:stream/consumers/:consumer/reset?from=new|all
// Recreates the consumer with the same config. "new" (default) skips the
// whole backlog, "all" replays every message in the stream.
func ResetConsumerHandler(c fiber.Ctx) error {
	js, err := broker.NatsClient.JetStream()
	if err != nil {
		return jsErr(c, fiber.StatusServiceUnavailable, err)
	}

	stream, name := c.Params("stream"), c.Params("consumer")

	info, err := js.ConsumerInfo(stream, name)
	if err != nil {
		return jsErr(c, fiber.StatusNotFound, err)
	}

	cfg := info.Config
	if c.Query("from") == "all" {
		cfg.DeliverPolicy = nats.DeliverAllPolicy
	} else {
		cfg.DeliverPolicy = nats.DeliverNewPolicy
	}
	cfg.OptStartSeq = 0
	cfg.OptStartTime = nil

	if err := js.DeleteConsumer(stream, name); err != nil {
		return jsErr(c, fiber.StatusInternalServerError, err)
	}
	if _, err := js.AddConsumer(stream, &cfg); err != nil {
		return jsErr(c, fiber.StatusInternalServerError, err)
	}
	return c.JSON(fiber.Map{"status": "consumer reset"})
}
