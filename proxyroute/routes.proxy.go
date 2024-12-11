package proxyroute

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"go-gerbang/config"
	"go-gerbang/database"
	"go-gerbang/handlers"
	"go-gerbang/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/valyala/fasthttp"
)

func MainProxyRoutes(app *fiber.App) {
	// file, err := os.Open(config.BasePath + config.ConfigPath)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()

	// var MapMicroService []types.Service
	// err = json.NewDecoder(file).Decode(&MapMicroService)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	var err error
	// handlers.MapMicroService, err = handlers.LoadConfig(config.BasePath + config.ConfigPath)
	configProxy, err := handlers.LoadConfig(config.BasePath + config.ConfigPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	handlers.SaveToRedis("proxy-route", configProxy)

	configProxyJSON, err := database.RedisDb.Get(handlers.Ctx, "proxy-route").Result()
	if err != nil {
		log.Fatalf("Error loading redis: %v", err)
	}

	err = json.Unmarshal([]byte(configProxyJSON), &handlers.MapMicroService)
	if err != nil {
		log.Fatalf("could not deserialize config: %w", err)
	}

	done := make(chan bool)
	go handlers.WatchConfigFile(config.BasePath+config.ConfigPath, done)

	RegisterRoutes(app)
}

func RegisterRoutes(app *fiber.App) {
	handlers.MapMicroServiceMutex.RLock()
	defer handlers.MapMicroServiceMutex.RUnlock()

	proxy.WithClient(&fasthttp.Client{
		NoDefaultUserAgentHeader: true,
		DisablePathNormalizing:   true,
	})

	for _, data := range handlers.MapMicroService.Services {
		if data.AuthProtection && data.CsrfProtection {
			app.Use(data.Path, middleware.CsrfProtection, middleware.Auth, func(c *fiber.Ctx) error {
				path := c.OriginalURL()
				params := strings.TrimPrefix(path, data.Path)
				url := data.Url + params

				if err := proxy.DoTimeout(c, url, 30*time.Second); err != nil {
					log.Println("Error: %s\n", err.Error())
					return handlers.InternalServerErrorResponse(c, fmt.Errorf("endpoint is not running"))
				}

				c.Response().Header.Del(fiber.HeaderServer)
				return nil
			})
		} else {
			app.Use(data.Path, func(c *fiber.Ctx) error {
				path := c.OriginalURL()
				params := strings.TrimPrefix(path, data.Path)
				url := data.Url + params

				if err := proxy.Do(c, url); err != nil {
					log.Println("Error: %s\n", err.Error())
					return handlers.InternalServerErrorResponse(c, fmt.Errorf("endpoint is not running"))
				}

				c.Response().Header.Del(fiber.HeaderServer)
				return nil
			})
		}
	}
}

func RegisterDynamicRoutes(app *fiber.App) {
	app.Use(func(c *fiber.Ctx) error {
		handlers.MapMicroServiceMutex.RLock()
		defer handlers.MapMicroServiceMutex.RUnlock()

		for _, data := range handlers.MapMicroService.Services {
			if strings.HasPrefix(c.OriginalURL(), data.Path) {
				params := strings.TrimPrefix(c.OriginalURL(), data.Path)
				url := data.Url + params

				log.Println(data.AuthProtection)

				// if data.AuthProtection {
				// Apply AuthProtection and CsrfProtection
				// if err := middleware.CsrfProtection(c); err != nil {
				// 	return err
				// }
				// if err := middleware.Auth(c); err != nil {
				// 	return err
				// }
				// }
				middleware.Auth(c)

				// Proxy request
				if err := proxy.DoTimeout(c, url, 30*time.Second); err != nil {
					log.Println("Error: %s\n", err.Error())
					return handlers.InternalServerErrorResponse(c, fmt.Errorf("endpoint is not running"))
				}

				// c.Response().Header.Del(fiber.HeaderServer)
				return nil
			}
		}

		return c.Next() // If no matching route, proceed to the next middleware
	})
}
