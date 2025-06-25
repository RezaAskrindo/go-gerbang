package services

import (
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-gerbang/config"
	"go-gerbang/database"
	"go-gerbang/handlers"
	"go-gerbang/middleware"
	"go-gerbang/models"
	"go-gerbang/types"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Signup(c *fiber.Ctx) error {
	user := new(models.User)

	if err := handlers.ParseBody(c, user); err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	if err := handlers.ValidateStruct(*user); err != nil {
		return c.Status(fiber.StatusOK).JSON(err)
	}

	userExist := new(models.User)
	if err := models.FindUserByIdentity(userExist, user.Username, user.Email, user.PhoneNumber, user.IdentityNumber); err == nil {
		return handlers.SuccessResponse(c, false, "Account Already Exist", fiber.Map{
			"username": fiber.Map{
				"invalid": true,
				"desc":    "This Account Already Exist",
				"descRaw": "This Account Already Exist",
			},
		}, nil)
	} // Reserve Algoritm

	user.PasswordHash = handlers.GeneratePasswordHash(user.Password)

	isActive := c.QueryBool("active")
	if isActive {
		user.StatusAccount = 10
	}

	if err := models.CreateUser(user); err.Error != nil {
		return handlers.ConflictErrorResponse(c, err.Error)
	}

	sendNotification := c.QueryBool("notif")
	QuerySender := c.Query("sender")
	sendPass := c.QueryBool("sendPass")

	if sendNotification {
		Sender := "GOGERBANG"
		if QuerySender != "" {
			Sender = QuerySender
		}

		textPass := ``
		htmlPass := ``
		if sendPass {
			textPass = `password: ` + user.Password
			htmlPass = `<div>password: <strong>` + user.Password + `</strong></div>`
		}

		sendEmail := new(types.SendingEmailToBroker)
		sendEmail.Sender = Sender
		sendEmail.Subject = "Create Account Success"
		sendEmail.Title = "Akun Anda Berhasil Di Buat"
		sendEmail.BodyText = `
			Hi, ` + user.FullName + `, berikut inform
			for reset password pleasasi akun anda:e click link below

			username: ` + user.Username + `
			email: ` + user.Email + textPass + `

			Tetap jaga rahasia akun anda, mohon untuk jangan diberikan kepada siapapun termasuk Admin.
		`
		sendEmail.Body = `
			<div class="font-family: Roboto-Regular, Helvetica, Arial, sans-serif; font-size: 14px; color: rgba(0, 0, 0, 0.87); padding-top: 20px; text-align: center;">Hi, ` + user.FullName + `, berikut informasi akun anda:</div>
			<br/>
			<div>username: <strong>` + user.FullName + `</strong></div>
			<div>email: <strong>` + user.Email + `</strong></div>` + htmlPass + `
			<br/>
			<div class="padding-top: 20px; font-size: 12px; line-height: 16px; color: rgb(95, 99, 104); letter-spacing: 0.3px; text-align: center;">Tetap jaga rahasia akun anda, mohon untuk jangan diberikan kepada siapapun termasuk Admin.</div>
		`
		sendEmail.Footer = "ini merupakan email otomatis dari " + Sender
		sendEmail.Emails = []types.Email{
			{
				Name:      user.FullName,
				EmailAddr: user.Email,
			},
		}

		// go PublishServiceEmail(sendEmail)
		PublishEvent("user.notification", sendEmail)
	}

	return handlers.SuccessResponse(c, true, "Success Create User", nil, nil)
}

// @Summary Login
// @Description Login Api Explaination
// @Tags auth
// @Accept json
// @Produce json
// @Param user body types.LoginInput true "Login Input"
// @Success 200 {object} ResponseHTTP{data=SuccessStruct}
// @Failure 400 {object} ResponseHTTP{}
// @Router /api/v1/auth/login [post]
func Login(c *fiber.Ctx) error {
	captcha := c.QueryBool("captcha")
	block := c.QueryBool("block")
	session := c.QueryBool("session")
	httponly := c.QueryBool("httponly")
	domain := c.Query("domain")
	validate_ip := c.QueryBool("validate_ip")
	single_login := c.QueryBool("single_login")

	input := new(types.LoginInput)

	if err := handlers.ParseBody(c, input); err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	if err := handlers.ValidateStruct(*input); err != nil {
		return c.Status(fiber.StatusOK).JSON(err)
	}

	if captcha { // CAPTCHA QUERY
		sess, err := middleware.CaptchaStore.Get(c)
		if err != nil {
			return handlers.InternalServerErrorResponse(c, err)
		}

		sessionCaptcha := sess.Get("captcha")
		if sessionCaptcha == nil {
			return handlers.InternalServerErrorResponse(c, fmt.Errorf("session captcha is null"))
		}

		captcha := sessionCaptcha.(string)
		intCaptcha, _ := strconv.Atoi(captcha)

		if input.Captcha != intCaptcha {
			return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("invalid captcha"))
		}
	}

	identity := input.Identity

	user := new(models.User)
	if err := models.FindUserByIdentity(user, identity, identity, identity, identity); err != nil {
		return handlers.NotFoundErrorResponse(c, err)
	}

	if user.StatusAccount == 0 {
		return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("your account is not active or blocked"))
	}

	password := input.Password
	u := new(models.User)

	if !handlers.CheckPasswordHash(password, user.PasswordHash) {
		u.LoginAttempts = user.LoginAttempts + 1
		u.LoginIp = c.IP()
		if block { // BLOCK Query
			if user.LoginAttempts >= 3 {
				models.BlockUser(user.IdAccount)
				return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("you're account has block, you're already fill wrong password 3 time"))
			} else {
				models.UpdateUser(user.IdAccount, u)
				return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("%s", "wrong password, you have "+strconv.Itoa(4-(int(u.LoginAttempts)))+" chances left"))
			}
		} else {
			if user.PasswordHash == "" && user.IsGoogleAccount == 10 {
				return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("try login with Google"))
			}
			models.UpdateUser(user.IdAccount, u)
			return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("wrong password"))
		}
	} else {
		u.LoginAttempts = 1
		u.LoginIp = c.IP()
		u.LoginTime = handlers.TimeNow.Unix()
		models.UpdateUser(user.IdAccount, u)
	}

	randString := handlers.RandomString(32)

	user_data := handlers.SendSafeUserData(user, randString)

	if err := models.GenerateAuthKeyUser(user_data.IdAccount, user_data.AuthKey).Error; err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	if session { // SESSION QUERY
		err := middleware.SaveUserSession(c, user_data, single_login) // SINGLE LOGIN
		if err != nil {
			return handlers.InternalServerErrorResponse(c, err)
		}
	}

	if validate_ip { // VALIDATE IP QUERY
		errValidate := handlers.ValidateUserLoginIp(user_data, c)
		if errValidate != nil {
			return handlers.SuccessResponse(c, true, errValidate.Error(), user_data, nil)
		}
	}

	// refreshToken, err := handlers.GenerateRefreshToken(user_data)
	refreshToken, err := handlers.GenerateTokenJWT(user_data, true)
	if err != nil {
		return handlers.InternalServerErrorResponse(c, fmt.Errorf("failed to generate new refresh token"))
	}

	token, err := handlers.GenerateTokenJWT(user_data, false)
	if err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	if httponly { // HTTPONLY QUERY
		if domain == "" {
			return handlers.UnprocessableEntityErrorResponse(c, fmt.Errorf("need domain params"))
		} else {
			c.Cookie(&fiber.Cookie{
				Name:     middleware.CookieRefreshJWT,
				Value:    refreshToken,
				HTTPOnly: true,
				Secure:   true,
				SameSite: "Strict",
				Expires:  time.Now().Add(config.RefreshAuthTimeCache),
			})

			cookie := new(fiber.Cookie)
			cookie.Name = middleware.CookieJWT
			cookie.Value = "Bearer " + token
			cookie.Expires = time.Now().Add(config.AuthTimeCache)
			cookie.HTTPOnly = true
			cookie.Domain = domain
			cookie.Secure = config.SecureCookies
			cookie.SameSite = config.CookieSameSite
			cookie.SessionOnly = false
			c.Cookie(cookie)

			return handlers.SuccessResponse(c, true, "Success Login for domain:"+domain, user_data, nil)
		}
	}

	res := fiber.Map{
		"userData":     user_data,
		"token":        token,
		"refreshToken": refreshToken,
	}

	return handlers.SuccessResponse(c, true, "Success Login", res, nil)
}

