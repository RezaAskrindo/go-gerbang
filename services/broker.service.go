package services

import (
	"encoding/json"
	"log"

	"go-gerbang/broker"
	"go-gerbang/handlers"
	"go-gerbang/types"

	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
)

func PublishService(c *fiber.Ctx) error {
	input := new(types.SendingEmailToBroker)
	if err := handlers.ParseBody(c, input); err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	messageBytes, err := json.Marshal(input)
	if err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	subject := "send_mail"
	if err := broker.NatsClient.Publish(subject, []byte(messageBytes)); err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	return handlers.SuccessResponse(c, true, "on publish", nil, nil)
}

func SubscribeService(c *fiber.Ctx) error {
	subject := "send_mail"
	// broker.NatsClient.Subscribe(subject, func(msg *nats.Msg) {
	// 	fmt.Printf("Received message on %s: %s\n", subject, string(msg.Data))
	// })

	return handlers.SuccessResponse(c, true, "on subscribe", subject, nil)
}

// func PublishServiceEmail(input *types.SendingEmailToBroker) error {
// 	messageBytes, err := json.Marshal(input)
// 	if err != nil {
// 		return err
// 	}

// 	subject := "send_mail"
// 	if err := broker.NatsClient.Publish(subject, []byte(messageBytes)); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func SubscribeServiceEmail() {
// 	subject := "send_mail"
// 	broker.NatsClient.Subscribe(subject, func(msg *nats.Msg) {
// 		var email types.SendingEmailToBroker

// 		if err := json.Unmarshal(msg.Data, &email); err != nil {
// 			log.Println("Failed to parse message data:", err)
// 			return
// 		}

// 		dataSend := new(types.ListEmail)
// 		dataSend.Sender = email.Sender
// 		dataSend.Subject = email.Subject
// 		dataSend.BodyTemplateText = email.Title + `

// 		` + email.BodyText + `

// 		` + email.Footer
// 		dataSend.BodyTemplateHtml = `<div style="margin: 0px; padding: 0px;" bgcolor="#FFFFFF">
// 			<table width="100%" height="100%" style="min-width: 348px;" border="0" cellspacing="0" cellpadding="0" lang="en">
// 				<tbody>
// 					<tr height="32" style="height: 32px;">
// 						<td></td>
// 					</tr>
// 					<tr align="center">
// 						<td>
// 							<table border="0" cellspacing="0" cellpadding="0" style="padding-bottom: 20px; max-width: 516px; min-width: 220px;">
// 								<tbody>
// 									<tr>
// 										<td width="8" style="width: 8px;"></td>
// 										<td>
// 											<div style="border-style: solid; border-width: thin; border-color: rgb(218, 220, 224); border-radius: 8px; padding: 40px 20px;" align="center">
// 												<div style="font-family: Google Sans, Roboto, RobotoDraft, Helvetica, Arial, sans-serif; border-bottom: thin solid rgb(218, 220, 224); color: rgba(0, 0, 0, 0.87); line-height: 32px; padding-bottom: 24px; text-align: center; word-break: break-word;">
// 													<div style="font-size: 24px;">
// 														` + email.Title + `
// 													</div>
// 												</div>
// 												<div style="font-family: Roboto-Regular, Helvetica, Arial, sans-serif; font-size: 14px; color: rgba(0, 0, 0, 0.87); line-height: 20px; padding-top: 20px; text-align: center;">
// 													` + email.Body + `
// 												</div>
// 											</div>
// 											<div style="font-family: Roboto-Regular, Helvetica, Arial, sans-serif; color: rgba(0, 0, 0, 0.54); font-size: 11px; line-height: 18px; padding-top: 12px; text-align: center;">
// 												` + email.Footer + `
// 											</div>
// 										</td>
// 										<td width="8" style="width: 8px;"></td>
// 									</tr>
// 								</tbody>
// 							</table>
// 						</td>
// 					</tr>
// 					<tr height="32" style="height: 32px;">
// 						<td></td>
// 					</tr>
// 				</tbody>
// 			</table>
// 		</div>`
// 		dataSend.Emails = email.Emails

