package services

import "github.com/gofiber/fiber/v2"

type SuccessStruct struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Total   *int64      `json:"total"`
}

type ErrorStruct struct {
	Message interface{} `json:"message"`
	Status  bool        `json:"status"`
	Code    int         `json:"code"`
}

func ParseBody(c *fiber.Ctx, body interface{}) error {
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&ErrorStruct{
			Message: "Failed To Parse Body",
			Status:  false,
			Code:    fiber.StatusBadRequest,
		})
	}

	return nil
}

func SuccessResponse(c *fiber.Ctx, message string, data interface{}, total *int64) error {
	return c.Status(fiber.StatusOK).JSON(&SuccessStruct{
		Status:  true,
		Message: message,
		Data:    data,
		Total:   total,
	})
}

func BadRequestErrorResponse(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusBadRequest).JSON(&ErrorStruct{
		Message: err.Error(),
		Status:  false,
		Code:    fiber.StatusBadRequest,
	})
}

func ConflictErrorResponse(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusConflict).JSON(&ErrorStruct{
		Message: err.Error(),
		Status:  false,
		Code:    fiber.StatusConflict,
	})
}

func InternalServerErrorResponse(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusInternalServerError).JSON(&ErrorStruct{
		Message: err.Error(),
		Status:  false,
		Code:    fiber.StatusInternalServerError,
	})
}

func UnauthorizedErrorResponse(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(&ErrorStruct{
		Message: err.Error(),
		Status:  false,
		Code:    fiber.StatusUnauthorized,
	})
}
