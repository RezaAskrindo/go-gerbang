package services

import (
	"fmt"

	"go-gerbang/handlers"
	"go-gerbang/models"

	"github.com/gofiber/fiber/v2"
)

func FindUserById(c *fiber.Ctx) error {
	userId := c.Params("userId")

	if userId == "" {
		return handlers.UnprocessableEntityErrorResponse(c, fmt.Errorf("need userId params"))
	}

	user := new(models.User)
	if err := models.FindUserById(user, userId); err != nil {
		return handlers.NotFoundErrorResponse(c, err)
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
		return handlers.UnprocessableEntityErrorResponse(c, fmt.Errorf("need userId params"))
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

func DeleteUser(c *fiber.Ctx) error {
	userId := c.Params("userId")

	if userId == "" {
		return handlers.UnprocessableEntityErrorResponse(c, fmt.Errorf("need userId params"))
	}

	if err := models.HardDeleteUser(userId); err != nil {
		return handlers.NotFoundErrorResponse(c, err.Error)
	}

	return handlers.SuccessResponse(c, true, "success to delete user", nil, nil)
}
