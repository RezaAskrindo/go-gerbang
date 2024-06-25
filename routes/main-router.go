package routes

import (
	"go-gerbang/middleware"
	"go-gerbang/services"

	"github.com/gofiber/fiber/v2"
)

func MainRoutes(app *fiber.App) {
	// GET CSRF TOKEN
	app.Get("/secure-gateway-c", middleware.CsrfProtection, services.IndexService)
	app.Get("/api/secure-gateway-c", middleware.CsrfProtection, services.IndexService)
	// PROTECT
	app.Get("/test-protect", middleware.Auth, services.ProtectService)

	app.Get("/migration", services.MigrationService)

	app.Get("/info/micro-service", services.InfoService)

	app.Get("/get-captcha", middleware.ValidateCaptcha, middleware.GenerateCaptcha)

	auth := app.Group("/api/login")
	auth.Post("/v1", middleware.ValidateCaptcha, middleware.CsrfProtection, services.Login)
	auth.Post("/with-google", services.LoginWithGoogle)

	authSession := app.Group("/api/session").Use(middleware.Auth)
	authSession.Get("/get-session", services.GetSessionJWT)
	authSession.Get("/auth-key/:token", services.AuthByJWT)
	authSession.Get("/logout", services.LogoutWeb)
}
