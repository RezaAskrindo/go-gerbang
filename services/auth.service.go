package services

import (
	"errors"
	"fmt"
	"strconv"

	"go-gerbang/config"
	"go-gerbang/handlers"
	"go-gerbang/middleware"
	"go-gerbang/models"
	"go-gerbang/types"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

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

	// CAPTCHA QUERY
	if captcha {
		sess, err := middleware.CaptchaStore.Get(c)
		if err != nil {
			return handlers.InternalServerErrorResponse(c, err)
		}

		sessionCaptcha := sess.Get("captcha")
		if sessionCaptcha == nil {
			return handlers.InternalServerErrorResponse(c, fmt.Errorf("session captcha is null"))
		}

		captcha := sessionCaptcha.(string)
		intCaptcha, err := strconv.Atoi(captcha)

		if input.Captcha != intCaptcha {
			return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("invalid captcha"))
		}
	}

	identity := input.Identity

	user, err := models.FindUserByIdentity(identity)
	if err != nil {
		return handlers.BadRequestErrorResponse(c, fmt.Errorf("error on find user"))
	}

	if user.StatusAccount == 0 {
		return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("your account is not active or blocked"))
	}

	password := input.Password
	u := new(models.User)

	if user.PasswordHash != "" {
		// passwordHash := *user.PasswordHash
		if !handlers.CheckPasswordHash(password, user.PasswordHash) {
			u.LoginAttempts = user.LoginAttempts + 1
			u.LoginIp = c.IP()
			// BLOCK Query
			if block {
				if user.LoginAttempts >= 3 {
					models.BlockUser(user.IdAccount)
					return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("you're account has block, you're already fill wrong password 3 time"))
				} else {
					models.UpdateUser(user.IdAccount, u)
					return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("wrong password, you have "+strconv.Itoa(4-(int(u.LoginAttempts)))+" chances left"))
				}
			} else {
				models.UpdateUser(user.IdAccount, u)
				return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("wrong password"))
			}
		} else {
			u.LoginAttempts = 1
			u.LoginIp = c.IP()
			u.LoginTime = handlers.TimeNow.Unix()
			models.UpdateUser(user.IdAccount, u)
		}
	}

	randString := handlers.RandomString(32)

	user_data := handlers.SendSafeUserData(user, randString)

	if err := models.GenerateAuthKeyUser(user_data.IdAccount, user_data.AuthKey).Error; err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	// SESSION QUERY
	if session {
		err := middleware.SaveUserSession(c, user_data, single_login) // SINGLE LOGIN
		if err != nil {
			return handlers.InternalServerErrorResponse(c, err)
		}
	}

	// VALIDATE IP QUERY
	if validate_ip {
		errValidate := handlers.ValidateUserLoginIp(user_data, c)
		if errValidate != nil {
			return handlers.SuccessResponse(c, true, errValidate.Error(), user_data, nil)
		}
	}

	token, err := handlers.GenerateTokenJWT(user_data, c)
	if err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	if httponly {
		if domain == "" {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "For httponly need domain")
		} else {
			cookie := new(fiber.Cookie)
			cookie.Name = middleware.CookieJWT
			cookie.Value = "Bearer " + token
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
		"userData": user_data,
		"token":    token,
	}

	return handlers.SuccessResponse(c, true, "Success Login", res, nil)
}

func ValidateUserPasswordById(c *fiber.Ctx) error {
	b := new(types.LoginInput)

	if err := handlers.ParseBody(c, b); err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	user, err := models.FindUserById(b.Id)
	if err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	if user.PasswordHash != "" {
		if !handlers.CheckPasswordHash(b.Password, user.PasswordHash) {
			return handlers.UnauthorizedErrorResponse(c, err)
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

	err := models.FindUserByIdRaw(user, b.Id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
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

	user, err := models.FindUserByIdentity(accountEmail)
	if err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	randomReset := handlers.GenerateResetRandom(64)

	if err := models.CeneratePasswordResetToken(user.IdAccount, randomReset).Error; err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	QuerySender := c.Query("sender")
	QueryUrl := c.Query("url")

	Sender := "Go Gerbang"
	if QuerySender != "" {
		Sender = QuerySender
	}
	BaseUrl := "//localhost"
	if QueryUrl != "" {
		BaseUrl = QueryUrl
	}

	sendEmail := new(types.SendingEmail)
	sendEmail.Sender = Sender
	sendEmail.Subject = "Reset Password"
	sendEmail.Title = "You are request for reset password"
	sendEmail.Body = "Yth " + user.FullName + " <br/> Mohon klik link di bawah ini untuk reset Password <br/><br/> <a href='" + BaseUrl + "/forget-password?token=" + randomReset + "'>reset link</a> <br/><br/> Link Ini aktif hanya dalam 24 Jam."
	sendEmail.Footer = "ini merupakan email otomatis dari " + Sender
	sendEmail.Emails = []types.Email{
		{
			Name:      user.FullName,
			EmailAddr: accountEmail,
		},
	}

	go PublishServiceEmail(sendEmail)

	return handlers.SuccessResponse(c, true, "Silahkan Cek Email", nil, nil)
}

func ResetPassword(c *fiber.Ctx) error {
	Token := c.Query("token")

	if Token == "" {
		return handlers.UnprocessableEntityErrorResponse(c, fmt.Errorf("need token params"))
	}

	user := new(models.User)

	err := models.FindUserByPasswordReset(user, Token).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return handlers.NotFoundErrorResponse(c, fmt.Errorf("your token is not exist"))
	}

	if !handlers.IsPasswordResetTokenValid(Token) {
		return handlers.InternalServerErrorResponse(c, fmt.Errorf("your token is not valid anymore"))
	}

	input := new(types.ResetPasswordInput)

	if err := handlers.ParseBody(c, input); err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	if err := handlers.ValidateStruct(*input); err != nil {
		return c.Status(fiber.StatusOK).JSON(err)
	}

	user.PasswordHash = handlers.GeneratePasswordHash(input.Password)
	user.PasswordResetToken = nil
	if err := models.UpdateUserPassword(user.IdAccount, user).Error; err != nil {
		return handlers.NotFoundErrorResponse(c, err)
	}

	return handlers.SuccessResponse(c, true, "Reset Password Berhasil, Silahkan Login Kembali", nil, nil)
}
