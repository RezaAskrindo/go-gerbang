package services

import (
	"fmt"

	"go-gerbang/handlers"
	"go-gerbang/models"
	"go-gerbang/types"

	"github.com/gofiber/fiber/v3"
)

// baseURL := "https://api2.callmebot.com/text.php"
// 		params := map[string]string{
// 			"user": "@Reza_aceh",
// 			"text": "config email not setup for " + list.Sender,
// 		}
// 		fullURL, err := handlers.BuildURL(baseURL, params)
// 		if err != nil {
// 			log.Println("Error building URL:", err)
// 		}
// 		SendGetRequest(fullURL)

// func SendGetRequest(url string) {
// 	req := fasthttp.AcquireRequest()
// 	req.SetRequestURI(url)
// 	req.Header.SetMethod(fasthttp.MethodGet)
// 	resp := fasthttp.AcquireResponse()
// 	readTimeout, _ := time.ParseDuration("500ms")
// 	writeTimeout, _ := time.ParseDuration("500ms")
// 	maxIdleConnDuration, _ := time.ParseDuration("1h")
// 	client := &fasthttp.Client{
// 		ReadTimeout:                   readTimeout,
// 		WriteTimeout:                  writeTimeout,
// 		MaxIdleConnDuration:           maxIdleConnDuration,
// 		NoDefaultUserAgentHeader:      true, // Don't send: User-Agent: fasthttp
// 		DisableHeaderNamesNormalizing: true, // If you set the case on your headers correctly you can enable this
// 		DisablePathNormalizing:        true,
// 		// increase DNS cache time to an hour instead of default minute
// 		Dial: (&fasthttp.TCPDialer{
// 			Concurrency:      4096,
// 			DNSCacheDuration: time.Hour,
// 		}).Dial,
// 	}
// 	err := client.Do(req, resp)
// 	fasthttp.ReleaseRequest(req)
// 	if err != nil {
// 		log.Printf("ERR HTTP Connection error: %v\n", err)
// 	}
// 	fasthttp.ReleaseResponse(resp)
// }

func MailTesting(c fiber.Ctx) error {
	to := c.Query("to")
	if to == "" {
		return handlers.InternalServerErrorResponse(c, fmt.Errorf("need params to"))
	}

	appName := c.Query("appName")
	if appName == "" {
		return handlers.InternalServerErrorResponse(c, fmt.Errorf("need params appName"))
	}

	provider := c.Query("provider")
	if provider == "" {
		return handlers.InternalServerErrorResponse(c, fmt.Errorf("need params provider"))
	}

	if !handlers.IsValidEmail(to) {
		return handlers.InternalServerErrorResponse(c, fmt.Errorf("please provide valid email"))
	}

	dataSend := &types.ListEmail{
		Sender:           appName,
		Subject:          "Testing Email!",
		BodyTemplateText: "Testing email...",
		BodyTemplateHtml: "<p>Testing email...</p>",
		Emails: []types.Email{
			{Name: "Test User", EmailAddr: to},
		},
	}

	tipe := c.Query("type")
	if tipe == "event" {
		var sendToEvent types.SendingEmailToBroker
		sendToEvent = types.SendingEmailToBroker{
			Sender:   dataSend.Sender,
			Provider: provider,
			Subject:  dataSend.Subject,
			Title:    dataSend.Subject,
			BodyText: dataSend.BodyTemplateText,
			Body:     dataSend.BodyTemplateHtml,
			Footer:   "",
			Emails:   dataSend.Emails,
		}
		PublishEvent("user.notification", sendToEvent)
		return handlers.SuccessResponse(c, true, "Send Mail On Event Success", nil, nil)
	}

	if provider == "Resend" {
		if !handlers.SendResendMail(dataSend) {
			return handlers.InternalServerErrorResponse(c, fmt.Errorf("failed to send email using "+provider))
		}
	} else if provider == "SMTP" && !handlers.SendSMTPMail(dataSend) {
		return handlers.InternalServerErrorResponse(c, fmt.Errorf("failed to send email using "+provider))
	} else {
		return handlers.SuccessResponse(c, true, "Send Mail Not On Configuration", nil, nil)
	}

	return handlers.SuccessResponse(c, true, "Check Mail Success", nil, nil)
}

func QueueUserInformation(providerNotification string, querySender string, user *models.User, sendPass bool) bool {
	Sender := "GOGERBANG"
	if querySender != "" {
		Sender = querySender
	}

	textPass := ``
	htmlPass := ``
	if sendPass {
		textPass = `password: ` + user.Password
		htmlPass = `<div>password: <strong>` + user.Password + `</strong></div>`
	}

	sendEmail := new(types.SendingEmailToBroker)
	sendEmail.Sender = Sender
	sendEmail.Provider = providerNotification
	sendEmail.Subject = "Create Account Success"
	sendEmail.Title = "Akun Anda Berhasil Di Buat"
	sendEmail.BodyText = `
		Hi, ` + user.FullName + `, berikut informasi akun anda:

		username: ` + user.Username + `
		email: ` + user.Email + textPass + `

		Tetap jaga rahasia akun anda, mohon untuk jangan diberikan kepada siapapun termasuk Admin.
	`
	sendEmail.Body = `
		<div class="font-family: Roboto-Regular, Helvetica, Arial, sans-serif; font-size: 14px; color: rgba(0, 0, 0, 0.87); padding-top: 20px; text-align: center;">Hi, ` + user.FullName + `, berikut informasi akun anda:</div>
		<br/>
		<div>username: <strong>` + user.FullName + `</strong></div>
		<div>email: <strong>` + user.Email + `</strong></div>` + htmlPass + `
		<br/>
		<div class="padding-top: 20px; font-size: 12px; line-height: 16px; color: rgb(95, 99, 104); letter-spacing: 0.3px; text-align: center;">Tetap jaga rahasia akun anda, mohon untuk jangan diberikan kepada siapapun termasuk Admin.</div>
	`
	sendEmail.Footer = "ini merupakan email otomatis dari " + Sender
	sendEmail.Emails = []types.Email{
		{
			Name:      user.FullName,
			EmailAddr: user.Email,
		},
	}

	PublishEvent("user.notification", sendEmail)

	return true
}