func ValidateUserPasswordById(c *fiber.Ctx) error {
	b := new(types.LoginInput)

	if err := handlers.ParseBody(c, b); err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	user := new(models.User)
	if err := models.FindUserById(user, b.Id); err != nil {
		return handlers.NotFoundErrorResponse(c, err)
	}

	if user.PasswordHash != "" {
		if !handlers.CheckPasswordHash(b.Password, user.PasswordHash) {
			return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("your password is wrong"))
		}
	}

	return handlers.SuccessResponse(c, true, "Password Is Valid", nil, nil)
}

func ChangePassword(c *fiber.Ctx) error {
	b := new(types.ResetPasswordInput)

	if err := handlers.ParseBody(c, b); err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	user := new(models.User)
	if err := models.FindUserById(user, b.Id); err != nil {
		return handlers.NotFoundErrorResponse(c, err)
	}

	user.PasswordHash = handlers.GeneratePasswordHash(b.Password)
	user.PasswordResetToken = nil
	if err := models.UpdateUserPassword(user.IdAccount, user).Error; err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	return handlers.SuccessResponse(c, true, "Change Password Is Success", nil, nil)
}

func RequestResetPassword(c *fiber.Ctx) error {
	input := new(types.LoginInput)

	if err := handlers.ParseBody(c, input); err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	accountEmail := input.Identity

	user := new(models.User)
	if err := models.FindUserByIdentity(user, accountEmail, accountEmail, accountEmail, accountEmail); err != nil {
		return handlers.NotFoundErrorResponse(c, err)
	}

	randomReset := handlers.GenerateResetRandom(64)

	if err := models.CeneratePasswordResetToken(user.IdAccount, randomReset).Error; err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	QuerySender := c.Query("sender")

	Sender := "GOGERBANG"
	if QuerySender != "" {
		Sender = QuerySender
	}

	BaseUrl := c.Query("url")
	if BaseUrl == "" {
		return handlers.UnprocessableEntityErrorResponse(c, fmt.Errorf("need base url params"))
	}

	sendEmail := new(types.SendingEmailToBroker)
	sendEmail.Sender = Sender
	sendEmail.Subject = "Reset Password"
	sendEmail.Title = "You are request for reset password"
	sendEmail.BodyText = `Hi ` + user.FullName + `
	
	for reset password please click link below:
		
	` + BaseUrl + `/forget-password?token=` + randomReset + `

	this link only active in 24 hours`
	sendEmail.Body = `Hi ` + user.FullName + `<br/> 
		for reset password please click link below:
		<div style="padding-top:30px;padding-bottom:28px;text-align:center">
			<a href='` + BaseUrl + `/forget-password?token=` + randomReset + `' style="font-family:'Google Sans',Roboto,RobotoDraft,Helvetica,Arial,sans-serif;line-height:16px;color:#ffffff;font-weight:400;text-decoration:none;font-size:14px;display:inline-block;padding:10px 24px;background-color:#171717;border-radius:5px;min-width:90px">Reset Link</a>
		</div>
		this link only active in 24 hours`
	sendEmail.Footer = "you are receiving this mail from " + Sender
	sendEmail.Emails = []types.Email{
		{
			Name:      user.FullName,
			EmailAddr: accountEmail,
		},
	}

	// go PublishServiceEmail(sendEmail)
	PublishEvent("user.notification", sendEmail)

	return handlers.SuccessResponse(c, true, "Silahkan Cek Email", nil, nil)
}

