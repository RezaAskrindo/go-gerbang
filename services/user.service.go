package services

import (
	"fmt"

	"go-gerbang/handlers"
	"go-gerbang/models"

	"github.com/gofiber/fiber/v2"
)

func GetAllUser(c *fiber.Ctx) error {
	d := &[]models.UserData{}

	err := models.FindAllUser(d).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Find All user Is Error"})
	}

	return c.JSON(&models.UserDataResponses{
		Items: d,
	})
}

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

	return handlers.SuccessResponse(c, true, "success to create user", u, nil)
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
		return handlers.NotFoundErrorResponse(c, err)
	}

	return handlers.SuccessResponse(c, true, "success to delete user", nil, nil)
}

// USER ASSIGNMENT ROLE
func GetAllUserAssignments(c *fiber.Ctx) error {
	var count int64
	assignments := &[]models.UserAssignment{}

	err := models.FindUserAssignment(assignments).Error
	if err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	err = models.CountFindUserAssignment(&count)
	if err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	return handlers.SuccessResponse(c, true, "success to get all user assignments", assignments, &count)
}

func CreateUserAssignment(c *fiber.Ctx) error {
	u := new(models.UserAssignment)

	if err := handlers.ParseBody(c, u); err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	if err := handlers.ValidateStruct(*u); err != nil {
		return handlers.SuccessResponse(c, false, "error validation user assignment", err, nil)
	}

	if err := models.CreateUserAssignment(u).Error; err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	return handlers.SuccessResponse(c, true, "success to create user assignment", nil, nil)
}

func CreateUserAssignmentsBulk(c *fiber.Ctx) error {
	assignments := new([]models.UserAssignment)

	if err := handlers.ParseBody(c, assignments); err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	for _, a := range *assignments {
		if err := handlers.ValidateStruct(a); err != nil {
			return handlers.SuccessResponse(c, false, "validation failed on bulk create", err, nil)
		}
	}

	if err := models.CreateUserAssignments(*assignments).Error; err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	return handlers.SuccessResponse(c, true, "success to bulk create user assignments", nil, nil)
}

func UpdateUserAssignment(c *fiber.Ctx) error {
	u := new(models.UserAssignment)

	if err := handlers.ParseBody(c, u); err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	if err := handlers.ValidateStruct(*u); err != nil {
		return handlers.SuccessResponse(c, false, "error validation user assignment", err, nil)
	}

	if err := models.UpdateUserAssignment(u).Error; err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	return handlers.SuccessResponse(c, true, "success to update user assignment", nil, nil)
}

func UpdateUserAssignmentsBulk(c *fiber.Ctx) error {
	assignments := new([]models.UserAssignment)

	if err := handlers.ParseBody(c, assignments); err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	for _, a := range *assignments {
		if err := handlers.ValidateStruct(a); err != nil {
			return handlers.SuccessResponse(c, false, "validation failed on bulk update", err, nil)
		}
	}

	if err := models.UpdateUserAssignments(*assignments); err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	return handlers.SuccessResponse(c, true, "success to bulk update user assignments", nil, nil)
}

func DeleteUserAssignment(c *fiber.Ctx) error {
	accountId := c.Params("account_id")
	authRoleId, _ := c.ParamsInt("auth_role_id")

	if accountId == "" || authRoleId == 0 {
		return handlers.UnprocessableEntityErrorResponse(c, fmt.Errorf("need account_id and auth_role_id params"))
	}

	if err := models.DeleteUserAssignment(accountId, authRoleId).Error; err != nil {
		return handlers.NotFoundErrorResponse(c, err)
	}

	return handlers.SuccessResponse(c, true, "success to delete user assignment", nil, nil)
}
