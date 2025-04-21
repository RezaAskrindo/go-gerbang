package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"go-gerbang/broker"
	"go-gerbang/config"
	"go-gerbang/database"
	"go-gerbang/proxyroute"
	"go-gerbang/routes"

	// "go-gerbang/docs"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/earlydata"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	// "github.com/gofiber/swagger"
	// "github.com/gofiber/contrib/fiberzap/v2"
	// "go.uber.org/zap"
)

const (
	appName = "GO Gerbang"
)

// @termsOfService http://swagger.io/terms/
// @contact.name Muhammad Reza
// @contact.email m.reza911992@gmail.com

func main() {
	logFile, err := os.OpenFile("go-gerbang.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v\n", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	natsServer, err := broker.StartingNatsServer()
	if err != nil {
		log.Fatalf("Error starting NATS server: %v", err)
	}
	defer natsServer.Shutdown()

	database.ConnectGormDB()
	broker.StartingNatsClient()

	app := fiber.New(fiber.Config{
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		BodyLimit:             100 * 1024 * 1024, // this is the default limit of 100MB
		ServerHeader:          appName,
		AppName:               appName,
		CaseSensitive:         true,
		DisableStartupMessage: true,
		// StrictRouting:         true,
		// Prefork:       true,
	})
	// defer app.Shutdown()c

	app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/",
	}))

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

	config.SecureCookies, err = strconv.ParseBool(config.SecureCookiesString)
	if err != nil {
		config.SecureCookies = false
	}

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

	app.Use(earlydata.New())

	// logger, _ := zap.NewProduction()
	// defer logger.Sync()

	// app.Use(fiberzap.New(fiberzap.Config{
	// 	Logger: logger,
	// }))

	// docs.SwaggerInfo.Title = appName
	// docs.SwaggerInfo.Description = "This is an API for GO GERBANG Apigateway"
	// docs.SwaggerInfo.Version = "1.0"
	// docs.SwaggerInfo.Host = "localhost:" + config.Config("PORT_APIGATEWAY")
	// docs.SwaggerInfo.BasePath = "/"

	// app.Get("/swagger/*", swagger.HandlerDefault)

	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.Send([]byte("Welcome to GO GERBANG"))
	// })

	routes.MainRoutes(app)
	routes.AuthRoutes(app)

	proxyroute.MainProxyRoutes(app)

	app.Use("*", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"code": 400, "status": "error", "message": "Not Found Services"})
	})

	fmt.Println("server running " + config.Config("PORT_APIGATEWAY"))
	if err := app.Listen(config.Config("PORT_APIGATEWAY")); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
