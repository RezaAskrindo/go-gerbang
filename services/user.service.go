package services

import (
	"errors"

	// "go-gerbang/config"
	"go-gerbang/handlers"
	"go-gerbang/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func FindUserById(c *fiber.Ctx) error {
	userId := c.Params("userId")

	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(&handlers.ErrorStruct{
			Message: "Need Params userId",
			Status:  false,
		})
	}

	user := new(models.User)

	err := models.FindUserByIdRaw(user, userId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(&handlers.ErrorStruct{
			Message: "User Not Found",
			Status:  false,
		})
	}

	return handlers.SuccessResponse(c, true, "success to get users", user, nil)
}

func CreateUser(c *fiber.Ctx) error {
	u := new(models.User)

	if err := handlers.ParseBody(c, u); err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	if err := handlers.ValidateStruct(*u); err != nil {
		return handlers.SuccessResponse(c, false, "error validation user", err, nil)
	}

	if err := models.CreateUser(u).Error; err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	return handlers.SuccessResponse(c, true, "success to create user", nil, nil)
}

func UpdateUser(c *fiber.Ctx) error {
	userId := c.Params("userId")

	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(&handlers.ErrorStruct{
			Message: "Need Params userId",
			Status:  false,
		})
	}

	u := new(models.User)

	if err := handlers.ParseBody(c, u); err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	if err := handlers.ValidateStruct(*u); err != nil {
		return handlers.SuccessResponse(c, false, "error validation user", err, nil)
	}

	if err := models.UpdateUser(userId, u).Error; err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	return handlers.SuccessResponse(c, true, "success to update user", nil, nil)
}
