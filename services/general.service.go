package services

import (
	"fmt"

	"go-gerbang/config"
	"go-gerbang/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

// @Summary Get CSRF Cookie
// @Description Get CSRF Cookie
// @Tags security
// @Accept json
// @Produce json
// @Router /secure-gateway-c [get]
func IndexService(c *fiber.Ctx) error {
	csrfToken, ok := c.Locals("token_csrf").(string)
	if !ok {
		return handlers.InternalServerErrorResponse(c, fmt.Errorf("error getting csrf"))
	}
	return handlers.SuccessResponse(c, "success getting csrf", csrfToken, nil)
}

func ProtectService(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"title": "Testing Protect Route"})
}

func InfoService(c *fiber.Ctx) error {
	var err error
	handlers.MapMicroService, err = handlers.LoadConfig(config.BasePath + config.ConfigPath)
	if err != nil {
		return err
	}

	done := make(chan bool)
	go handlers.WatchConfigFile(config.BasePath+config.ConfigPath, done)

	for i := range handlers.MapMicroService.Services {
		req := fasthttp.AcquireRequest()
		res := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(res)

		req.SetRequestURI(handlers.MapMicroService.Services[i].Url)

		handlers.MapMicroService.Services[i].Status = true

		client := &fasthttp.Client{}
		if err := client.Do(req, res); err != nil {
			// fmt.Printf("Error: %s\n", err)
			handlers.MapMicroService.Services[i].Status = false
		}

		if res.StatusCode() != fiber.StatusOK {
			handlers.MapMicroService.Services[i].Status = false
		}
	}

	return c.JSON(handlers.MapMicroService.Services)
}
