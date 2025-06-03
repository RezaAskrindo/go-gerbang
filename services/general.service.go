package services

import (
	"log"
	"time"

	"go-gerbang/config"
	"go-gerbang/handlers"
	"go-gerbang/middleware"

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
	csrfToken, _ := c.Locals(middleware.CsrfContextKey).(string)
	return handlers.SuccessResponse(c, true, "success get csrf token", csrfToken, nil)
	// if !ok {
	// 	return handlers.InternalServerErrorResponse(c, fmt.Errorf("error getting csrf"))
	// }
	// return c.SendString(csrfToken)
}

func GetCSRFTokenService(c *fiber.Ctx) error {
	_ = c.Locals(middleware.CsrfContextKey).(string)
	return c.SendStatus(fiber.StatusNoContent)
}

func ProtectService(c *fiber.Ctx) error {
	return c.SendString("Testing Protect Route")
}

func InfoService(c *fiber.Ctx) error {
	var err error
	handlers.MapMicroService, err = handlers.LoadConfig(config.BasePath + config.ConfigPath)
	if err != nil {
		return err
	}

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

func SendGetRequest(url string) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodGet)
	resp := fasthttp.AcquireResponse()
	readTimeout, _ := time.ParseDuration("500ms")
	writeTimeout, _ := time.ParseDuration("500ms")
	maxIdleConnDuration, _ := time.ParseDuration("1h")
	client := &fasthttp.Client{
		ReadTimeout:                   readTimeout,
		WriteTimeout:                  writeTimeout,
		MaxIdleConnDuration:           maxIdleConnDuration,
		NoDefaultUserAgentHeader:      true, // Don't send: User-Agent: fasthttp
		DisableHeaderNamesNormalizing: true, // If you set the case on your headers correctly you can enable this
		DisablePathNormalizing:        true,
		// increase DNS cache time to an hour instead of default minute
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
	}
	err := client.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	if err != nil {
		log.Printf("ERR Connection error: %v\n", err)
	}
	fasthttp.ReleaseResponse(resp)
}
