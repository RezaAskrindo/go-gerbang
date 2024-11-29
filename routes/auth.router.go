package routes

import (
	"go-gerbang/middleware"
	"go-gerbang/services"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	auth := app.Group("/api")
	auth.Post("/v1/auth/login", middleware.ValidateCaptcha, middleware.CsrfProtection, services.Login)

	auth.Post("/with-google", services.LoginWithGoogle)

	authSession := auth.Group("/v1/auth").Use(middleware.Auth)
	authSession.Get("/get-session", services.GetSessionJWT)
	authSession.Get("/auth-key/:token", services.AuthByJWT)
	authSession.Get("/logout", services.LogoutWeb)
}