func ResetPassword(c *fiber.Ctx) error {
	Token := c.Query("token")

	if Token == "" {
		return handlers.UnprocessableEntityErrorResponse(c, fmt.Errorf("need token params"))
	}

	input := new(types.ResetPasswordInput)

	if err := handlers.ParseBody(c, input); err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	if err := handlers.ValidateStruct(*input); err != nil {
		return c.Status(fiber.StatusOK).JSON(err)
	}

	apiKey := "Go_Gerbang_Is_Key"

	if config.Config("RESET_PASSWORD_DEFAULT_KEY") != "" {
		apiKey = config.Config("RESET_PASSWORD_DEFAULT_KEY")
	}

	user := new(models.User)

	hashedAPIKey := sha256.Sum256([]byte(apiKey))
	hashedKey := sha256.Sum256([]byte(Token))

	if subtle.ConstantTimeCompare(hashedAPIKey[:], hashedKey[:]) == 1 {
		if err := models.FindUserById(user, input.Id); err != nil {
			return handlers.NotFoundErrorResponse(c, err)
		}
	} else { // if not found token
		err := models.FindUserByPasswordReset(user, Token).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return handlers.NotFoundErrorResponse(c, fmt.Errorf("your token is not exist"))
		}

		if !handlers.IsPasswordResetTokenValid(Token) {
			return handlers.InternalServerErrorResponse(c, fmt.Errorf("your token is not valid anymore"))
		}
	}

	user.PasswordHash = handlers.GeneratePasswordHash(input.Password)
	user.PasswordResetToken = nil
	if err := models.UpdateUserPassword(user.IdAccount, user).Error; err != nil {
		return handlers.NotFoundErrorResponse(c, err)
	}

	return handlers.SuccessResponse(c, true, "Reset Password Berhasil, Silahkan Login Kembali", nil, nil)
}

