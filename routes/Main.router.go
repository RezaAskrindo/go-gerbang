package routes

import (
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

	app.Get("/monitor", monitor.New(monitor.Config{Title: "GO GERBANG Monitor Page"}))

	// app.Static("/info", "./info")

	// GET CSRF TOKEN
	app.Get("/secure-gateway-c", middleware.CsrfProtection, services.IndexService)
	// app.Get("/api/secure-gateway-c", middleware.CsrfProtection, services.IndexService)
	// PROTECT
	// app.Get("/test-protect", middleware.Auth, services.ProtectService)

	app.Get("/check-migration", services.CheckMigrationStatus)
	app.Get("/migration", services.MigrationService)

	app.Get("/info/micro-service", basicauth.New(basicauth.Config{
		Users: map[string]string{
			"admin": "9192",
		},
	}), services.InfoService)

	// PUB / SUB
	app.Post("/publish", services.PublishService)
	app.Get("/subscribe", services.SubscribeService)

	services.SubscribeServiceEmail()
}
