package routes

import (
	"go-gerbang/middleware"
	"go-gerbang/services"

	"github.com/gofiber/contrib/v3/monitor"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/basicauth"
)

var baseConfig = basicauth.Config{
	Users: map[string]string{
		"admin": "{SHA256}ITby5cH1QqwpO5S2UvXJLErtCVlFyVewpW97RXHUNKI=", //@dmin9192
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
	app.Post("/migration-admin", basicauth.New(baseConfig), services.MigrateAdminUser)

	app.Get("/info", services.InfoService)

	app.Get("/check-local-service", services.CheckLocalService)
	app.Get("/proxy-local-service", services.ProxyLocalService)

	app.Get("/Configuration/execute", services.ConfigExecuteScript)
	app.Get("/Configuration/:group", services.GetConfigurationByGroup)
	app.Post("/Configuration", services.UpsertConfiguration)
	app.Delete("/Configuration/:group", services.DeleteConfiguration)

	// PUB / SUB
	app.Post("/publish", services.PublishService)
	app.Get("/subscribe", services.SubscribeService)

	// MAIL
	app.Get("/check-mail", services.MailTesting)

	// SERVICE
	app.Post("/restart", services.RestartHandler)
	app.Post("/config-file", services.HandleConfigFile)
	app.Post("/upload-file", services.HandleFileUpload)

	app.Get("/log-stats-proxy", services.GetStatsLogger)
	app.Get("/metrics", monitor.New(monitor.Config{APIOnly: true}))

	// services.SubscribeServiceEmail()
	services.SubscribeEvent()
}
