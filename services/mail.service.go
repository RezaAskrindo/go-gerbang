package services

import (
	"context"
	"fmt"
	ht "html/template"
	"log"
	"regexp"
	tt "text/template"

	"go-gerbang/config"
	"go-gerbang/handlers"
	"go-gerbang/types"

	"github.com/gofiber/fiber/v2"
	"github.com/resend/resend-go/v2"
	"github.com/wneessen/go-mail"
)

func CreateMailClient(appName string) (*mail.Client, error) {
	CONFIG_SMTP_HOST := config.Config(appName + "_CONFIG_SMTP_HOST")
	CONFIG_AUTH_EMAIL := config.Config(appName + "_CONFIG_AUTH_EMAIL")
	CONFIG_AUTH_PASSWORD := config.Config(appName + "_CONFIG_AUTH_PASSWORD")
	return mail.NewClient(CONFIG_SMTP_HOST,
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithUsername(CONFIG_AUTH_EMAIL),
		mail.WithPassword(CONFIG_AUTH_PASSWORD),
	)
}

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

	appName := c.Query("appName")
	if appName == "" {
		return handlers.InternalServerErrorResponse(c, fmt.Errorf("need params appName"))
	}

	if !isValidEmail(to) {
		return handlers.InternalServerErrorResponse(c, fmt.Errorf("please provide valid email"))
	}

	// from := config.Config(appName + "_CONFIG_AUTH_EMAIL")

	// if !isValidEmail(from) {
	// 	baseURL := "https://api2.callmebot.com/text.php"
	// 	params := map[string]string{
	// 		"user": "@Reza_aceh",
	// 		"text": "config email not setup for " + appName,
	// 	}
	// 	fullURL, err := handlers.BuildURL(baseURL, params)
	// 	if err != nil {
	// 		log.Println("Error building URL:", err)
	// 	}
	// 	SendGetRequest(fullURL)
	// 	return handlers.InternalServerErrorResponse(c, fmt.Errorf("config email not setup for "+appName))
	// }

	// message := mail.NewMsg()
	// if err := message.From(config.Config(appName + "_CONFIG_AUTH_EMAIL")); err != nil {
	// 	log.Printf("failed to set From address: %s", err)
	// 	return handlers.InternalServerErrorResponse(c, err)
	// }
	// if err := message.To(to); err != nil {
	// 	log.Printf("failed to set To address: %s", err)
	// 	return handlers.InternalServerErrorResponse(c, err)
	// }

	// message.Subject("Testing Email!")
	// message.SetBodyString(mail.TypeTextPlain, "Testing email...")
	// client, err := CreateMailClient(appName)
	// if err != nil {
	// 	log.Printf("failed to create mail client: %s", err)
	// }
	// if err := client.DialAndSend(message); err != nil {
	// 	log.Printf("failed to send mail: %s", err)
	// }

	send := SendResendMail(&types.ListEmail{
		Sender:           appName,
		Subject:          "Testing Email!",
		BodyTemplateText: "Testing email...",
		BodyTemplateHtml: "<p>Testing email...</p>",
		Emails: []types.Email{
			{Name: "Test User", EmailAddr: to},
		},
	})

	if !send {
		log.Printf("failed to send mail using Resend: %s", "error sending email")
		return handlers.InternalServerErrorResponse(c, fmt.Errorf("failed to send email using Resend"))
	}

	return handlers.SuccessResponse(c, true, "Check Mail Success", nil, nil)
}

func SendMail(list *types.ListEmail) bool {
	userList := list.Emails

	textTpl, err := tt.New("texttpl").Parse(list.BodyTemplateText)
	if err != nil {
		log.Printf("failed to parse text template: %s", err)
	}

	htmlTpl, err := ht.New("htmltpl").Parse(list.BodyTemplateHtml)
	if err != nil {
		log.Printf("failed to parse text template: %s", err)
	}

	from := config.Config(list.Sender + "_CONFIG_AUTH_EMAIL")
	if !isValidEmail(from) {
		baseURL := "https://api2.callmebot.com/text.php"
		params := map[string]string{
			"user": "@Reza_aceh",
			"text": "config email not setup for " + list.Sender,
		}
		fullURL, err := handlers.BuildURL(baseURL, params)
		if err != nil {
			log.Println("Error building URL:", err)
		}
		SendGetRequest(fullURL)
	} else {
		var messages []*mail.Msg
		for _, user := range userList {
			message := mail.NewMsg()
			if err := message.EnvelopeFrom(config.Config(list.Sender + "_CONFIG_AUTH_EMAIL")); err != nil {
				log.Printf("failed to set ENVELOPE FROM address: %s", err)
			}
			if err := message.FromFormat(list.Sender, config.Config(list.Sender+"_CONFIG_AUTH_EMAIL")); err != nil {
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

		// client, err := mail.NewClient(config.Config("CONFIG_SMTP_HOST"),
		// 	mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithTLSPortPolicy(mail.TLSMandatory),
		// 	mail.WithUsername(config.Config("CONFIG_AUTH_EMAIL")), mail.WithPassword(config.Config("CONFIG_AUTH_PASSWORD")),
		// )
		client, _ := CreateMailClient(list.Sender)
		if err := client.DialAndSend(messages...); err != nil {
			log.Printf("failed to deliver mail: %s", err)
		}
	}

	return true
}

// RESEND API
var APIResent []types.ResendKey = []types.ResendKey{
	{Sender: "SISKOR", Key: "re_DByKnU6f_PXKdrRBzANLfZhE7Yn5i7Kmq", Email: "notification@siskor.web.id"},
}

func GetApiResendKey(sender string) *types.ResendKey {
	for _, key := range APIResent {
		if key.Sender == sender {
			return &key
		}
	}
	return nil
}

var ResendClient *resend.Client

func extractEmailAddrs(listEmail types.ListEmail) []string {
	var emailAddrs []string
	for _, email := range listEmail.Emails {
		emailAddrs = append(emailAddrs, email.EmailAddr)
	}
	return emailAddrs
}

func SendResendMail(list *types.ListEmail) bool {
	APIResentKey := GetApiResendKey(list.Sender)

	if APIResentKey == nil {
		return false
	}

	if ResendClient == nil {
		ResendClient = resend.NewClient(APIResentKey.Key)
	}

	if list.TypeBatchAddress == "all" {
		emailAddrs := extractEmailAddrs(*list)

		params := &resend.SendEmailRequest{
			From:    list.Sender + " <" + APIResentKey.Email + ">",
			To:      emailAddrs,
			Subject: list.Subject,
			Text:    list.BodyTemplateText,
			Html:    list.BodyTemplateHtml,
		}

		_, err := ResendClient.Emails.Send(params)
		if err != nil {
			log.Printf("failed to send email using Resend: %s", err)
			return false
		}
	} else {
		ctx := context.TODO()

		var batchEmails []*resend.SendEmailRequest
		for _, email := range list.Emails {
			req := &resend.SendEmailRequest{
				From:    list.Sender + " <" + APIResentKey.Email + ">",
				To:      []string{email.EmailAddr},
				Subject: list.Subject,
				Text:    list.BodyTemplateText,
				Html:    list.BodyTemplateHtml,
			}
			batchEmails = append(batchEmails, req)
		}

		_, err := ResendClient.Batch.SendWithContext(ctx, batchEmails)
		if err != nil {
			log.Printf("failed to send email using Resend: %s", err)
			return false
		}
	}

	return true
}
