package handlers

import (
	"fmt"
	"html/template"
	"regexp"
	"strconv"
	"sync"
	"time"

	"go-gerbang/models"
	"go-gerbang/types"

	"github.com/valyala/bytebufferpool"
)

var (
	emailServiceName             = make(map[string]struct{})
	emailServiceNameLastRefresh  time.Time
	emailServiceUrl              string
	emailServiceUrlLastRefresh   time.Time
	emailResendConfig            = []types.ResendKey{}
	emailResendConfigLastRefresh time.Time
	emailSMTPConfig              = []types.SMTPConfig{}
	emailSMTPConfigLastRefresh   time.Time
	mailMutex                    sync.RWMutex
	cacheTTL                     = 5 * time.Minute
	EmailSuccess                 = "info-mail"
	EmailSuccessCode             = 200
	EmailError                   = "error-mail"
	EmailErrorCode               = 304
	EmailErrorNotInDB            = " (appName) not registered on Database"
)

func IsValidEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)

	return re.MatchString(email)
}

func GetEmailSendApi() string {
	mailMutex.RLock()
	if time.Since(emailServiceUrlLastRefresh) < cacheTTL && len(emailServiceUrl) > 0 {
		result := emailServiceUrl
		mailMutex.RUnlock()
		return result
	}
	mailMutex.RUnlock()

	d := &[]models.Configuration{}
	err := models.FindConfiguration(d, "configuration_group = ?", "EMAIL_SERVICE_URL").Error
	if err != nil {
		return ""
	}

	config := models.ParseConfiguration(d)

	var result string
	if v, ok := config.(map[string]interface{}); ok {
		for _, value := range v {
			if s, ok := value.(string); ok {
				result = s
				break
			}
		}
	}

	mailMutex.Lock()
	emailServiceUrl = result
	emailServiceUrlLastRefresh = time.Now()
	mailMutex.Unlock()

	return result
}

func GetEmailResendConfig() []types.ResendKey {
	mailMutex.RLock()
	if time.Since(emailResendConfigLastRefresh) < cacheTTL && len(emailResendConfig) > 0 {
		result := make([]types.ResendKey, len(emailResendConfig))
		copy(result, emailResendConfig)
		mailMutex.RUnlock()
		return result
	}
	mailMutex.RUnlock()

	d := &[]models.Configuration{}
	err := models.FindConfiguration(d, "configuration_group = ?", "EMAIL_RESEND_CONFIG").Error
	if err != nil {
		return nil
	}

	config := models.ParseConfiguration(d)

	resendKeys := []types.ResendKey{}
	if configList, ok := config.([]map[string]string); ok {
		for _, item := range configList {
			resendKeys = append(resendKeys, types.ResendKey{
				Sender:       item["sender"],
				Email:        item["email"],
				Key:          item["key"],
				ImageElement: StringPtr(item["image_element"]),
			})
		}
	}

	mailMutex.Lock()
	emailResendConfig = resendKeys
	emailResendConfigLastRefresh = time.Now()

	result := make([]types.ResendKey, len(emailResendConfig))
	copy(result, emailResendConfig)
	mailMutex.Unlock()

	return result
}

func GetEmailSMTPConfig() []types.SMTPConfig {
	mailMutex.RLock()
	if time.Since(emailSMTPConfigLastRefresh) < cacheTTL && len(emailSMTPConfig) > 0 {
		result := make([]types.SMTPConfig, len(emailSMTPConfig))
		copy(result, emailSMTPConfig)
		mailMutex.RUnlock()
		return result
	}
	mailMutex.RUnlock()

	d := &[]models.Configuration{}
	err := models.FindConfiguration(d, "configuration_group = ?", "EMAIL_SMTP_CONFIG").Error
	if err != nil {
		return nil
	}

	config := models.ParseConfiguration(d)

	smtpConfig := []types.SMTPConfig{}
	if configList, ok := config.([]map[string]string); ok {
		for _, item := range configList {
			portInt, _ := strconv.Atoi(item["port"])
			smtpConfig = append(smtpConfig, types.SMTPConfig{
				Sender:       item["sender"],
				SMTPUser:     item["email"],
				SMTPHost:     item["host"],
				SMTPPort:     portInt,
				SMTPPassword: item["password"],
				ImageElement: StringPtr(item["image_element"]),
			})
		}
	}

	mailMutex.Lock()
	emailSMTPConfig = smtpConfig
	emailSMTPConfigLastRefresh = time.Now()

	result := make([]types.SMTPConfig, len(emailSMTPConfig))
	copy(result, emailSMTPConfig)
	mailMutex.Unlock()

	return result
}

type Renderer struct {
	Templates map[string]*template.Template
}

var MailRenderer = NewRenderer()

func NewRenderer() *Renderer {
	r := &Renderer{
		Templates: map[string]*template.Template{},
	}

	r.load("testing-mail")
	r.load("user_created")

	return r
}

func (r *Renderer) load(name string) {
	t := template.Must(template.ParseFiles(
		"mail-templates/layout.html",
		"mail-templates/"+name+".html",
	))

	fmt.Println("Loaded templates:")

	for k := range r.Templates {
		fmt.Println("-", k)
	}

	r.Templates[name] = t
}

func (r *Renderer) Render(name string, data any) (string, error) {
	tmpl, ok := r.Templates[name]
	if !ok {
		return "", fmt.Errorf("template %s not found", name)
	}

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	if err := tmpl.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
