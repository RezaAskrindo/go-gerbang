package types

type ValueMicroService struct {
	Path              string `json:"path"`
	Url               string `json:"url"`
	AuthProtection    bool   `json:"auth_protection"`
	CsrfProtection    bool   `json:"csrf_protection"`
	SessionProtection bool   `json:"session_protection"`
	JwtProtection     bool   `json:"jwt_protection"`
	Status            bool   `json:"status"`
}

// type ValueMicroServiceResponses struct {
// 	Menu []ValueMicroService `json:"menu"`
// }

type LoginInput struct {
	Id       int    `json:"id"`
	Identity string `json:"identity"`
	Password string `json:"password"`
	Captcha  int    `json:"captcha"`
}

type GoogleLogin struct {
	IdToken string `json:"id_token"`
}
