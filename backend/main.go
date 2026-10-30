package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"go-gerbang/broker"
	"go-gerbang/config"
	"go-gerbang/database"
	"go-gerbang/handlers"
	"go-gerbang/middleware"
	"go-gerbang/proxyroute"
	"go-gerbang/routes"
	"go-gerbang/services"

	"github.com/KimMachineGun/automemlimit/memlimit"
	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/v3/circuitbreaker"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/encryptcookie"
	"github.com/gofiber/fiber/v3/middleware/etag"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/idempotency"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"go.uber.org/automaxprocs/maxprocs"
)

func main() {
	// runtime.GOMAXPROCS(coreCPU) // change to maxprocs
	if _, err := maxprocs.Set(); err != nil {
		log.Printf("automaxprocs: %v", err)
	}
	_, err := memlimit.Set(
		memlimit.WithRatio(0.50),
		memlimit.WithMin(20*1024*1024),
		memlimit.WithProvider(
			memlimit.ApplyFallback(
				memlimit.FromCgroup,
				memlimit.FromSystem,
			),
		),
	)
	if err != nil {
		debug.SetMemoryLimit(config.MemoryLimit << 20) // change to memlimit
		log.Printf("automemlimit: %v", err)
	}

	debug.SetGCPercent(50)
	debug.SetMaxStack(config.MemoryLimit << 20)

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

	database.ConnectGormDB()
	broker.StartingNatsClient()
	defer broker.NatsClient.Drain()

	app := fiber.New(fiber.Config{
		JSONEncoder:       json.Marshal,
		JSONDecoder:       json.Unmarshal,
		BodyLimit:         50 * 1024 * 1024, // this is the default limit of 50MB
		ServerHeader:      config.AppName,
		AppName:           config.AppName,
		CaseSensitive:     true,
		ProxyHeader:       "X-Forwarded-For",
		ReduceMemoryUsage: true,
		// NEW CONFIG
		ReadTimeout:     50 * time.Second,
		WriteTimeout:    50 * time.Second,
		IdleTimeout:     50 * time.Second,
		Concurrency:     2048,
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
	})

	app.Use(idempotency.New(idempotency.Config{
		Storage: middleware.StorageIdempotency,
		Next: func(c fiber.Ctx) bool {
			m := c.Method()
			return m != fiber.MethodPost
		},
		KeyHeader: "Gerbang-Idempotency-Key",
		Lifetime:  24 * time.Hour,
	}))

	app.Use(requestid.New(requestid.Config{
		Header: "Gerbang-Request-ID",
		Generator: func() string {
			return "static-id"
		},
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: config.GetTrustedOrigins(),
		AllowHeaders: []string{
			"Authorization",
			"Content-Type",
			"X-Sgcsrf-Token",
			"Gerbang-Idempotency-Key",
			"Gerbang-Request-ID",
		},
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		ExposeHeaders:    []string{"Gerbang-Request-ID", "Retry-After", "ETag", "Location"},
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

	app.Use(limiter.New(limiter.Config{
		Storage:    middleware.StorageLimiter,
		Max:        500,
		Expiration: 60 * time.Second,
		KeyGenerator: func(c fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c fiber.Ctx) error {
			return c.SendString("be slow bro...")
		},
	}))

	middleware.InitCSRF()

	app.Get("/", services.MainService)

	ctx := context.Background()
	err = handlers.InitLogger(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize zap logger: %v", err)
	}

	app.All("/live", healthcheck.New())

	cb := circuitbreaker.New(circuitbreaker.Config{
		FailureThreshold: 7,
		Timeout:          15 * time.Second,
		SuccessThreshold: 4,
		IsFailure: func(c fiber.Ctx, err error) bool {
			return c.Response().StatusCode() >= http.StatusInternalServerError
		},
	})

	app.Get("/health/circuit", cb.HealthHandler())
	app.Get("/metrics/circuit", func(c fiber.Ctx) error {
		return c.JSON(cb.GetStateStats())
	})

	app.Use(circuitbreaker.Middleware(cb))

	proxyroute.MainProxyRoutes(app)

	routes.MainRoutes(app)
	routes.AuthRoutes(app)

	app.Use("*", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"code": 400, "status": "error", "message": "Not Found Services"})
	})

	fmt.Println("[INFO] Server running " + config.APP_PORT)

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- app.Listen(config.APP_PORT, fiber.ListenConfig{
			EnablePrefork:         false, // required — shutdown doesn't work with prefork
			DisableStartupMessage: true,
		})
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		log.Println("shutdown signal received, draining...")
	case err := <-serverErr:
		stop()
		if err != nil {
			log.Fatalf("Fiber server failed: %v", err)
		}
		return
	}
	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Printf("Fiber shutdown error: %v", err)
	}

	if err := <-serverErr; err != nil {
		log.Printf("Listen error after shutdown: %v", err)
	}

	natsDone := make(chan struct{})
	go func() {
		natsServer.Shutdown()
		close(natsDone)
	}()
	select {
	case <-natsDone:
		log.Println("NATS server stopped")
	case <-time.After(5 * time.Second):
		log.Println("NATS shutdown timed out")
	}

	log.Println("clean shutdown complete")
}