func RefreshAuth(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		token = c.Cookies(middleware.CookieRefreshJWT)
	}

	if token == "" {
		return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("missing refresh token"))
	}

	var refreshToken string
	if strings.HasPrefix(token, "Bearer ") {
		refreshToken = token[len("Bearer "):]
	} else {
		refreshToken = token
	}

	user, err := middleware.Verify(refreshToken, "refresh")
	if err != nil {
		return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("invalid refresh token"))
	}

	jti := user.Jti
	if jti == nil || *jti == "" || handlers.IsTokenBlacklisted(*jti) {
		return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("this refresh token is no longer valid"))
	}

	database.RedisDb.Del(database.RedisCtx, "refresh:"+user.IdAccount)

	newRefreshToken, err := handlers.GenerateTokenJWT(*user, true)
	if err != nil {
		return handlers.InternalServerErrorResponse(c, fmt.Errorf("failed to generate new refresh token"))
	}

	accessToken, err := handlers.GenerateTokenJWT(*user, false)
	if err != nil {
		return handlers.InternalServerErrorResponse(c, fmt.Errorf("failed to generate access token"))
	}

	httponly := c.QueryBool("httponly")
	domain := c.Query("domain")

	if httponly { // HTTPONLY QUERY
		if domain == "" {
			return handlers.UnprocessableEntityErrorResponse(c, fmt.Errorf("need domain params"))
		} else {
			c.Cookie(&fiber.Cookie{
				Name:     middleware.CookieRefreshJWT,
				Value:    newRefreshToken,
				HTTPOnly: true,
				Secure:   true,
				SameSite: "Strict",
				Expires:  time.Now().Add(config.RefreshAuthTimeCache),
			})

			cookie := new(fiber.Cookie)
			cookie.Name = middleware.CookieJWT
			cookie.Value = "Bearer " + accessToken
			cookie.Expires = time.Now().Add(config.AuthTimeCache)
			cookie.HTTPOnly = true
			cookie.Domain = domain
			cookie.Secure = config.SecureCookies
			cookie.SameSite = config.CookieSameSite
			cookie.SessionOnly = false
			c.Cookie(cookie)

			return handlers.SuccessResponse(c, true, "Success Login for domain:"+domain, user, nil)
		}
	}

	res := fiber.Map{
		"userData":     user,
		"token":        accessToken,
		"refreshToken": newRefreshToken,
	}

	return handlers.SuccessResponse(c, true, "Refresh token rotated", res, nil)
}
