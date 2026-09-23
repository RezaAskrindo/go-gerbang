package services

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"go-gerbang/handlers"
	"go-gerbang/models"

	"github.com/gofiber/fiber/v3"
)

func GetConfigurationByGroup(c fiber.Ctx) error {
	group := c.Params("group")

	if group == "" {
		return handlers.UnprocessableEntityErrorResponse(c, fmt.Errorf("need group params"))
	}

	d := &[]models.Configuration{}

	useId := false

	config_name := c.Query("config_name")

	if config_name != "" {
		useId = true
		err := models.FindConfiguration(d, "configuration_group = ? AND configuration_name = ?", group, config_name).Error
		if err != nil {
			return handlers.InternalServerErrorResponse(c, err)
		}
	} else {
		err := models.FindConfiguration(d, "configuration_group = ?", group).Error
		if err != nil {
			return handlers.InternalServerErrorResponse(c, err)
		}
	}

	if len(*d) > 0 {
		result := models.ParseConfiguration(d, useId)
		var count *int64
		if arr, ok := result.([]map[string]string); ok {
			c := int64(len(arr))
			count = &c
		}

		if useId {
			parseConfig := fmt.Sprintf("%s-%s", group, config_name)
			_ = handlers.SaveToRedis(parseConfig, result)
		}

		return handlers.SuccessResponse(c, true, "success to get config", result, count)
	}

	return handlers.SuccessResponse(c, true, "config is empty", nil, nil)
}

func UpsertConfiguration(c fiber.Ctx) error {
	body := []models.Configuration{}

	if err := handlers.ParseBody(c, &body); err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	if err := models.CreateConfiguration(body).Error; err != nil {
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
	d := &[]models.GroupConfiguration{}
	err := models.FindGroupConfiguration(d, "configuration_group = ?", group).Error
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

type ExecuteScriptQuery struct {
	Name  string `query:"name"`
	Index *int   `query:"index"`
}

func ConfigExecuteScript(c fiber.Ctx) error {
	group := c.Params("group")

	if group == "" {
		return handlers.UnprocessableEntityErrorResponse(c, fmt.Errorf("need group params"))
	}

	q := new(ExecuteScriptQuery)
	if err := c.Bind().Query(q); err != nil {
		return handlers.UnprocessableEntityErrorResponse(c, err)
	}

	if q.Name == "" {
		return handlers.UnprocessableEntityErrorResponse(c, fmt.Errorf("missing name param"))
	}

	whereArgs := []interface{}{
		"configuration_group = ? AND configuration_name = ?",
		group, q.Name,
	}

	if q.Index != nil {
		whereArgs[0] = whereArgs[0].(string) + " AND configuration_index = ?"
		whereArgs = append(whereArgs, *q.Index)
	}

	rows := &[]models.Configuration{}
	if err := models.FindConfiguration(rows, whereArgs...).Error; err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	if len(*rows) == 0 {
		return handlers.UnprocessableEntityErrorResponse(c, fmt.Errorf("unknown script: %s", q.Name))
	}

	byIndex := make(map[int]map[string]string)
	for _, row := range *rows {
		idx := 0
		if row.ConfigurationIndex != nil {
			idx = *row.ConfigurationIndex
		}
		if byIndex[idx] == nil {
			byIndex[idx] = make(map[string]string)
		}
		if row.ConfigurationValue != nil {
			byIndex[idx][row.ConfigurationKey] = *row.ConfigurationValue
		}
	}

	if len(byIndex) > 1 {
		return handlers.UnprocessableEntityErrorResponse(c, fmt.Errorf("ambiguous script %q — pass index", q.Name))
	}

	var keys map[string]string
	for _, v := range byIndex {
		keys = v
	}

	workDirVal, fileVal := keys["location"], keys["desist"]
	if workDirVal == "" || fileVal == "" {
		return handlers.InternalServerErrorResponse(c, fmt.Errorf("script %q is not fully configured", q.Name))
	}

	baseFile := strings.TrimSuffix(fileVal, filepath.Ext(fileVal))
	base := filepath.Join(workDirVal, baseFile)

	scriptPath, workDir, err := handlers.ResolveScriptForOS(base)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error: script not found for this OS",
			"err":     err.Error(),
		})
	}

	go func() {
		if err := handlers.ExecuteScript(scriptPath, workDir); err != nil {
			log.Printf("script %q failed: %v", q.Name, err)
		}
	}()

	return c.JSON(fiber.Map{
		"message": "script is executing",
		"err":     nil,
	})
}
