package services

import (
	"go-gerbang/handlers"
	"go-gerbang/middleware"
	"go-gerbang/models"
	"go-gerbang/types"

	"github.com/gofiber/fiber/v2"
)

func LoginWithGoogle(c *fiber.Ctx) error {
	session := c.QueryBool("session")
	validate_ip := c.QueryBool("validate_ip")

	b := new(types.GoogleLogin)

	if err := handlers.ParseBody(c, b); err != nil {
		return handlers.BadRequestErrorResponse(c, err)
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

	return c.JSON(fiber.Map{"success": true, "message": "Success Login", "token": token, "data": user_data})
}
