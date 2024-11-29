package proxyroute

import (
	"log"
	"strings"
	"time"

	"go-gerbang/config"
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
	handlers.MapMicroService, err = handlers.LoadConfig(config.BasePath + config.ConfigPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	done := make(chan bool)
	go handlers.WatchConfigFile(config.BasePath+config.ConfigPath, done)

	proxy.WithClient(&fasthttp.Client{
		NoDefaultUserAgentHeader: true,
		DisablePathNormalizing:   true,
	})

	for _, data := range handlers.MapMicroService.Services {
		if data.AuthProtection {
			app.Use(data.Path, middleware.CsrfProtection, middleware.Auth, func(c *fiber.Ctx) error {
				path := c.OriginalURL()
				params := strings.TrimPrefix(path, data.Path)
				url := data.Url + params

				// if err := proxy.Do(c, url); err != nil {
				if err := proxy.DoTimeout(c, url, 30*time.Second); err != nil {
					log.Println("Error: %s\n", err.Error())
					// return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "endpoint is not running"})
					// return nil
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
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "endpoint is not running"})
				}

				c.Response().Header.Del(fiber.HeaderServer)
				return nil
			})

		}
	}
}
