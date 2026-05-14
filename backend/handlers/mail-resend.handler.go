package handlers

import (
	"context"
	"fmt"
	"time"

	"go-gerbang/types"

	"github.com/resend/resend-go/v2"
	"go.uber.org/zap"
)

var resendEmailStat = "send-email-resend"

func GetApiResendKey(sender string) *types.ResendKey {
	apiResend := GetEmailResendConfig()
	for _, key := range apiResend {
		if key.Sender == sender {
			return &key
		}
	}
	return nil
}

var ResendClient *resend.Client

func ExtractEmailAddrs(listEmail types.ListEmail) []string {
	var emailAddrs []string
	for _, email := range listEmail.Emails {
		emailAddrs = append(emailAddrs, email.EmailAddr)
	}
	return emailAddrs
}

func SendResendMail(list *types.ListEmail) bool {
	start := time.Now()

	APIResentKey := GetApiResendKey(list.Sender)

	if APIResentKey == nil {
		duration := time.Since(start)
		ZapLogger.Error(EmailError,
			zap.String("path", resendEmailStat),
			zap.Int("status", EmailErrorCode),
			zap.Duration("duration", duration),
			zap.Error(fmt.Errorf(list.Sender+EmailErrorNotInDB)),
		)
		return false
	}

	if ResendClient == nil {
		ResendClient = resend.NewClient(APIResentKey.Key)
	}

	emailAddrs := ExtractEmailAddrs(*list)
	var batchEmails []*resend.SendEmailRequest

	if list.TypeBatchAddress == "all" {
		params := &resend.SendEmailRequest{
			From:    list.Sender + " <" + APIResentKey.Email + ">",
			To:      emailAddrs,
			Subject: list.Subject,
			Text:    list.BodyTemplateText,
			Html:    list.BodyTemplateHtml,
		}

		_, err := ResendClient.Emails.Send(params)
		if err != nil {
			duration := time.Since(start)
			ZapLogger.Error(EmailError,
				zap.String("path", resendEmailStat),
				zap.Int("status", EmailErrorCode),
				zap.Duration("duration", duration),
				zap.Error(err),
			)
			return false
		}
	} else {
		ctx := context.TODO()

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
			duration := time.Since(start)
			ZapLogger.Error(EmailError,
				zap.String("path", resendEmailStat),
				zap.Int("status", EmailErrorCode),
				zap.Duration("duration", duration),
				zap.Error(err),
			)
			return false
		}

	}

	duration := time.Since(start)
	ZapLogger.Info(EmailSuccess,
		zap.String("path", resendEmailStat),
		zap.Int("status", EmailSuccessCode),
		zap.Duration("duration", duration),
		zap.Any("request", batchEmails),
		zap.Any("request", emailAddrs),
	)

	return true
}
