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

	"github.com/KimMachineGun/automemlimit/memlimit"
	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/v3/circuitbreaker"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/earlydata"
	"github.com/gofiber/fiber/v3/middleware/encryptcookie"
	"github.com/gofiber/fiber/v3/middleware/etag"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"go.uber.org/automaxprocs/maxprocs"
	// "github.com/gofiber/fiber/v3/middleware/idempotency"
	// "github.com/gofiber/fiber/v3/middleware/requestid"
)

// NOTE: FOR LOW VPS
const (
	appName = "GO Gerbang"
	// coreCPU     = 1   // for VCPU is 1
	memoryLimit = 128 // for Memory Limit
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
	// _, err := memlimit.SetGoMemLimitWithOpts(memlimit.WithRatio(0.50));
	if err != nil {
		debug.SetMemoryLimit(memoryLimit << 20) // change to memlimit
		log.Printf("automemlimit: %v", err)
	}

	debug.SetGCPercent(50)
	debug.SetMaxStack(memoryLimit << 20)

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
	defer broker.NatsClient.Drain()

	app := fiber.New(fiber.Config{
		JSONEncoder:       json.Marshal,
		JSONDecoder:       json.Unmarshal,
		BodyLimit:         50 * 1024 * 1024, // this is the default limit of 50MB
		ServerHeader:      appName,
		AppName:           appName,
		CaseSensitive:     true,
		ProxyHeader:       "X-Forwarded-For",
		ReduceMemoryUsage: true,
		// NEW CONFIG
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	})

	// app.Use(idempotency.New(idempotency.Config{
	// 	Storage: middleware.StorageIdempotency,
	// }))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.GetTrustedOrigins(),
		AllowHeaders:     []string{"Authorization", "Content-Type", "X-Sgcsrf-Token"},
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
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

	// app.Use(requestid.New())

	app.Use(limiter.New(limiter.Config{
		Storage: middleware.StorageLimiter,
		Next: func(c fiber.Ctx) bool {
			return c.IP() == "127.0.0.1" // limit will apply to this IP
		},
		Max:        1000,
		Expiration: 60 * time.Second,
		KeyGenerator: func(c fiber.Ctx) string {
			// return c.Get("X-forwarded-for")
			return c.IP()
		},
		LimitReached: func(c fiber.Ctx) error {
			return c.SendString("be slow bro...")
		},
	}))

	app.Use(earlydata.New())

	middleware.InitCSRF()

	app.Get("/", func(c fiber.Ctx) error {
		return c.Send([]byte("Welcome to GO GERBANG API GATEWAY - by Muhammad Reza"))
	})

	ctx := context.Background()
	err = handlers.InitLogger(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize zap logger: %v", err)
	}

	proxyroute.MainProxyRoutes(app)

	app.All("/live", healthcheck.New())

	cb := circuitbreaker.New(circuitbreaker.Config{
		FailureThreshold: 3,               // Max failures before opening the circuit
		Timeout:          5 * time.Second, // Wait time before retrying
		SuccessThreshold: 2,               // Required successes to move back to closed state
		IsFailure: func(c fiber.Ctx, err error) bool {
			return c.Response().StatusCode() >= http.StatusInternalServerError
		},
	})

	app.Get("/health/circuit", cb.HealthHandler())
	app.Get("/metrics/circuit", func(c fiber.Ctx) error {
		return c.JSON(cb.GetStateStats())
	})

	app.Use(circuitbreaker.Middleware(cb))
	routes.MainRoutes(app)
	routes.AuthRoutes(app)

	// For reducing memory and move to caddy
	// app.Get("/*", static.New("./web"))
	// SPA fallback — catches everything else
	// app.Get("*", static.New("./web/index.html"))

	app.Use("*", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"code": 400, "status": "error", "message": "Not Found Services"})
	})

	fmt.Println("[INFO] Server running " + config.APP_PORT)
	// if err := app.Listen(config.APP_PORT, fiber.ListenConfig{
	// 	EnablePrefork:         false,
	// 	DisableStartupMessage: true,
	// }); err != nil {
	// 	log.Fatalf("Error starting server: %v", err)
	// }

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
		natsServer.Shutdown() // don't skip cleanup on the fatal path
		if err != nil {
			log.Fatalf("Fiber server failed: %v", err)
		}
		return
	}
	stop() // restore default signal handling — a second Ctrl+C now force-exits

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Printf("Fiber shutdown error: %v", err)
	}

	// Wait for Listen to actually return (it should return nil after shutdown)
	if err := <-serverErr; err != nil {
		log.Printf("Listen error after shutdown: %v", err)
	}

	// NATS shutdown with its own timeout
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
