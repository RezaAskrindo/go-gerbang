package services

import (
	"context"
	"fmt"

	"go-gerbang/handlers"
	"go-gerbang/middleware"
	"go-gerbang/models"
	"go-gerbang/types"

	"github.com/gofiber/fiber/v3"
)

func GoogleOAuthLogin(c fiber.Ctx) error {
	query := new(types.LoginQuery)

	if err := c.Bind().Query(query); err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	if query.ClientId == "" {
		return handlers.UnprocessableEntityErrorResponse(c, fmt.Errorf("need client_id params"))
	}

	req := new(handlers.GoogleLoginRequest)

	if err := handlers.ParseBody(c, req); err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	// Verify token on backend
	claims, err := handlers.VerifyIdTokenGoogle(context.Background(), req.Token, query.ClientId)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid token",
			"error":   err.Error(),
		})
	}

	email := claims.Email
	name := claims.Name

	user := new(models.User)
	userExists := true

	if err := models.FindUserByIdentity(user, email, email, "", email); err != nil {
		userExists = false
		user.FullName = name
		user.Username = email
		user.Email = email
		user.IsGoogleAccount = 10
		user.StatusAccount = 10

		if query.CreateNew {
			if err := models.CreateUser(user); err.Error != nil {
				return handlers.ConflictErrorResponse(c, err.Error)
			}
		} else {
			return handlers.NotFoundErrorResponse(c, fmt.Errorf("User is not found"))
		}
	}

	if user.StatusAccount == 0 {
		return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("your account is not active or blocked"))
	}

	randString := handlers.RandomStringV1(24)
	userData := handlers.SendSafeUserData(user, randString)

	if err := models.GenerateAuthKeyUser(userData.IdAccount, userData.AuthKey).Error; err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	if query.Session { // SESSION QUERY
		err := middleware.SaveUserSession(c, userData, query.SingleLogin) // SINGLE LOGIN
		if err != nil {
			return handlers.InternalServerErrorResponse(c, err)
		}
	}

	if query.ValidateIp { // VALIDATE IP QUERY
		errValidate := handlers.ValidateUserLoginIp(userData, c)
		if errValidate != nil {
			return handlers.SuccessResponse(c, true, errValidate.Error(), userData, nil)
		}
	}

	refreshToken, err := handlers.GenerateTokenJWT(userData, true)
	if err != nil {
		return handlers.InternalServerErrorResponse(c, fmt.Errorf("failed to generate new refresh token"))
	}

	token, err := handlers.GenerateTokenJWT(userData, false)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to generate token",
		})
	}

	if query.HttpOnly { // HTTPONLY QUERY
		domain := query.Domain

		if domain == "" {
			domain = getMainDomain(c.Host())
		}

		middleware.SetAuthCookies(c, domain, refreshToken, token)

		return handlers.SuccessResponse(c, true, "Success Login for domain:"+domain, userData, nil)
	}

	return handlers.SuccessResponse(c, true, "Login successful", fiber.Map{
		"user":      userData,
		"token":     token,
		"isNewUser": !userExists,
	}, nil)
}
