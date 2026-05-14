package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"

	"go-gerbang/config"
	"go-gerbang/handlers"
	"go-gerbang/models"
	"go-gerbang/proxyroute"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/proxy"
)

func IndexService(c fiber.Ctx) error {
	csrfToken := csrf.TokenFromContext(c)
	return handlers.SuccessResponse(c, true, "success get csrf token", csrfToken, nil)
}

func GetCSRFTokenService(c fiber.Ctx) error {
	csrfToken := csrf.TokenFromContext(c)
	return handlers.SuccessResponse(c, true, "success get csrf token", csrfToken, nil)
}

func ProtectService(c fiber.Ctx) error {
	return c.SendString("Testing Protect Route")
}

func InfoService(c fiber.Ctx) error {
	var err error
	handlers.MapMicroService, err = handlers.LoadConfig(config.BasePath + config.ConfigPath)
	if err != nil {
		return err
	}

	for i := range handlers.MapMicroService.Services {
		// USING NET/HTTP
		req, err := http.NewRequest("GET", handlers.MapMicroService.Services[i].Url, nil)
		if err != nil {
			handlers.MapMicroService.Services[i].Status = false
			continue
		}

		resp, err := proxyroute.ProxyClient.Do(req)
		if err != nil {
			handlers.MapMicroService.Services[i].Status = false
			continue
		}
		defer resp.Body.Close()

		handlers.MapMicroService.Services[i].Status = resp.StatusCode == http.StatusOK

		// USING FASTHTTP
		// req := fasthttp.AcquireRequest()
		// res := fasthttp.AcquireResponse()
		// defer fasthttp.ReleaseRequest(req)
		// defer fasthttp.ReleaseResponse(res)

		// req.SetRequestURI(handlers.MapMicroService.Services[i].Url)

		// handlers.MapMicroService.Services[i].Status = true

		// if err := proxyroute.ProxyClient.Do(req, res); err != nil {
		// 	handlers.MapMicroService.Services[i].Status = false
		// }

		// if res.StatusCode() != fiber.StatusOK {
		// 	handlers.MapMicroService.Services[i].Status = false
		// }
	}

	return c.JSON(handlers.MapMicroService.Services)
}

func CheckLocalService(c fiber.Ctx) error {
	url := c.Query("url")
	if url == "" {
		return handlers.UnprocessableEntityErrorResponse(c, fmt.Errorf("need url params"))
	}

	getResponse := fiber.Query[bool](c, "getRes")

	resp, err := http.Get(url)
	if err != nil {
		return handlers.SuccessResponse(c, true, "url is not active", false, nil)
	}
	defer resp.Body.Close()

	if getResponse {
		body, _ := io.ReadAll(resp.Body)
		c.Set("Content-Type", "application/json")
		c.Status(resp.StatusCode)
		return c.Send(body)
	}

	return handlers.SuccessResponse(c, true, "url is active", true, nil)
}

func RestartHandler(c fiber.Ctx) error {
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

func GetStatsLogger(c fiber.Ctx) error {
	layout := "2006-01-02T15:04:05.000Z07:00"

	isDetail := fiber.Query[bool](c, "detail")

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
		d := &[]models.Logger{}
		service := c.Query("service")
		method := c.Query("method")
		path := c.Query("path")
		status := c.Query("status")

		err := models.FindLogger(d, service, method, path, status, from, to).Error
		if err != nil {
			return handlers.InternalServerErrorResponse(c, err)
		}

		return handlers.SuccessResponse(c, true, "success to get detail log proxy", d, nil)
	} else {
		d := &[]models.PathStats{}

		err := models.FindStatsLogger(d, from, to).Error
		if err != nil {
			return handlers.InternalServerErrorResponse(c, err)
		}

		return handlers.SuccessResponse(c, true, "success to get stats log proxy", d, nil)
	}
}

func ProxyLocalService(c fiber.Ctx) error {
	urlQuery := c.Query("url")
	if urlQuery == "" {
		return handlers.UnprocessableEntityErrorResponse(c, fmt.Errorf("need url params"))
	}

	// target, _ := url.Parse(urlQuery)
	// proxy := httputil.NewSingleHostReverseProxy(target)

	// proxy.WithClient(proxyroute.ProxyClient)

	if err := proxy.Do(c, urlQuery); err != nil {
		return err
	}

	c.Response().Header.ContentType()
	// c.Response().Header.Del(fiber.HeaderServer)
	return nil
}
