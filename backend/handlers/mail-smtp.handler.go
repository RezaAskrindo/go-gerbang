package handlers

import (
	"bytes"
	"fmt"
	"time"

	"go-gerbang/types"

	mail "github.com/wneessen/go-mail"
	"go.uber.org/zap"
)

var smtpEmailStat = "send-email-smtp"

func GetSMTPSetup(sender string) *types.SMTPConfig {
	mailSMTP := GetEmailSMTPConfig()
	for _, user := range mailSMTP {
		if user.Sender == sender {
			return &user
		}
	}
	return nil
}

func SendSMTPMail(list *types.ListEmail) bool {
	start := time.Now()

	s := GetSMTPSetup(list.Sender)

	if s == nil {
		duration := time.Since(start)
		ZapLogger.Error(EmailError,
			zap.String("path", smtpEmailStat),
			zap.Int("status", EmailErrorCode),
			zap.Duration("duration", duration),
			zap.Error(fmt.Errorf("%s", list.Sender+EmailErrorNotInDB)),
		)
		return false
	}

	emailAddrs := ExtractEmailAddrs(*list)

	msg := mail.NewMsg()

	if err := msg.FromFormat(list.Sender, s.SMTPUser); err != nil {
		logSMTPError(start, err)
		return false
	}

	if err := msg.To(emailAddrs...); err != nil {
		logSMTPError(start, err)
		return false
	}

	msg.Subject(list.Subject)

	// Plain text body
	if list.BodyTemplateText != "" {
		msg.SetBodyString(mail.TypeTextPlain, list.BodyTemplateText)
	}

	// HTML alternative
	if list.BodyTemplateHtml != "" {
		msg.AddAlternativeString(mail.TypeTextHTML, list.BodyTemplateHtml)
	}

	// Optional attachments
	// Assuming ListEmail has:
	// Attachments []string
	for _, att := range list.Attachments {
		opts := []mail.FileOption{}

		if att.ContentType != "" {
			opts = append(opts, mail.WithFileContentType(mail.ContentType(att.ContentType)))
		}

		msg.AttachReadSeeker(att.Name, bytes.NewReader(att.Data), opts...)
	}

	client, err := mail.NewClient(
		s.SMTPHost,
		mail.WithPort(s.SMTPPort),
		mail.WithSSL(),
		mail.WithUsername(s.SMTPUser),
		mail.WithPassword(s.SMTPPassword),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithTimeout(10*time.Second),
	)
	if err != nil {
		logSMTPError(start, err)
		return false
	}

	if err := client.DialAndSend(msg); err != nil {
		logSMTPError(start, err)
		return false
	}

	duration := time.Since(start)
	ZapLogger.Info(EmailSuccess,
		zap.String("path", smtpEmailStat),
		zap.Int("status", EmailSuccessCode),
		zap.Duration("duration", duration),
		zap.Strings("request", emailAddrs),
	)

	return true
}

func logSMTPError(start time.Time, err error) {
	duration := time.Since(start)

	ZapLogger.Error(EmailError,
		zap.String("path", smtpEmailStat),
		zap.Int("status", EmailErrorCode),
		zap.Duration("duration", duration),
		zap.Error(err),
	)
}

type SMTPService struct {
	config *types.SMTPConfig
}

func NewSMTPService(config *types.SMTPConfig) *SMTPService {
	return &SMTPService{config: config}
}
