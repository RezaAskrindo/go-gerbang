package routes

import (
	"go-gerbang/middleware"
	"go-gerbang/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

var baseConfig = basicauth.Config{
	Users: map[string]string{
		"admin": "@dmin9192",
	},
}

func MainRoutes(app *fiber.App) {
	// GET CSRF TOKEN
	app.Get("/secure-gateway-c", middleware.CsrfProtection, services.IndexService)
	app.Get("/secure-gateway-c-cookie", middleware.CsrfProtectionCookies, services.GetCSRFTokenService)
	// PROTECT
	app.Get("/test-protect", middleware.Auth, services.ProtectService)

	app.Get("/check-migration", services.CheckMigrationStatus)
	app.Get("/migration", basicauth.New(baseConfig), services.MigrationService)
	app.Get("/migration-admin", basicauth.New(baseConfig), services.MigrateAdminUser)

	app.Get("/info", services.InfoService)

	// PUB / SUB
	app.Post("/publish", services.PublishService)
	app.Get("/subscribe", services.SubscribeService)

	// MAIL
	app.Get("/check-mail", services.MailTesting)

	// SERVICE
	app.Post("/restart", services.RestartHandler)
	app.Post("/config-file", services.HandleConfigFile)

	app.Get("/log-stats-proxy", services.GetStatsLogProxy)
	app.Get("/metrics", monitor.New(monitor.Config{APIOnly: true}))

	// services.SubscribeServiceEmail()
	services.SubscribeEvent()
}
