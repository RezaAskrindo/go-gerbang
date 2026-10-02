package services

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"

	"go-gerbang/broker"
	"go-gerbang/config"
	"go-gerbang/handlers"
	"go-gerbang/models"
	"go-gerbang/proxyroute"
	"go-gerbang/types"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/nats-io/nats.go"
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

	isGroup := fiber.Query[bool](c, "group", false)

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

		err := models.FindStatsLogger(d, from, to, isGroup).Error
		if err != nil {
			return handlers.InternalServerErrorResponse(c, err)
		}

		return handlers.SuccessResponse(c, true, "success to get stats log proxy", d, nil)
	}
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

type proxyRule struct {
	Methods map[string]bool
	Paths   map[string]bool // exact-match paths; nil/empty = any path
}

var allowedProxyTargets = map[string]proxyRule{
	"localhost:2019": {
		Methods: map[string]bool{
			"GET":  true,
			"POST": true,
		},
		Paths: map[string]bool{
			"/config/": true,
			"/load":    true, // replace config wholesale — the only mutating path allowed
		},
	},
}

func isProxyTargetAllowed(host, method, path string) bool {
	rule, ok := allowedProxyTargets[host]
	if !ok {
		return false
	}
	if !rule.Methods[strings.ToUpper(method)] {
		return false
	}
	if len(rule.Paths) > 0 && !rule.Paths[path] {
		return false
	}
	return true
}

func ProxyLocalService(c fiber.Ctx) error {
	urlQuery := c.Query("url")
	if urlQuery == "" {
		return handlers.UnprocessableEntityErrorResponse(c, fmt.Errorf("need url params"))
	}

	parsed, err := url.Parse(urlQuery)
	if err != nil {
		return handlers.UnprocessableEntityErrorResponse(c, fmt.Errorf("invalid url"))
	}

	if !isProxyTargetAllowed(parsed.Host, c.Method(), parsed.Path) {
		return handlers.UnprocessableEntityErrorResponse(c, fmt.Errorf("target not allowed"))
	}

	body := bytes.NewReader(c.Request().Body())

	req, err := http.NewRequest(c.Method(), urlQuery, body)
	if err != nil {
		return handlers.SuccessResponse(c, true, "invalid request", false, nil)
	}

	headersToKeep := []string{
		"Content-Type",
		"Accept",
		"Content-Length",
		"Authorization",
	}

	for _, headerName := range headersToKeep {
		if value := c.Request().Header.Peek(headerName); value != nil {
			req.Header.Set(headerName, string(value))
		}
	}

	req.Header.Set("User-Agent", "GO GERBANG CLIENT")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return handlers.SuccessResponse(c, true, "url is not active", false, nil)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return handlers.SuccessResponse(c, true, "error reading response", false, nil)
	}

	for key, values := range resp.Header {
		for _, value := range values {
			c.Response().Header.Add(key, value)
		}
	}
	c.Status(resp.StatusCode)

	if len(respBody) == 0 {
		return handlers.SuccessResponse(c, true, "response body is empty", false, nil)
	}

	return c.Send(respBody)
}

func JetStreamMetricsHandler(c fiber.Ctx) error {
	js, err := broker.NatsClient.JetStream()
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error":   "JetStream unavailable",
			"details": err.Error(),
		})
	}

	metrics, err := GetDetailedJetStreamMetrics(js)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve metrics",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(metrics)
}

func GetDetailedJetStreamMetrics(js nats.JetStreamContext) (*types.DetailedMetrics, error) {
	ai, err := js.AccountInfo()
	if err != nil {
		return nil, err
	}

	m := &types.DetailedMetrics{
		MemoryUsed:     ai.Memory,
		StoreUsed:      ai.Store,
		MaxMemory:      ai.Limits.MaxMemory,
		MaxStore:       ai.Limits.MaxStore,
		TotalStreams:   ai.Streams,
		TotalConsumers: ai.Consumers,
		Streams:        []types.StreamMetrics{},
	}

	for si := range js.StreamsInfo() {
		sm := types.StreamMetrics{
			Name:      si.Config.Name,
			Subjects:  si.Config.Subjects,
			Messages:  si.State.Msgs,
			Bytes:     si.State.Bytes,
			FirstSeq:  si.State.FirstSeq,
			LastSeq:   si.State.LastSeq,
			Created:   si.Created,
			Consumers: []types.ConsumerMetrics{},
		}

		for ci := range js.ConsumersInfo(si.Config.Name) {
			sm.Consumers = append(sm.Consumers, types.ConsumerMetrics{
				Name:         ci.Name,
				Pending:      ci.NumPending,
				AckPending:   ci.NumAckPending,
				Redelivered:  ci.NumRedelivered,
				Waiting:      ci.NumWaiting,
				DeliveredSeq: ci.Delivered.Stream,
				AckFloorSeq:  ci.AckFloor.Stream,
			})
		}

		sm.Schedules = getScheduleMetrics(js, si)

		m.Streams = append(m.Streams, sm)
	}

	return m, nil
}

const maxScheduleScan = 200

func getScheduleMetrics(js nats.JetStreamContext, si *nats.StreamInfo) *types.ScheduleMetrics {
	out := &types.ScheduleMetrics{Items: []types.ScheduledMessage{}}
	if si.State.Msgs == 0 {
		return out
	}

	full, err := js.StreamInfo(si.Config.Name, &nats.StreamInfoRequest{SubjectsFilter: ">"})
	if err != nil {
		return out
	}

	scanned := 0
	for subj := range full.State.Subjects {
		if scanned >= maxScheduleScan {
			out.Truncated = true
			break
		}
		scanned++

		raw, err := js.GetLastMsg(si.Config.Name, subj)
		if err != nil {
			continue
		}
		sched := raw.Header.Get("Nats-Schedule")
		if sched == "" {
			continue
		}

		out.Items = append(out.Items, types.ScheduledMessage{
			Subject:  subj,
			Seq:      raw.Sequence,
			Schedule: sched,
			Target:   raw.Header.Get("Nats-Schedule-Target"),
			Source:   raw.Header.Get("Nats-Schedule-Source"),
			TTL:      raw.Header.Get("Nats-Schedule-TTL"),
			TimeZone: raw.Header.Get("Nats-Schedule-Time-Zone"),
			Stored:   raw.Time,
		})
	}
	out.Count = len(out.Items)
	return out
}
