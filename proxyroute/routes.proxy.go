package proxyroute

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"go-gerbang/config"
	"go-gerbang/handlers"
	"go-gerbang/middleware"
	"go-gerbang/models"
	"go-gerbang/types"

	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/gofiber/contrib/casbin"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func MainProxyRoutes(app *fiber.App) {
	var err error

	cfg, err := handlers.LoadConfig(config.BasePath + config.ConfigPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	handlers.MapMicroServiceMutex.Lock()
	handlers.MapMicroService = cfg
	handlers.MapMicroServiceMutex.Unlock()

	go handlers.WatchConfigFile(config.BasePath + config.ConfigPath)

	RegisterRoutes(app)
}

var ProxyClient = &fasthttp.Client{
	NoDefaultUserAgentHeader: true,
	DisablePathNormalizing:   true,
	MaxConnsPerHost:          10000,
	MaxIdleConnDuration:      90 * time.Second,
	ReadTimeout:              30 * time.Second,
	WriteTimeout:             30 * time.Second,
	MaxConnWaitTimeout:       10 * time.Second,
}

func RegisterRoutes(app *fiber.App) {
	handlers.MapMicroServiceMutex.RLock()
	defer handlers.MapMicroServiceMutex.RUnlock()

	proxy.WithClient(ProxyClient)

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

	sort.Slice(handlers.MapMicroService.Services, func(i, j int) bool {
		return len(handlers.MapMicroService.Services[i].Path) > len(handlers.MapMicroService.Services[j].Path)
	})

	for _, service := range handlers.MapMicroService.Services {
		// if service.AuthProtection && service.CsrfProtection && service.SessionProtection && service.RbacProtection {
		// 	app.Use(service.Path, middleware.CsrfProtection, middleware.Auth, middleware.ValidateSession, authz.RoutePermission(), proxyHandler(service))
		// } else if service.AuthProtection && service.CsrfProtection && service.RbacProtection {
		// 	app.Use(service.Path, middleware.CsrfProtection, middleware.Auth, authz.RoutePermission(), proxyHandler(service))
		// } else if service.AuthProtection && service.CsrfProtection && service.SessionProtection {
		// 	app.Use(service.Path, middleware.CsrfProtection, middleware.Auth, middleware.ValidateSession, proxyHandler(service))
		// } else if service.AuthProtection && service.CsrfProtection {
		// 	app.Use(service.Path, middleware.CsrfProtection, middleware.Auth, proxyHandler(service))
		// } else if service.AuthProtection {
		// 	app.Use(service.Path, middleware.Auth, proxyHandler(service))
		// } else {
		// 	app.Use(service.Path, proxyHandler(service))
		// }

		middlewares := []fiber.Handler{}

		if service.CsrfProtection {
			middlewares = append(middlewares, middleware.CsrfProtection)
		}
		if service.AuthProtection {
			middlewares = append(middlewares, middleware.Auth)
		}
		if service.SessionProtection {
			middlewares = append(middlewares, middleware.ValidateSession)
		}
		if service.RbacProtection {
			middlewares = append(middlewares, authz.RoutePermission())
		}

		// Build args properly
		if len(middlewares) > 0 {
			args := []interface{}{service.Path + "*"}
			for _, m := range middlewares {
				args = append(args, m)
			}
			app.Use(args...)
		}

		// Always add proxy handler
		app.All(service.Path+"*", proxyHandler(service))
	}
}

func proxyHandler(service types.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		requestBody := string(c.Body())
		method := c.Method()
		path := c.OriginalURL()

		prefixLen := len(service.Path)
		url := service.Url + c.OriginalURL()[prefixLen:]
		err := proxy.DoDeadline(c, url, time.Now().Add(10*time.Second))

		responseBody := string(c.Response().Body())
		status := c.Response().StatusCode()
		duration := time.Since(start)

		userField := zap.Skip()
		if user, ok := c.Locals("user").(*models.UserData); ok {
			fmt.Println(user)
			userField = zap.String("user", user.Username)
		}

		if err != nil {
			if status == 0 || status == fiber.StatusOK {
				status = fiber.StatusBadGateway
			}
			handlers.ZapLogger.Error(service.Service,
				zap.String("method", method),
				zap.String("path", path),
				zap.Int("status", status),
				zap.Duration("duration", duration),
				zap.Error(err),
				userField,
			)
			return handlers.InternalServerErrorResponse(c, fmt.Errorf("upstream unavailable: %v", err))
		}

		contentType := string(c.Response().Header.ContentType())

		var respLog zap.Field
		if strings.Contains(contentType, "application/json") || strings.HasPrefix(contentType, "text/") {
			maxLen := 1024 // 1KB
			respStr := string(responseBody)
			if len(respStr) > maxLen {
				respStr = respStr[:maxLen] + "...(truncated)"
			}
			respLog = zap.String("response", respStr)
		} else {
			respLog = zap.String("response_skipped", "binary or large response")
		}

		handlers.ZapLogger.Info(service.Service,
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Duration("duration", duration),
			zap.String("request", requestBody),
			respLog,
			zap.Int("response_size", len(responseBody)),
			zap.String("content_type", contentType),
			userField,
		)

		return nil
	}
}