// 		if !SendMail(dataSend) {
// 			log.Println("failed to send email")
// 		}

// 		// req := fasthttp.AcquireRequest()
// 		// res := fasthttp.AcquireResponse()
// 		// defer fasthttp.ReleaseRequest(req)
// 		// defer fasthttp.ReleaseResponse(res)

// 		// req.Header.SetContentType("application/json")
// 		// req.Header.SetMethod("POST")
// 		// data := handlers.ToMarshal(dataSend)
// 		// req.SetBody(data)
// 		// req.SetRequestURI(config.Config("EMAIL_SERVICE") + "/send-email")

// 		// if err := handlers.Client.DoTimeout(req, res, 300*time.Second); err != nil {
// 		// 	log.Printf("error %s\n", err)
// 		// }

// 		log.Println("Success Subscribe Message")
// 	})
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
		"user.notification": handleMsg,
	}

	for subject, handler := range handlers {
		if _, err := broker.NatsClient.Subscribe(subject, handler); err != nil {
			log.Println("Error on subcribe NATS:", err)
		}
	}

	// log.Println("Listening for subcribe events...")
}

func handleMsg(msg *nats.Msg) {
	var email types.SendingEmailToBroker

	if err := json.Unmarshal(msg.Data, &email); err != nil {
		log.Println("Failed to parse message data:", err)
		return
	}

	dataSend := new(types.ListEmail)
	dataSend.Sender = email.Sender
	dataSend.Subject = email.Subject
	dataSend.BodyTemplateText = email.Title + `
	
	` + email.BodyText + `
	
	
	` + email.Footer
	dataSend.BodyTemplateHtml = `<div style="margin: 0px; padding: 0px;" bgcolor="#FFFFFF">
		<table width="100%" height="100%" style="min-width: 348px;" border="0" cellspacing="0" cellpadding="0" lang="en">
			<tbody>
				<tr height="32" style="height: 32px;">
					<td></td>
				</tr>
				<tr align="center">
					<td>
						<table border="0" cellspacing="0" cellpadding="0" style="padding-bottom: 20px; max-width: 516px; min-width: 220px;">
							<tbody>
								<tr>
									<td width="8" style="width: 8px;"></td>
									<td>
										<div style="border-style: solid; border-width: thin; border-color: rgb(218, 220, 224); border-radius: 8px; padding: 40px 20px;" align="center">
											<div style="font-family: Google Sans, Roboto, RobotoDraft, Helvetica, Arial, sans-serif; border-bottom: thin solid rgb(218, 220, 224); color: rgba(0, 0, 0, 0.87); line-height: 32px; padding-bottom: 24px; text-align: center; word-break: break-word;">
												<div style="font-size: 24px;">
													` + email.Title + `
												</div>
											</div>
											<div style="font-family: Roboto-Regular, Helvetica, Arial, sans-serif; font-size: 14px; color: rgba(0, 0, 0, 0.87); line-height: 20px; padding-top: 20px; text-align: center;">
												` + email.Body + `
											</div>
										</div>
										<div style="font-family: Roboto-Regular, Helvetica, Arial, sans-serif; color: rgba(0, 0, 0, 0.54); font-size: 11px; line-height: 18px; padding-top: 12px; text-align: center;">
											` + email.Footer + `
										</div>
									</td>
									<td width="8" style="width: 8px;"></td>
								</tr>
							</tbody>
						</table>
					</td>
				</tr>
				<tr height="32" style="height: 32px;">
					<td></td>
				</tr>
			</tbody>
		</table>
	</div>`
	dataSend.Emails = email.Emails

	if !SendMail(dataSend) {
		log.Println("failed to send email")
	}

	// log.Println("Success Subscribe Message")
}
