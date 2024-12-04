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
	Id       int    `json:"id"`
	Identity string `json:"identity" example:"Muhammad Reza"`
	Password string `json:"password" example:"12345"`
	Captcha  int    `json:"captcha"`
}

type GoogleLogin struct {
	IdToken string `json:"id_token"`
}
