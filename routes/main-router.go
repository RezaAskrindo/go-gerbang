package routes

import (
	// "go-gerbang/embed"
	"go-gerbang/middleware"
	"go-gerbang/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	// "github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func MainRoutes(app *fiber.App) {
	// app.All("/info/*", filesystem.New(filesystem.Config{
	// 	Root:         embed.Dist(),
	// 	NotFoundFile: "index.html",
	// 	Index:        "index.html",
	// }))

	app.Get("/metrics", monitor.New(monitor.Config{Title: "GO GERBANG Metrics Page"}))

	// app.Static("/info", "./info")

	// GET CSRF TOKEN
	app.Get("/secure-gateway-c", middleware.CsrfProtection, services.IndexService)
	app.Get("/api/secure-gateway-c", middleware.CsrfProtection, services.IndexService)
	// PROTECT
	app.Get("/test-protect", middleware.Auth, services.ProtectService)

	app.Get("/check-migration", services.CheckMigrationStatus)
	app.Get("/migration", services.MigrationService)

	app.Get("/info/micro-service", basicauth.New(basicauth.Config{
		Users: map[string]string{
			"admin": "9192",
		},
	}), services.InfoService)

	app.Get("/get-captcha", middleware.ValidateCaptcha, middleware.GenerateCaptcha)

	// auth := app.Group("/api/login")
	// auth.Post("/v1", middleware.ValidateCaptcha, middleware.CsrfProtection, services.Login)
	// auth.Post("/with-google", services.LoginWithGoogle)

	// authSession := app.Group("/api/session").Use(middleware.Auth)
	// authSession.Get("/get-session", services.GetSessionJWT)
	// authSession.Get("/auth-key/:token", services.AuthByJWT)
	// authSession.Get("/logout", services.LogoutWeb)
}
