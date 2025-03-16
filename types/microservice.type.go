package types

type Service struct {
	Path              string `json:"path"`
	Url               string `json:"url"`
	AuthProtection    bool   `json:"auth_protection"`
	CsrfProtection    bool   `json:"csrf_protection"`
	SessionProtection bool   `json:"session_protection"`
	JwtProtection     bool   `json:"jwt_protection"`
	Status            bool   `json:"status"`
}

type ConfigServices struct {
	Services []Service `json:"services"`
}

type LoginInput struct {
	Id       string `json:"id"`
	Identity string `json:"identity" validate:"required"`
	Password string `json:"password" validate:"required"`
	Captcha  int    `json:"captcha"`
}

type GoogleLogin struct {
	IdToken string `json:"id_token"`
}

type ResetPasswordRequest struct {
	Identity string `json:"identity" validate:"required"`
}
type ResetPasswordInput struct {
	Id              string `json:"id"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"passwordConfirm" validate:"required"`
}

type Email struct {
	Name      string `json:"name"`
	EmailAddr string `json:"email_addr"`
}

type ListEmail struct {
	Sender           string  `json:"sender"`
	Subject          string  `json:"subject"`
	BodyTemplateText string  `json:"body_template_text"`
	BodyTemplateHtml string  `json:"body_template_html"`
	Emails           []Email `json:"emails"`
}

type SendingEmailToBroker struct {
	Sender   string  `json:"sender"`
	Subject  string  `json:"subject"`
	Title    string  `json:"title"`
	BodyText string  `json:"bodyText"`
	Body     string  `json:"body"`
	Footer   string  `json:"footer"`
	Emails   []Email `json:"emails"`
}
