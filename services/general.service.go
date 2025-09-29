package services

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"

	"go-gerbang/config"
	"go-gerbang/handlers"
	"go-gerbang/middleware"
	"go-gerbang/models"
	"go-gerbang/proxyroute"
	"go-gerbang/types"

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

		if err := proxyroute.ProxyClient.Do(req, res); err != nil {
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
		log.Printf("ERR HTTP Connection error: %v\n", err)
	}
	fasthttp.ReleaseResponse(resp)
}

func RestartHandler(c *fiber.Ctx) error {
	go func() {
		time.Sleep(1 * time.Second)

		// Detect restart mode from ENV (optional override)
		mode := config.Config("RESTART_MODE") // "exit", "exec", or "auto"

		// Docker / K8s / systemd case
		if mode == "exit" || os.Getenv("IN_DOCKER") == "true" {
			os.Exit(0)
			return
		}

		// Auto-detect by OS
		if runtime.GOOS == "windows" {
			// 🔹 Windows: spawn a new process, then exit
			exe, err := os.Executable()
			if err != nil {
				panic(err)
			}
			args := os.Args[1:]
			cmd := exec.Command(exe, args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin

			if err := cmd.Start(); err != nil {
				panic(err)
			}
			os.Exit(0)

		} else {
			// 🔹 Linux/macOS: replace process in-place
			exe, err := os.Executable()
			if err != nil {
				panic(err)
			}
			args := os.Args
			env := os.Environ()

			if err := syscall.Exec(exe, args, env); err != nil {
				panic(err)
			}
		}
	}()

	return c.JSON(fiber.Map{
		"message": "Service restarting...",
		"os":      runtime.GOOS,
	})
}

func GetStatsLogProxy(c *fiber.Ctx) error {
	layout := "2006-01-02T15:04:05.000Z07:00"

	isDetail := c.QueryBool("detail")

	fromStr := c.Query("from")
	toStr := c.Query("to")

	from := time.Now().AddDate(0, 0, -30)
	to := time.Now()

	if fromStr != "" {
		parsedFrom, err := time.Parse(layout, fromStr)
		if err != nil {
			return handlers.BadRequestErrorResponse(c, fmt.Errorf("invalid from param"))
		}
		from = parsedFrom
	}

	if toStr != "" {
		parsedTo, err := time.Parse(layout, toStr)
		if err != nil {
			return handlers.BadRequestErrorResponse(c, fmt.Errorf("invalid from param"))
		}
		to = parsedTo
	}

	if isDetail {
		d := &[]models.LogProxy{}
		service := c.Query("service")
		method := c.Query("method")
		path := c.Query("path")
		status := c.Query("status")

		err := models.FindLogProxy(d, service, method, path, status, from, to).Error
		if err != nil {
			return handlers.InternalServerErrorResponse(c, err)
		}

		return handlers.SuccessResponse(c, true, "success to get detail log proxy", d, nil)
	} else {
		d := &[]models.PathStats{}

		err := models.FindStatsLogProxy(d, from, to).Error
		if err != nil {
			return handlers.InternalServerErrorResponse(c, err)
		}

		return handlers.SuccessResponse(c, true, "success to get stats log proxy", d, nil)
	}
}

func HandleConfigFile(c *fiber.Ctx) error {
	u := new(types.ConfigServices)

	if err := handlers.ParseBody(c, u); err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	if err := handlers.SaveConfig(config.BasePath+config.ConfigPath, u); err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	return handlers.SuccessResponse(c, true, "success update config file", u, nil)
}
