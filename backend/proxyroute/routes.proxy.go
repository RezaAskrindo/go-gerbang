package proxyroute

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"go-gerbang/config"
	"go-gerbang/handlers"
	"go-gerbang/middleware"
	"go-gerbang/models"
	"go-gerbang/types"

	// fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
	// "github.com/gofiber/fiber/v3/middleware/proxy"
	// "github.com/valyala/fasthttp"
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

// var ProxyClient = &fasthttp.Client{
// 	NoDefaultUserAgentHeader: true,
// 	DisablePathNormalizing:   true,
// 	MaxConnsPerHost:          10000,
// 	MaxIdleConnDuration:      90 * time.Second,
// 	ReadTimeout:              30 * time.Second,
// 	WriteTimeout:             30 * time.Second,
// 	MaxConnWaitTimeout:       10 * time.Second,
// }

var ProxyClient = &http.Client{
	Timeout: 5 * time.Minute,
	Transport: &http.Transport{
		MaxConnsPerHost:     10000,
		IdleConnTimeout:     90 * time.Second,
		DisableKeepAlives:   false,
		DisableCompression:  false,
		MaxIdleConns:        10000,
		MaxIdleConnsPerHost: 10000,
	},
}

type CasbinRule struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Ptype string `gorm:"size:255;uniqueIndex:unique_index"`
	V0    string `gorm:"size:255;uniqueIndex:unique_index"`
	V1    string `gorm:"size:255;uniqueIndex:unique_index"`
	V2    string `gorm:"size:128;uniqueIndex:unique_index"`
	V3    string `gorm:"default:null;size:128;uniqueIndex:unique_index"`
	V4    string `gorm:"default:null;size:128;uniqueIndex:unique_index"`
	V5    string `gorm:"default:null;size:128;uniqueIndex:unique_index"`
}

