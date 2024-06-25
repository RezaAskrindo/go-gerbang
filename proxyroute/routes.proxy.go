package proxyroute

import (
	"encoding/json"
	"fmt"
	"go-gerbang/config"
	"go-gerbang/middleware"
	"go-gerbang/types"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

func MainProxyRoutes(app *fiber.App) {
	file, err := os.Open(config.BasePath + config.ConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var MapMicroService []types.ValueMicroService
	err = json.NewDecoder(file).Decode(&MapMicroService)
	if err != nil {
		log.Fatal(err)
	}

	for _, data := range MapMicroService {
		if data.AuthProtection == true {
			app.Use(data.Path, middleware.CsrfProtection, middleware.Auth, func(c *fiber.Ctx) error {
				path := c.OriginalURL()
				params := strings.TrimPrefix(path, data.Path)
				url := data.Url + params

				// if err := proxy.DoTimeout(c, url, 30*time.Second); err != nil {
				if err := proxy.Do(c, url); err != nil {
					fmt.Printf("Error: %s\n", url)
					// return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "endpoint is not running"})
					// return nil
				}

				c.Response().Header.Del(fiber.HeaderServer)
				return nil
			})
		} else {
			app.Use(data.Path, middleware.CsrfProtection, func(c *fiber.Ctx) error {
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
