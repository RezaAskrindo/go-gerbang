package services

import (
	"encoding/json"
	"fmt"
	ht "html/template"
	"log"
	"net/http"
	tt "text/template"

	"go-gerbang/config"
	"go-gerbang/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/wneessen/go-mail"
)

var CONFIG_SMTP_HOST = config.Config("CONFIG_SMTP_HOST")
var CONFIG_SMTP_PORT = config.Config("CONFIG_SMTP_PORT")
var CONFIG_SENDER_NAME = config.Config("CONFIG_SENDER_NAME")
var CONFIG_AUTH_EMAIL = config.Config("CONFIG_AUTH_EMAIL")
var CONFIG_AUTH_PASSWORD = config.Config("CONFIG_AUTH_PASSWORD")

type Email struct {
	Name      string `json:"name"`
	EmailAddr string `json:"email_addr"`
}

type ListEmail struct {
	Sender           string  `json:"sender"`
	Subject          string  `json:"subject"`
	BodyTemplateText string  `json:"body_template_text"`
	BodyTemplateHtml string  `json:"body_template_html"`
	Emails           []Email `json:"emails"`
}

func SendMail(list *ListEmail) bool {
	userList := list.Emails

	textTpl, err := tt.New("texttpl").Parse(list.BodyTemplateText)
	if err != nil {
		log.Fatalf("failed to parse text template: %s", err)
	}
	htmlTpl, err := ht.New("htmltpl").Parse(list.BodyTemplateHtml)
	if err != nil {
		log.Fatalf("failed to parse text template: %s", err)
	}

	var messages []*mail.Msg
	for _, user := range userList {
		message := mail.NewMsg()
		if err := message.EnvelopeFrom(config.Config("CONFIG_AUTH_EMAIL")); err != nil {
			log.Fatalf("failed to set ENVELOPE FROM address: %s", err)
		}
		if err := message.FromFormat(list.Sender, config.Config("CONFIG_AUTH_EMAIL")); err != nil {
			log.Fatalf("failed to set formatted FROM address: %s", err)
		}
		if err := message.AddToFormat(user.Name, user.EmailAddr); err != nil {
			log.Fatalf("failed to set formatted TO address: %s", err)
		}
		message.SetMessageID()
		message.SetDate()
		message.SetBulk()
		message.Subject(list.Subject)
		if err := message.SetBodyTextTemplate(textTpl, user); err != nil {
			log.Fatalf("failed to add text template to mail body: %s", err)
		}
		if err := message.AddAlternativeHTMLTemplate(htmlTpl, user); err != nil {
			log.Fatalf("failed to add HTML template to mail body: %s", err)
		}

		messages = append(messages, message)
	}

	client, err := mail.NewClient(config.Config("CONFIG_SMTP_HOST"),
		mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithUsername(config.Config("CONFIG_AUTH_EMAIL")), mail.WithPassword(config.Config("CONFIG_AUTH_PASSWORD")),
	)
	if err := client.DialAndSend(messages...); err != nil {
		log.Fatalf("failed to deliver mail: %s", err)
	}
	log.Printf("Bulk mailing successfully delivered.")

	return err == nil
}

type Response struct {
	Response interface{} `json:"response"`
	Status   int         `json:"status"`
}

func JsonResponseHandler(response Response, statusCode int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
		}
	}
}

func SendEmailHandler(c *fiber.Ctx) error {
	u := new(ListEmail)

	if err := handlers.ParseBody(c, u); err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	if !SendMail(u) {
		return handlers.InternalServerErrorResponse(c, fmt.Errorf("failed"))
	}

	return handlers.SuccessResponse(c, true, "success send email", nil, nil)
}
