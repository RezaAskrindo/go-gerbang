package routes

import (
	// "go-gerbang/embed"
	"go-gerbang/middleware"
	"go-gerbang/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	// "github.com/gofiber/fiber/v2/middleware/filesystem"
)

var baseConfig = basicauth.Config{
	Users: map[string]string{
		"admin": "@dmin9192",
	},
}

func MainRoutes(app *fiber.App) {
	// app.All("/embed/*", filesystem.New(filesystem.Config{
	// 	Root:         embed.Dist(),
	// 	NotFoundFile: "index.html",
	// 	Index:        "index.html",
	// }))

	// GET CSRF TOKEN
	app.Get("/secure-gateway-c", middleware.CsrfProtection, services.IndexService)
	app.Get("/secure-gateway-c-cookie", middleware.CsrfProtectionCookies, services.GetCSRFTokenService)
	// PROTECT
	app.Get("/test-protect", middleware.Auth, services.ProtectService)

	app.Get("/check-migration", services.CheckMigrationStatus)
	app.Get("/migration", basicauth.New(baseConfig), services.MigrationService)
	app.Get("/migration-admin", basicauth.New(baseConfig), services.MigrateAdminUser)

	app.Get("/info", basicauth.New(baseConfig), services.InfoService)

	// PUB / SUB
	app.Post("/publish", services.PublishService)
	app.Get("/subscribe", services.SubscribeService)

	// MAIL
	app.Get("/check-mail", services.MailTesting)

	// services.SubscribeServiceEmail()
	services.SubscribeEvent()
}
