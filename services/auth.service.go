package services

import (
	"errors"
	"fmt"
	"strconv"
	"time"

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
		return handlers.ParseBodyErrorResponse(c, err)
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
			return handlers.SuccessResponse(c, errValidate.Error(), user_data, nil)
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

			return handlers.SuccessResponse(c, "Success Login for domain:"+domain, user_data, nil)
		}
	}

	res := fiber.Map{
		"userData": user_data,
		"token":    token,
	}

	return handlers.SuccessResponse(c, "Success Login", res, nil)
}

func AuthByJWT(c *fiber.Ctx) error {
	AuthKey := c.Params("token")
	url := c.Query("url")

	if url == "" {
		return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("need url params"))
	}

	user := new(models.User)

	err := models.FindUserByAuthKey(user, AuthKey).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return handlers.NotFoundErrorResponse(c, fmt.Errorf("auth key is not found"))
	}

	currSession, err := middleware.SessionStore.Get(c)
	if err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}
	defer currSession.Save()

	err = currSession.Regenerate()
	if err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	userSession := currSession.Get("User")

	if userSession == nil {
		currSession.Set(middleware.UserId, user.IdAccount)
		currSession.Set(middleware.AuthKey, user.AuthKey)
		currSession.Set(middleware.Username, user.Username)

		err := currSession.Save()
		if err != nil {
			return handlers.InternalServerErrorResponse(c, err)
		}
	}

	redirectUrl := c.Query("redirectUrl")
	if redirectUrl != "" {
		return c.Redirect(redirectUrl)
	}

	return handlers.SuccessResponse(c, "current session", userSession, nil)
}

func GetSessionJWT(c *fiber.Ctx) error {
	data := new(models.UserData)
	data = &models.UserData{
		IdAccount:       c.Locals("id_account").(string),
		IdentityNumber:  c.Locals("identity_number").(string),
		Username:        c.Locals("username").(string),
		FullName:        c.Locals("full_name").(string),
		Email:           c.Locals("email").(string),
		PhoneNumber:     c.Locals("phone_number").(string),
		DateOfBirth:     c.Locals("date_of_birth").(*time.Time),
		AuthKey:         c.Locals("auth_key").(string),
		UsedPin:         c.Locals("used_pin").(int8),
		IsGoogleAccount: c.Locals("is_google_account").(int8),
		StatusAccount:   c.Locals("status_account").(int8),
		LoginIp:         c.Locals("login_ip").(string),
		LoginAttempts:   c.Locals("login_attempts").(int8),
		LoginTime:       c.Locals("login_time").(int64),
		CreatedAt:       c.Locals("created_at").(int),
		UpdatedAt:       c.Locals("updated_at").(int),
	}

	// return c.JSON(fiber.Map{"success": true, "message": data})
	return handlers.SuccessResponse(c, "Success Get JWT", data, nil)
}

func GetSession(c *fiber.Ctx) error {
	currSession, err := middleware.SessionStore.Get(c)
	if err != nil {
		return handlers.UnauthorizedErrorResponse(c, err)
	}

	if len(currSession.Keys()) > 0 {

		userId, ok := currSession.Get(middleware.UserId).(string)
		if !ok {
			return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("session userId is null"))
		}
		authKey, _ := currSession.Get(middleware.AuthKey).(string)
		username, _ := currSession.Get(middleware.Username).(string)
		fullName, _ := currSession.Get(middleware.FullName).(string)

		data := new(models.UserData)
		data = &models.UserData{
			IdAccount: userId,
			AuthKey:   authKey,
			Username:  username,
			FullName:  fullName,
		}

		return handlers.SuccessResponse(c, "success get session", data, nil)
	}

	return handlers.SuccessResponse(c, "no session", nil, nil)
}

func LogoutWeb(c *fiber.Ctx) error {
	currSession, err := middleware.SessionStore.Get(c)
	if err != nil {
		return handlers.UnauthorizedErrorResponse(c, err)
	}

	if err := currSession.Destroy(); err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	currSession.SetExpiry(time.Second * -60)

	cookieJWT := new(fiber.Cookie)
	cookieJWT.Name = middleware.CookieJWT
	// cookieJWT.Expires = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	cookieJWT.Expires = time.Now().Add(-(time.Hour * 2))
	c.Cookie(cookieJWT)

	redirectUrl := c.Query("redirectUrl")
	if redirectUrl != "" {
		return c.Redirect(redirectUrl)
	}

	return handlers.SuccessResponse(c, "success logout", currSession, nil)
}
