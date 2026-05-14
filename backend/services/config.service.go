package services

import (
	"fmt"
	"os"

	"go-gerbang/handlers"
	"go-gerbang/models"

	"github.com/gofiber/fiber/v3"
)

func UpsertConfig(c fiber.Ctx) error {
	body := new(models.Config)

	if err := handlers.ParseBody(c, &body); err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	if err := models.CreatConfig(body).Error; err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	return handlers.SuccessResponse(c, true, "success to insert config", body, nil)
}

// CONFIGURATIONS
func GetConfigurationByGroup(c fiber.Ctx) error {
	group := c.Params("group")

	if group == "" {
		return handlers.UnprocessableEntityErrorResponse(c, fmt.Errorf("need group params"))
	}

	d := &[]models.Configurations{}

	config_name := c.Query("config_name")

	if config_name != "" {
		err := models.FindConfigurations(d, "configuration_group = ? AND configuration_name = ?", group, config_name).Error
		if err != nil {
			return handlers.InternalServerErrorResponse(c, err)
		}
	} else {
		err := models.FindConfigurations(d, "configuration_group = ?", group).Error
		if err != nil {
			return handlers.InternalServerErrorResponse(c, err)
		}
	}

	if len(*d) > 0 {
		result := models.ParseConfigurations(d)
		var count *int64
		if arr, ok := result.([]map[string]string); ok {
			c := int64(len(arr))
			count = &c
		}
		return handlers.SuccessResponse(c, true, "success to get config", result, count)
	}

	return handlers.SuccessResponse(c, true, "config is empty", nil, nil)
}

func UpsertConfiguration(c fiber.Ctx) error {
	body := []models.Configurations{}

	if err := handlers.ParseBody(c, &body); err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	if err := models.CreateConfigurations(body).Error; err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	return handlers.SuccessResponse(c, true, "success to insert config", body, nil)
}

func DeleteConfiguration(c fiber.Ctx) error {
	group := c.Params("group")
	if group == "" {
		return handlers.UnprocessableEntityErrorResponse(c, fmt.Errorf("need group params"))
	}

	config_name := c.Query("config_name")
	if config_name == "" {
		return handlers.UnprocessableEntityErrorResponse(c, fmt.Errorf("need config_name params"))
	}

	if err := models.DeleteConfigurationByConfName(group, config_name); err != nil {
		return handlers.NotFoundErrorResponse(c, err)
	}

	// RESET INDEX
	d := &[]models.GroupConfigurations{}
	err := models.FindGroupConfigurations(d, "configuration_group = ?", group).Error
	if err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	if len(*d) > 0 {
		for i, item := range *d {
			models.UpdateConfigurationIndex(item.ConfigurationGroup, item.ConfigurationName, i)
		}
	}

	return handlers.SuccessResponse(c, true, "success to delete config", nil, nil)
}

func ConfigExecuteScript(c fiber.Ctx) error {
	config_work_dir := c.Query("work_dir")
	config_file := c.Query("file")
	if config_work_dir == "" && config_file == "" {
		return handlers.UnprocessableEntityErrorResponse(c, fmt.Errorf("need config_url params"))
	}

	if _, err := os.Stat(config_work_dir + "\\" + config_file); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error: file not found",
			"err":     err.Error(),
		})
	}

	// Run in background goroutine (non-blocking)
	go func() {
		handlers.ExecuteScript(config_work_dir+"\\"+config_file, config_work_dir)
	}()

	return c.JSON(fiber.Map{
		"message": "script is execute",
		"err":     nil,
	})
}

// GENERIC OBJECT PARSE
// func CreateConfiguration(c fiber.Ctx) error {
// 	group := c.Params("group")

// 	if group == "" {
// 		return handlers.UnprocessableEntityErrorResponse(c, fmt.Errorf("need group params"))
// 	}

// 	config_name := c.Query("config_name")
// 	config_index := fiber.Query[int](c, "config_index")

// 	var body map[string]interface{}
// 	if err := handlers.ParseBody(c, &body); err != nil {
// 		return handlers.BadRequestErrorResponse(c, err)
// 	}

// 	data := []models.Configurations{}
// 	for k, v := range body {
// 		var valPtr *string
// 		if str, ok := v.(string); ok {
// 			valPtr = &str
// 		}
// 		data = append(data, models.Configurations{
// 			ConfigurationGroup: group,
// 			ConfigurationKey:   k,
// 			ConfigurationValue: valPtr,
// 			ConfigurationName:  &config_name,
// 			ConfigurationIndex: &config_index,
// 		})
// 	}

// 	count := int64(len(data))

// 	if err := models.CreateConfigurations(data).Error; err != nil {
// 		return handlers.InternalServerErrorResponse(c, err)
// 	}

// 	return handlers.SuccessResponse(c, true, "success to insert config "+group, data, &count)
// }
