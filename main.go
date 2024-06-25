package main

import (
	"log"
	"time"

	"go-gerbang/config"
	"go-gerbang/database"
	"go-gerbang/middleware"

	"go-gerbang/proxyroute"
	"go-gerbang/routes"

	"github.com/bytedance/sonic"
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
)

func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder:   sonic.Marshal,
		JSONDecoder:   sonic.Unmarshal,
		BodyLimit:     100 * 1024 * 1024, // this is the default limit of 100MB
		CaseSensitive: true,
		ServerHeader:  "AZER CORP",
		AppName:       "Gateway App v1.0.0",
		// Prefork:       true,
		// StrictRouting: true,
		// DisableStartupMessage: true,
	})

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

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"code": 200, "status": "live", "message": config.Config("TEST_SCRIPT")})
	})

	database.ConnectGormDB()

	routes.MainRoutes(app)
	proxyroute.MainProxyRoutes(app)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"code": 400, "status": "error", "message": "Not Found Services"})
	})

	log.Fatal(app.Listen(config.Config("PORT_APIGATEWAY")))

	middleware.SessionStore.Storage.Close()
}
