package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"go-gerbang/broker"
	"go-gerbang/config"
	"go-gerbang/database"
	"go-gerbang/handlers"
	"go-gerbang/proxyroute"
	"go-gerbang/routes"

	// "go-gerbang/docs"

	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/circuitbreaker"
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
)

const (
	appName = "GO Gerbang"
)

var allowedOriginRegex *regexp.Regexp

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
		ServerHeader:          "",
		AppName:               appName,
		CaseSensitive:         true,
		DisableStartupMessage: true,
		ProxyHeader:           "X-Forwarded-For",
		// StrictRouting:         true,
		// Prefork:       true,
	})

	app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/live",
	}))

	app.Use(idempotency.New())

	allowedOriginRegex, _ = regexp.Compile(config.Config("ALLOWED_ORIGIN_REGEX"))
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return allowedOriginRegex.MatchString(origin)
		},
		AllowHeaders:     "Authorization, Content-Type, X-Sgcsrf-Token",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
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
		Max:        1000,
		Expiration: 60 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("X-forwarded-for")
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendString("be slow bro...")
		},
	}))

	app.Use(earlydata.New())

	// docs.SwaggerInfo.Title = appName
	// docs.SwaggerInfo.Description = "This is an API for GO GERBANG Apigateway"
	// docs.SwaggerInfo.Version = "1.0"
	// docs.SwaggerInfo.Host = "localhost:" + config.Config("PORT_APIGATEWAY")
	// docs.SwaggerInfo.BasePath = "/"

	// app.Get("/swagger/*", swagger.HandlerDefault)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Send([]byte("Welcome to GO GERBANG - by Muhammad Reza"))
	})

	routes.MainRoutes(app)

	cb := circuitbreaker.New(circuitbreaker.Config{
		FailureThreshold: 3,               // Max failures before opening the circuit
		Timeout:          5 * time.Second, // Wait time before retrying
		SuccessThreshold: 2,               // Required successes to move back to closed state
	})

	app.Get("/health/circuit", cb.HealthHandler())
	app.Get("/metrics/circuit", func(c *fiber.Ctx) error {
		return c.JSON(cb.GetStateStats())
	})

	app.Use(circuitbreaker.Middleware(cb))

	ctx := context.Background()
	err = handlers.InitLogger(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize zap logger: %v", err)
	}
	routes.AuthRoutes(app)
	proxyroute.MainProxyRoutes(app)

	app.Use("*", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"code": 400, "status": "error", "message": "Not Found Services"})
	})

	fmt.Println("✅ server running " + config.Config("PORT_APIGATEWAY"))
	if err := app.Listen(config.Config("PORT_APIGATEWAY")); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
