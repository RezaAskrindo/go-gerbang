package services

import (
	"context"
	"fmt"

	"go-gerbang/config"
	"go-gerbang/handlers"
	"go-gerbang/middleware"
	"go-gerbang/models"
	"go-gerbang/types"

	"github.com/gofiber/fiber/v2"
)

func LoadGoogleLoginClienId(c *fiber.Ctx) error {
	return handlers.SuccessResponse(c, true, "success getting google client id", config.Config("GOOGLE_LOGIN_CLIENT_ID"), nil)
}

func LoginWithGoogle(c *fiber.Ctx) error {
	session := c.QueryBool("session")
	validate_ip := c.QueryBool("validate_ip")
	httponly := c.QueryBool("httponly")
	domain := c.Query("domain")

	b := new(types.GoogleLogin)

	if err := handlers.ParseBody(c, b); err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	payload, err := handlers.VerifyIdTokenGoogle(context.Background(), b.IdToken, b.ClientId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	email, _ := payload.Claims["email"].(string)
	phone_number, _ := payload.Claims["phone_number"].(string)

	isActive := c.QueryBool("active")

	user := new(models.User)
	if err := models.FindUserByIdentity(user, email, email, phone_number, email); err != nil {
		name, _ := payload.Claims["name"].(string)

		user.FullName = name
		user.Username = email
		user.Email = email
		user.PhoneNumber = phone_number
		user.IsGoogleAccount = 10

		if isActive {
			user.StatusAccount = 10
		}

		if err := models.CreateUser(user); err.Error != nil {
			return handlers.ConflictErrorResponse(c, err.Error)
		}

		sendNotification := c.QueryBool("notif")
		QuerySender := c.Query("sender")

		if sendNotification {
			Sender := "GOGERBANG"
			if QuerySender != "" {
				Sender = QuerySender
			}

			sendEmail := new(types.SendingEmailToBroker)
			sendEmail.Sender = Sender
			sendEmail.Subject = "Create Account Success"
			sendEmail.Title = "Akun Anda Berhasil Di Buat"
			sendEmail.BodyText = `
				Hi, ` + user.FullName + `, berikut inform
				for reset password pleasasi akun anda:e click link below

				username: ` + user.Username + `
				email: ` + user.Email + `

				Tetap jaga rahasia akun anda, mohon untuk jangan diberikan kepada siapapun termasuk Admin.
			`
			sendEmail.Body = `
				<div class="font-family: Roboto-Regular, Helvetica, Arial, sans-serif; font-size: 14px; color: rgba(0, 0, 0, 0.87); padding-top: 20px; text-align: center;">Hi, ` + user.FullName + `, berikut informasi akun anda:</div>
				<br/>
				<div>username: <strong>` + user.FullName + `</strong></div>
				<div>email: <strong>` + user.Email + `</strong></div>` + `
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

		// return handlers.NotFoundErrorResponse(c, err)
	}

	if user.StatusAccount == 0 && isActive == false {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Your Account is not active or blocked"})
	}

	randString := handlers.RandomStringV1(32)

	user_data := handlers.SendSafeUserData(user, randString)

	if err := models.GenerateAuthKeyUser(user_data.IdAccount, user_data.AuthKey).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	// SESSION QUERY
	if session {
		err := middleware.SaveUserSession(c, user_data, false)
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

	if httponly { // HTTPONLY QUERY
		if domain == "" {
			return handlers.UnprocessableEntityErrorResponse(c, fmt.Errorf("need domain params"))
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
