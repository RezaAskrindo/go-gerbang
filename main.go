package main

import (
	"log"
	"time"

	"go-gerbang/config"
	"go-gerbang/database"
	"go-gerbang/docs"
	"go-gerbang/middleware"

	"go-gerbang/proxyroute"
	"go-gerbang/routes"

	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/earlydata"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"go.uber.org/zap"
)

const (
	appName = "GO Gerbang"
)

// @termsOfService http://swagger.io/terms/
// @contact.name Muhammad Reza
// @contact.email m.reza911992@gmail.com

func main() {
	database.ConnectGormDB()

	app := fiber.New(fiber.Config{
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
		BodyLimit:     100 * 1024 * 1024, // this is the default limit of 100MB
		ServerHeader:  appName,
		AppName:       appName,
		CaseSensitive: true,
		StrictRouting: true,
		// Prefork:       true,
		// DisableStartupMessage: true,
	})

	// defer app.Shutdown()

	app.Use(idempotency.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.Config("ALLOW_ORIGINS"),
		AllowCredentials: true,
	}))

	app.Use(helmet.New(helmet.Config{
		CrossOriginOpenerPolicy:   "cross-origin",
		CrossOriginResourcePolicy: "cross-origin",
	}))

	app.Use(recover.New())

	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: config.Config("KEY_COOKIE_APIGATEWAY"),
	}))

	app.Use(etag.New())

	app.Use(requestid.New())

	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1" // limit will apply to this IP
		},
		Max:        500,
		Expiration: 60 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("X-forwarded-for")
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendString("be slow bro...")
		},
	}))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))

	app.Use(earlydata.New())

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger,
	}))

	docs.SwaggerInfo.Title = appName
	docs.SwaggerInfo.Description = "This is an API for GO GERBANG Apigateway"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:" + config.Config("PORT_APIGATEWAY")
	docs.SwaggerInfo.BasePath = "/"

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Send([]byte("Welcome to GO GERBANG"))
		// return c.Status(fiber.StatusOK).JSON(fiber.Map{"code": 200, "status": "live", "message": config.Config("TEST_SCRIPT")})
	})

	proxyroute.MainProxyRoutes(app)
	routes.MainRoutes(app)
	routes.AuthRoutes(app)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"code": 400, "status": "error", "message": "Not Found Services"})
	})

	log.Fatal(app.Listen(config.Config("PORT_APIGATEWAY")))

	middleware.SessionStore.Storage.Close()
}