func RegisterRoutes(app *fiber.App) {
	handlers.MapMicroServiceMutex.RLock()
	defer handlers.MapMicroServiceMutex.RUnlock()

	// proxy.WithClient(ProxyClient)

	// RBAC PROTECTION
	// NEED AUTH PROTECTION
	// a, _ := gormadapter.NewAdapterByDBWithCustomTable(database.GDB, &CasbinRule{})
	// a := fileadapter.NewAdapter(config.BasePath + config.Config("CONFIG_PATH_CASBIN_POLICY"))
	// authz := casbin.New(casbin.Config{
	// 	ModelFilePath: config.BasePath + config.Config("CONFIG_PATH_CASBIN_MODEL"),
	// 	PolicyAdapter: a,
	// 	Lookup: func(c fiber.Ctx) string {
	// 		user, ok := c.Locals("user").(*models.UserData)
	// 		if !ok {
	// 			return ""
	// 		}
	// 		sub := strconv.Itoa(int(user.StatusAccount))

	// 		return sub
	// 	},
	// 	Unauthorized: func(c fiber.Ctx) error {
	// 		return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("your role don't have authorization"))
	// 	},
	// 	Forbidden: func(c fiber.Ctx) error {
	// 		// user, _ := c.Locals("user").(*models.UserData)
	// 		// fmt.Printf("[DEBUG] FORBIDDEN: sub=%v, obj=%v, act=%v\n",
	// 		// 	user.StatusAccount,
	// 		// 	c.Path(),
	// 		// 	c.Method(),
	// 		// )
	// 		return handlers.ForbiddenErrorResponse(c, fmt.Errorf("your role don't have access"))
	// 	},
	// })

	sort.Slice(handlers.MapMicroService.Services, func(i, j int) bool {
		return len(handlers.MapMicroService.Services[i].Path) > len(handlers.MapMicroService.Services[j].Path)
	})

	for _, service := range handlers.MapMicroService.Services {
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
			// middlewares = append(middlewares, middleware.Auth)
			middlewares = append(middlewares, middleware.AuthRBAC)
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

// USING FASTHTTP
// func proxyHandler(service types.Service) fiber.Handler {
// 	return func(c fiber.Ctx) error {
// 		start := time.Now()

// 		requestBody := string(c.Body())
// 		method := c.Method()
// 		path := c.OriginalURL()

// 		prefixLen := len(service.Path)
// 		url := service.Url + c.OriginalURL()[prefixLen:]
// 		// err := proxy.DoDeadline(c, url, time.Now().Add(10*time.Second))
// 		err := proxy.DoDeadline(c, url, time.Now().Add(5*time.Minute))

// 		responseBody := string(c.Response().Body())
// 		status := c.Response().StatusCode()
// 		duration := time.Since(start)

// 		userField := zap.Skip()
// 		if user, ok := c.Locals("user").(*models.UserData); ok {
// 			userField = zap.String("user", user.Username)
// 		}

// 		if err != nil {
// 			if status == 0 || status == fiber.StatusOK {
// 				status = fiber.StatusBadGateway
// 			}
// 			handlers.ZapLogger.Error(service.Service,
// 				zap.String("method", method),
// 				zap.String("path", path),
// 				zap.Int("status", status),
// 				zap.Duration("duration", duration),
// 				zap.Error(err),
// 				userField,
// 			)
// 			return handlers.InternalServerErrorResponse(c, fmt.Errorf("upstream unavailable: %v", err))
// 		}

// 		contentType := string(c.Response().Header.ContentType())

// 		var respLog zap.Field
// 		if strings.Contains(contentType, "application/json") || strings.HasPrefix(contentType, "text/") {
// 			maxLen := 1024 // 1KB
// 			respStr := string(responseBody)
// 			if len(respStr) > maxLen {
// 				respStr = respStr[:maxLen] + "...(truncated)"
// 			}
// 			respLog = zap.String("response", respStr)
// 		} else {
// 			respLog = zap.String("response_skipped", "binary or large response")
// 		}

// 		handlers.ZapLogger.Info(service.Service,
// 			zap.String("method", method),
// 			zap.String("path", path),
// 			zap.Int("status", status),
// 			zap.Duration("duration", duration),
// 			zap.String("request", requestBody),
// 			respLog,
// 			zap.Int("response_size", len(responseBody)),
// 			zap.String("content_type", contentType),
// 			userField,
// 		)

// 		return nil
// 	}
// }

// USING NET HTTP
func proxyHandler(service types.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		requestBody := string(c.Body())
		method := c.Method()
		path := c.OriginalURL()

		prefixLen := len(service.Path)
		upstreamURL := service.Url + c.OriginalURL()[prefixLen:]

		// Create HTTP request
		req, err := http.NewRequest(method, upstreamURL, strings.NewReader(requestBody))
		if err != nil {
			return handleProxyError(c, service, method, path, err)
		}

		// Copy headers
		for header, values := range c.GetReqHeaders() {
			for _, value := range values {
				req.Header.Add(header, value)
			}
		}

		// Execute request (automatically follows redirects)
		resp, err := ProxyClient.Do(req)
		if err != nil {
			return handleProxyError(c, service, method, path, err)
		}
		defer resp.Body.Close()

		// Read response body
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return handleProxyError(c, service, method, path, err)
		}

		// Copy response headers
		for key, values := range resp.Header {
			for _, value := range values {
				c.Append(key, value)
			}
		}

		// Set status and response
		c.Status(resp.StatusCode)
		duration := time.Since(start)

		// Logging
		logProxyRequest(service, method, path, resp.StatusCode, duration, requestBody, responseBody, c)

		return c.Send(responseBody)
	}
}

func handleProxyError(c fiber.Ctx, service types.Service, method, path string, err error) error {
	duration := time.Since(time.Time{}) // You might want to track this properly

	userField := zap.Skip()
	if user, ok := c.Locals("user").(*models.UserData); ok {
		userField = zap.String("user", user.Username)
	}

	handlers.ZapLogger.Error(service.Service,
		zap.String("method", method),
		zap.String("path", path),
		zap.Int("status", fiber.StatusBadGateway),
		zap.Duration("duration", duration),
		zap.Error(err),
		userField,
	)

	return handlers.InternalServerErrorResponse(c, fmt.Errorf("upstream unavailable: %v", err))
}

func logProxyRequest(service types.Service, method, path string, status int, duration time.Duration, reqBody string, respBody []byte, c fiber.Ctx) {
	userField := zap.Skip()
	if user, ok := c.Locals("user").(*models.UserData); ok {
		userField = zap.String("user", user.Username)
	}

	contentType := string(c.Response().Header.ContentType())

	var respLog zap.Field
	if strings.Contains(contentType, "application/json") || strings.HasPrefix(contentType, "text/") {
		maxLen := 1024
		respStr := string(respBody)
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
		zap.String("request", reqBody),
		respLog,
		zap.Int("response_size", len(respBody)),
		zap.String("content_type", contentType),
		userField,
	)
}
