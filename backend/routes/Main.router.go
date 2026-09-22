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
	app.Get("/check-migration", services.CheckMigrationStatus)

	app.Get("/log-stats-proxy", services.GetStatsLogger)
	app.Get("/metrics", monitor.New(monitor.Config{APIOnly: true}))

	app.Get("/Configuration/:group", services.GetConfigurationByGroup)

	app.Post("/publish", services.PublishService)
	services.SubscribeEvent()

	// GET CSRF TOKEN
	app.Get("/secure-gateway-c", middleware.CsrfProtection, services.IndexService)
	app.Get("/secure-gateway-c-cookie", middleware.CsrfProtectionCookies, services.GetCSRFTokenService)

	// AUTH BASE
	app.Get("/migration", basicauth.New(baseConfig), services.MigrationService)
	app.Post("/migration-admin", basicauth.New(baseConfig), services.MigrateAdminUser)

	// SERVICE AUTH with JWT
	app.Get("/test-protect", middleware.Auth, services.ProtectService)

	app.Get("/info", middleware.CsrfProtection, middleware.Auth, services.InfoService)

	app.Get("/check-local-service", middleware.CsrfProtection, middleware.Auth, services.CheckLocalService)
	app.All("/proxy-local-service", middleware.CsrfProtection, middleware.Auth, services.ProxyLocalService)

	app.Post("/restart", middleware.CsrfProtection, middleware.Auth, services.RestartHandler)
	app.Post("/config-file", middleware.CsrfProtection, middleware.Auth, services.HandleConfigFile)

	app.Post("/upload-file", middleware.Auth, services.HandleFileUpload)

	app.Get("/Configuration/:group/execute", middleware.CsrfProtection, middleware.Auth, services.ConfigExecuteScript)
	app.Post("/Configuration", middleware.CsrfProtection, middleware.Auth, services.UpsertConfiguration)
	app.Delete("/Configuration/:group", middleware.CsrfProtection, middleware.Auth, services.DeleteConfiguration)
}
