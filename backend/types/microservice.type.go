package types

type Service struct {
	Service           string `json:"service"`
	Path              string `json:"path"`
	Url               string `json:"url"`
	AuthProtection    bool   `json:"auth_protection"`
	SessionProtection bool   `json:"session_protection"`
	CsrfProtection    bool   `json:"csrf_protection"`
	RbacProtection    bool   `json:"rbac_protection"`
	Status            bool   `json:"status"`
	// JwtProtection     bool   `json:"jwt_protection"`
}

type ConfigServices struct {
	Services []Service `json:"services"`
}

type LoginInput struct {
	Id       string `json:"id"`
	Identity string `json:"identity" validate:"required"`
	Password string `json:"password" validate:"required"`
	Captcha  int    `json:"captcha"`
	OTP      bool   `json:"otp"`
}

// type GoogleLogin struct {
// 	IdToken  string `json:"id_token"`
// 	ClientId string `json:"client_id"`
// }

// type ResetPasswordRequest struct {
// 	Identity string `json:"identity" validate:"required"`
// }

type ResetPasswordInput struct {
	Id              string `json:"id"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"passwordConfirm" validate:"required"`
}

type Email struct {
	Name      string `json:"name"`
	EmailAddr string `json:"email_addr"`
}

type Attachment struct {
	Name        string `json:"name"`
	ContentType string `json:"content_type,omitempty"`
	Data        []byte `json:"-"`
}

type ListEmail struct {
	Sender           string       `json:"sender"`
	Subject          string       `json:"subject"`
	BodyTemplateText string       `json:"body_template_text"`
	BodyTemplateHtml string       `json:"body_template_html"`
	Emails           []Email      `json:"emails"`
	Attachments      []Attachment `json:"attachments,omitempty"`
	TypeBatchAddress string       `json:"type_batch_address"` // "all" or "single" default "single"
}

type SendingEmailToBroker struct {
	Sender   string `json:"sender"`
	Provider string `json:"provider"`
	Subject  string `json:"subject"`
	// Title    string  `json:"title"`
	// BodyText string  `json:"bodyText"`
	// Body     string  `json:"body"`
	// Footer   string  `json:"footer"`
	Template string
	Data     map[string]any
	// Data   json.RawMessage
	Emails []Email `json:"emails"`
}

type LoginQuery struct {
	Domain      string  `query:"domain"`
	Captcha     bool    `query:"captcha"`
	Block       bool    `query:"block"`
	Session     bool    `query:"session"`
	HttpOnly    bool    `query:"httponly"`
	ValidateIp  bool    `query:"validate_ip"`
	SingleLogin bool    `query:"single_login"`
	OTP         bool    `query:"otp"`
	AccountID   *string `query:"account_id"`
	Destination *string `query:"to_destination"`
	// FOR GOOGLE
	ClientId  string `query:"client_id"`
	CreateNew bool   `query:"create_new"`
}

type VerifyOTPRequest struct {
	UserId  string `json:"user_id" validate:"required"`
	OTPCode string `json:"otp_code" validate:"required"`
	OTPType string `json:"otp_type"` // "totp" or "simple"
}

type ApiInfo struct {
	Name    string
	Version string
	Maker   string
}
