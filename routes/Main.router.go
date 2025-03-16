package routes

import (
	"go-gerbang/middleware"
	"go-gerbang/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"

	// "github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

var baseConfig = basicauth.Config{
	Users: map[string]string{
		"admin": "@dmin9192",
	},
}

func MainRoutes(app *fiber.App) {
	// app.All("/info/*", filesystem.New(filesystem.Config{
	// 	Root:         embed.Dist(),
	// 	NotFoundFile: "index.html",
	// 	Index:        "index.html",
	// }))

	app.Get("/monitor", monitor.New(monitor.Config{Title: "GO GERBANG Monitor Page"}))

	// app.Static("/info", "./info")

	// GET CSRF TOKEN
	app.Get("/secure-gateway-c", middleware.CsrfProtection, services.IndexService)
	// app.Get("/api/secure-gateway-c", middleware.CsrfProtection, services.IndexService)
	// PROTECT
	// app.Get("/test-protect", middleware.Auth, services.ProtectService)

	app.Get("/check-migration", services.CheckMigrationStatus)
	app.Get("/migration", basicauth.New(baseConfig), services.MigrationService)
	app.Get("/migration-admin", basicauth.New(baseConfig), services.MigrateAdminUser)

	app.Get("/info", basicauth.New(baseConfig), services.InfoService)

	// PUB / SUB
	app.Post("/publish", services.PublishService)
	app.Get("/subscribe", services.SubscribeService)

	// MAIL
	app.Get("/check-mail", services.MailTesting)
	// app.Get("/send-mail", services.SendEmailHandler)

	services.SubscribeServiceEmail()
}
