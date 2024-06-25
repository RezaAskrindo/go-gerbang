package services

import (
	"encoding/json"
	"os"
	"sika_apigateway/types"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func IndexService(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}

func ProtectService(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"title": "Testing Protect Route"})
}

func InfoService(c *fiber.Ctx) error {
	file, err := os.Open("./config/config.json")
	if err != nil {
		return err
	}
	defer file.Close()

	var MapMicroService []types.ValueMicroService
	err = json.NewDecoder(file).Decode(&MapMicroService)
	if err != nil {
		return err
	}

	for i := range MapMicroService {
		req := fasthttp.AcquireRequest()
		res := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(res)

		req.SetRequestURI(MapMicroService[i].Url)

		MapMicroService[i].Status = true

		client := &fasthttp.Client{}
		if err := client.Do(req, res); err != nil {
			// fmt.Printf("Error: %s\n", err)
			MapMicroService[i].Status = false
		}

		if res.StatusCode() != fiber.StatusOK {
			MapMicroService[i].Status = false
		}
	}

	return c.JSON(MapMicroService)
}
