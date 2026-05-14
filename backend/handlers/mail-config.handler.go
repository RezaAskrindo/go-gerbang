package handlers

import (
	"regexp"
	"strconv"
	"sync"
	"time"

	"go-gerbang/models"
	"go-gerbang/types"
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

// func GetEmailServiceName() map[string]interface{} {
// 	mailMutex.RLock()
// 	if time.Since(emailServiceNameLastRefresh) < cacheTTL && len(emailServiceName) > 0 {
// 		result := make(map[string]interface{}, len(emailServiceName))
// 		for k, v := range emailServiceName {
// 			result[k] = v
// 		}
// 		mailMutex.RUnlock()
// 		return result
// 	}
// 	mailMutex.RUnlock()

// 	d := &[]models.Configuration{}
// 	err := models.FindConfiguration(d, "configuration_group = ?", "EMAIL_SERVICE_NAME").Error
// 	if err != nil {
// 		return nil
// 	}

// 	config := models.ParseConfiguration(d)

// 	mailMutex.Lock()
// 	if configMap, ok := config.(map[string]struct{}); ok {
// 		emailServiceName = configMap
// 	}
// 	emailServiceNameLastRefresh = time.Now()
// 	result := make(map[string]interface{}, len(emailServiceName))
// 	for k, v := range emailServiceName {
// 		result[k] = v
// 	}
// 	mailMutex.Unlock()

// 	return result
// }

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
