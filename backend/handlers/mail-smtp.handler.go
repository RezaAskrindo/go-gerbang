package handlers

import (
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"go-gerbang/types"

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
			zap.Error(fmt.Errorf(list.Sender+EmailErrorNotInDB)),
		)
		return false
	}

	emailAddrs := ExtractEmailAddrs(*list)

	// subject := fmt.Sprintf("Subject: %s \n", list.Subject)
	// body := fmt.Sprintf("Your verification code is %s", list.BodyTemplateHtml)
	// message := []byte(subject + "\n" + body)

	boundary := "mixed-boundary-123456"
	altBoundary := "alt-boundary-123456"

	header := ""
	header += fmt.Sprintf("From: %s\r\n", list.Sender)
	header += fmt.Sprintf("To: %s\r\n", strings.Join(emailAddrs, ","))
	header += fmt.Sprintf("Subject: %s\r\n", list.Subject)
	header += "MIME-Version: 1.0\r\n"
	header += fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", boundary)
	header += "\r\n"

	body := ""

	// Alternative part (text + html)
	body += fmt.Sprintf("--%s\r\n", boundary)
	body += fmt.Sprintf("Content-Type: multipart/alternative; boundary=%s\r\n\r\n", altBoundary)

	// Plain text
	body += fmt.Sprintf("--%s\r\n", altBoundary)
	body += "Content-Type: text/plain; charset=UTF-8\r\n\r\n"
	body += list.BodyTemplateText + "\r\n"

	// HTML
	body += fmt.Sprintf("--%s\r\n", altBoundary)
	body += "Content-Type: text/html; charset=UTF-8\r\n\r\n"
	body += list.BodyTemplateHtml + "\r\n"

	body += fmt.Sprintf("--%s--\r\n", altBoundary)

	// FOR FUTURE WORK
	// for _, file := range email.Attachments {

	// 	data, err := os.ReadFile(file)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	filename := filepath.Base(file)

	// 	body += fmt.Sprintf("--%s\r\n", boundary)
	// 	body += fmt.Sprintf("Content-Type: application/octet-stream\r\n")
	// 	body += "Content-Transfer-Encoding: base64\r\n"
	// 	body += fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n\r\n", filename)

	// 	b := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	// 	base64.StdEncoding.Encode(b, data)

	// 	body += string(b) + "\r\n"
	// }

	body += fmt.Sprintf("--%s--", boundary)

	message := []byte(header + body)

	auth := smtp.PlainAuth("", s.SMTPUser, s.SMTPPassword, s.SMTPHost)
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", s.SMTPHost, s.SMTPPort),
		auth,
		s.SMTPUser,
		emailAddrs,
		message,
	)
	if err != nil {
		duration := time.Since(start)
		ZapLogger.Error(EmailError,
			zap.String("path", smtpEmailStat),
			zap.Int("status", EmailErrorCode),
			zap.Duration("duration", duration),
			zap.Error(err),
		)
		return false
	}

	duration := time.Since(start)
	ZapLogger.Info(EmailSuccess,
		zap.String("path", smtpEmailStat),
		zap.Int("status", EmailSuccessCode),
		zap.Duration("duration", duration),
		zap.Any("request", emailAddrs),
	)
	return true
}

type SMTPService struct {
	config *types.SMTPConfig
}

func NewSMTPService(config *types.SMTPConfig) *SMTPService {
	return &SMTPService{config: config}
}

func (s *SMTPService) SendVerificationCode(to string, code string) error {
	subject := "Subject: Email Verification Code \n"
	body := fmt.Sprintf("Your verification code is %s", code)
	message := []byte(subject + "\n" + body)

	auth := smtp.PlainAuth("", s.config.SMTPUser, s.config.SMTPPassword, s.config.SMTPHost)
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort),
		auth,
		s.config.SMTPUser,
		[]string{to},
		message,
	)
	if err != nil {
		return err
	}

	return nil
}
