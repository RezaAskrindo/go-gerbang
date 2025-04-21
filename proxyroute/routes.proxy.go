package proxyroute

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"go-gerbang/config"
	"go-gerbang/handlers"
	"go-gerbang/middleware"
	"go-gerbang/types"

	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/gofiber/contrib/casbin"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/valyala/fasthttp"
)

func MainProxyRoutes(app *fiber.App) {
	var err error

	// USING FILE
	handlers.MapMicroService, err = handlers.LoadConfig(config.BasePath + config.ConfigPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// BEGIN USING REDIS
	// configProxy, err := handlers.LoadConfig(config.BasePath + config.ConfigPath)
	// if err != nil {
	// 	log.Fatalf("Error loading config: %v", err)
	// }

	// handlers.SaveToRedis("proxy-route", configProxy)

	// configProxyJSON, err := database.RedisDb.Get(handlers.Ctx, "proxy-route").Result()
	// if err != nil {
	// 	log.Fatalf("Error loading redis: %v", err)
	// }

	// err = json.Unmarshal([]byte(configProxyJSON), &handlers.MapMicroService)
	// if err != nil {
	// 	log.Fatalf("could not deserialize config: %v", err)
	// }
	// END USING REDIS

	// WATCH CONFIG FILE
	// done := make(chan bool)
	// go handlers.WatchConfigFile(config.BasePath+config.ConfigPath, done)

	RegisterRoutes(app)
}

func RegisterRoutes(app *fiber.App) {
	handlers.MapMicroServiceMutex.RLock()
	defer handlers.MapMicroServiceMutex.RUnlock()

	proxy.WithClient(&fasthttp.Client{
		NoDefaultUserAgentHeader: true,
		DisablePathNormalizing:   true,
	})

	// RBAC PROTECTION
	authz := casbin.New(casbin.Config{
		ModelFilePath: config.BasePath + config.Config("CONFIG_PATH_CASBIN_MODEL"),
		PolicyAdapter: fileadapter.NewAdapter(config.BasePath + config.Config("CONFIG_PATH_CASBIN_POLICY")),
		Lookup: func(c *fiber.Ctx) string {
			statusAccount, ok := c.Locals("status_account").(int8)
			if !ok {
				return ""
			}
			sub := strconv.Itoa(int(statusAccount))

			return sub
		},
		Unauthorized: func(c *fiber.Ctx) error {
			return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("your role don't have access"))
		},
	})

	for _, service := range handlers.MapMicroService.Services {
		if service.AuthProtection && service.CsrfProtection && service.RbacProtection {
			app.Use(service.Path, middleware.CsrfProtection, middleware.Auth, authz.RoutePermission(), proxyHandler(service))
		} else if service.AuthProtection && service.CsrfProtection {
			app.Use(service.Path, middleware.CsrfProtection, middleware.Auth, proxyHandler(service))
		} else {
			app.Use(service.Path, proxyHandler(service))
		}
	}
}

func proxyHandler(service types.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.OriginalURL()
		params := strings.TrimPrefix(path, service.Path)
		url := service.Url + params

		if err := proxy.DoTimeout(c, url, 30*time.Second); err != nil {
			log.Printf("error: %s\n", err.Error())
			return handlers.InternalServerErrorResponse(c, fmt.Errorf("endpoint is not running"))
		}

		c.Response().Header.Del(fiber.HeaderServer)
		return nil
	}
}
