package services

import (
	"fmt"
	"go-gerbang/handlers"
	"go-gerbang/models"

	"github.com/gofiber/fiber/v2"
)

func GetAllAuthRole(c *fiber.Ctx) error {
	var countAllRole int64

	d := &[]models.AuthRule{}

	err := models.FindAuthRule(d).Error
	if err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	err = models.CountFindAuthRule(&countAllRole)
	if err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	return handlers.SuccessResponse(c, true, "success to get all auth role", d, &countAllRole)
}

func CreateAuthRole(c *fiber.Ctx) error {
	u := new(models.AuthRule)

	if err := handlers.ParseBody(c, u); err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	if err := handlers.ValidateStruct(*u); err != nil {
		return handlers.SuccessResponse(c, false, "error validation auth role", err, nil)
	}

	if err := models.CreateAuthRule(u).Error; err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	return handlers.SuccessResponse(c, true, "success to create auth role", nil, nil)
}

func UpdateAuthRole(c *fiber.Ctx) error {
	authRoleId := c.Params("id")

	if authRoleId == "" {
		return handlers.UnprocessableEntityErrorResponse(c, fmt.Errorf("need id params"))
	}

	u := new(models.AuthRule)

	if err := handlers.ParseBody(c, u); err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	if err := handlers.ValidateStruct(*u); err != nil {
		return handlers.SuccessResponse(c, false, "error validation auth role", err, nil)
	}

	if err := models.UpdateAuthRule(authRoleId, u).Error; err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	return handlers.SuccessResponse(c, true, "success to update auth role", nil, nil)
}

func DeleteAuthRule(c *fiber.Ctx) error {
	authRoleId := c.Params("id")

	if authRoleId == "" {
		return handlers.UnprocessableEntityErrorResponse(c, fmt.Errorf("need id params"))
	}

	if err := models.DeleteAuthRule(authRoleId); err != nil {
		return handlers.NotFoundErrorResponse(c, err.Error)
	}

	return handlers.SuccessResponse(c, true, "success to delete auth role", nil, nil)
}
