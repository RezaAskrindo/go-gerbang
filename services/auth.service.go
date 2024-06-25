package services

import (
	"errors"
	"sika_apigateway/config"
	"sika_apigateway/handlers"
	"sika_apigateway/middleware"
	"sika_apigateway/models"
	"sika_apigateway/types"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Login(c *fiber.Ctx) error {
	captcha := c.QueryBool("captcha")
	block := c.QueryBool("block")
	session := c.QueryBool("session")
	httponly := c.QueryBool("httponly")
	domain := c.Query("domain")
	validate_ip := c.QueryBool("validate_ip")

	input := new(types.LoginInput)

	if err := handlers.ParseBody(c, input); err != nil {
		return handlers.ParseBodyErrorResponse(c, err)
	}

	// CAPTCHA QUERY
	if captcha {
		sess, err := middleware.CaptchaStore.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
		}

		sessionCaptcha := sess.Get("captcha")
		if sessionCaptcha == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Session Captcha is null"})
		}

		captcha := sessionCaptcha.(string)
		intCaptcha, err := strconv.Atoi(captcha)

		if input.Captcha != intCaptcha {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Invalid Captcha"})
		}
	}

	identity := input.Identity

	user, err := models.FindUserByIdentity(identity)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Error On Find User"})
	}

	if user.StatusAccount == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Your Account is not active or blocked"})
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
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "You're Account Has Block, You're Already Fill Wrong Password 3 Time."})
				} else {
					models.UpdateUser(user.IdAccount, u)
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Wrong password, you have " + strconv.Itoa(4-(int(u.LoginAttempts))) + " chances left"})
				}
			} else {
				models.UpdateUser(user.IdAccount, u)
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Wrong password"})
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	// SESSION QUERY
	if session {
		err := middleware.SaveUserSession(user_data, c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
	}

	// VALIDATE IP QUERY
	if validate_ip {
		errValidate := handlers.ValidateUserLoginIp(user_data, c)
		if errValidate != nil {
			return c.JSON(fiber.Map{"success": false, "message": errValidate.Error(), "data": user_data})
		}
	}

	token, err := handlers.GenerateTokenJWT(user_data, c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	if httponly {
		if domain == "" {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "For httponly need domain")
		} else {
			cookie := new(fiber.Cookie)
			cookie.Name = "__SecureGatewayJ"
			cookie.Value = "Bearer " + token
			cookie.HTTPOnly = true
			cookie.Domain = domain
			cookie.Secure = config.SecureCookies
			c.Cookie(cookie)

			return c.JSON(fiber.Map{"success": true, "message": "Success Login for domain:" + domain, "data": user_data})
		}
	}

	return c.JSON(fiber.Map{"success": true, "message": "Success Login", "token": token, "data": user_data})
}

func LoginWithGoogle(c *fiber.Ctx) error {
	session := c.QueryBool("session")
	validate_ip := c.QueryBool("validate_ip")

	b := new(types.GoogleLogin)

	if err := handlers.ParseBody(c, b); err != nil {
		return handlers.ParseBodyErrorResponse(c, err)
	}

	tokenInfo, err := handlers.VerifyIdTokenGoogle(b.IdToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	user, err := models.FindUserByIdentity(tokenInfo.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	if user.StatusAccount == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Your Account is not active or blocked"})
	}

	randString := handlers.RandomString(32)

	user_data := handlers.SendSafeUserData(user, randString)

	if err := models.GenerateAuthKeyUser(user_data.IdAccount, user_data.AuthKey).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	// SESSION QUERY
	if session {
		err := middleware.SaveUserSession(user_data, c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
	}

	// VALIDATE IP QUERY
	if validate_ip {
		errValidate := handlers.ValidateUserLoginIp(user_data, c)
		if errValidate != nil {
			return c.JSON(fiber.Map{"success": false, "message": errValidate.Error(), "data": user_data})
		}
	}

	token, err := handlers.GenerateTokenJWT(user_data, c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Success Login", "token": token, "data": user_data})
}

func AuthByJWT(c *fiber.Ctx) error {
	AuthKey := c.Params("token")
	url := c.Query("url")

	if AuthKey == "" {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Need AuthKey")
	}

	if url == "" {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Need url")
	}

	user := new(models.User)

	err := models.FindUserByAuthKey(user, AuthKey).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.SendStatus(fiber.StatusNotFound)
	}

	currSession, err := middleware.SessionStore.Get(c)
	if err != nil {
		return err
	}
	defer currSession.Save()

	err = currSession.Regenerate()
	if err != nil {
		return err
	}

	userSession := currSession.Get("User")

	if userSession == nil {
		currSession.Set("User", user.Username)
		currSession.Set("UserID", user.IdAccount)
		currSession.Set("auth_key", user.AuthKey)
	}

	return c.Redirect(url)
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

	return c.JSON(fiber.Map{"success": true, "message": data})
}

func LogoutWeb(c *fiber.Ctx) error {
	redirectUrl := c.Query("redirectUrl")

	if redirectUrl == "" {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Need redirectUrl")
	}

	currSession, err := middleware.SessionStore.Get(c)
	if err != nil {
		return err
	}
	defer currSession.Save()

	err = currSession.Regenerate()
	if err != nil {
		return err
	}

	userSession := currSession.Get("User")

	if userSession != nil {
		currSession.Delete("User")
		if err := currSession.Destroy(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
		}
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "__SecureGatewayJ"
	cookie.Expires = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	c.Cookie(cookie)

	return c.Redirect(redirectUrl)
}
