package services

import (
	"encoding/json"
	"fmt"
	ht "html/template"
	"log"
	"net/http"
	"regexp"
	tt "text/template"

	"go-gerbang/config"
	"go-gerbang/handlers"
	"go-gerbang/types"

	"github.com/gofiber/fiber/v2"
	"github.com/wneessen/go-mail"
)

// var CONFIG_SMTP_PORT = config.Config("CONFIG_SMTP_PORT")
// var CONFIG_SENDER_NAME = config.Config("CONFIG_SENDER_NAME")
var CONFIG_SMTP_HOST = config.Config("CONFIG_SMTP_HOST")
var CONFIG_AUTH_EMAIL = config.Config("CONFIG_AUTH_EMAIL")
var CONFIG_AUTH_PASSWORD = config.Config("CONFIG_AUTH_PASSWORD")

func isValidEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)

	return re.MatchString(email)
}

func MailTesting(c *fiber.Ctx) error {
	to := c.Query("to")

	if to == "" {
		return handlers.InternalServerErrorResponse(c, fmt.Errorf("need params to"))
	}

	if !isValidEmail(to) {
		return handlers.InternalServerErrorResponse(c, fmt.Errorf("please provide valid email"))
	}

	message := mail.NewMsg()
	if err := message.From(config.Config("CONFIG_AUTH_EMAIL")); err != nil {
		log.Printf("failed to set From address: %s", err)
		return handlers.InternalServerErrorResponse(c, err)
	}
	if err := message.To(to); err != nil {
		log.Printf("failed to set To address: %s", err)
		return handlers.InternalServerErrorResponse(c, err)
	}

	message.Subject("Testing Email!")
	message.SetBodyString(mail.TypeTextPlain, "Testing email...")
	client, err := mail.NewClient(config.Config("CONFIG_SMTP_HOST"),
		mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithUsername(config.Config("CONFIG_AUTH_EMAIL")), mail.WithPassword(config.Config("CONFIG_AUTH_PASSWORD")))
	if err != nil {
		log.Printf("failed to create mail client: %s", err)
	}
	if err := client.DialAndSend(message); err != nil {
		log.Printf("failed to send mail: %s", err)
	}

	return handlers.SuccessResponse(c, true, "Check Mail Success", nil, nil)
}

func SendMail(list *types.ListEmail) bool {
	userList := list.Emails

	log.Println("bulk mailing begin.")

	textTpl, err := tt.New("texttpl").Parse(list.BodyTemplateText)
	if err != nil {
		log.Printf("failed to parse text template: %s", err)
	}

	htmlTpl, err := ht.New("htmltpl").Parse(list.BodyTemplateHtml)
	if err != nil {
		log.Printf("failed to parse text template: %s", err)
	}

	var messages []*mail.Msg
	for _, user := range userList {
		message := mail.NewMsg()
		if err := message.EnvelopeFrom(config.Config("CONFIG_AUTH_EMAIL")); err != nil {
			log.Printf("failed to set ENVELOPE FROM address: %s", err)
		}
		if err := message.FromFormat(list.Sender, config.Config("CONFIG_AUTH_EMAIL")); err != nil {
			log.Printf("failed to set formatted FROM address: %s", err)
		}
		if err := message.AddToFormat(user.Name, user.EmailAddr); err != nil {
			log.Printf("failed to set formatted TO address: %s", err)
		}
		message.SetMessageID()
		message.SetDate()
		message.SetBulk()
		message.Subject(list.Subject)
		if err := message.SetBodyTextTemplate(textTpl, user); err != nil {
			log.Printf("failed to add text template to mail body: %s", err)
		}
		if err := message.AddAlternativeHTMLTemplate(htmlTpl, user); err != nil {
			log.Printf("failed to add HTML template to mail body: %s", err)
		}

		messages = append(messages, message)
	}

	client, err := mail.NewClient(config.Config("CONFIG_SMTP_HOST"),
		mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithUsername(config.Config("CONFIG_AUTH_EMAIL")), mail.WithPassword(config.Config("CONFIG_AUTH_PASSWORD")),
	)
	if err := client.DialAndSend(messages...); err != nil {
		log.Printf("failed to deliver mail: %s", err)
	}
	log.Println("bulk mailing successfully delivered.")

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

// NOT USED ANYMORE BECAUSE USING DIRECT FUNCTION TO SUBCRIBE
// func SendEmailHandler(c *fiber.Ctx) error {
// 	u := new(types.ListEmail)

// 	if err := handlers.ParseBody(c, u); err != nil {
// 		return handlers.BadRequestErrorResponse(c, err)
// 	}

// 	if !SendMail(u) {
// 		return handlers.InternalServerErrorResponse(c, fmt.Errorf("failed"))
// 	}

// 	return handlers.SuccessResponse(c, true, "success send email", nil, nil)
// }
