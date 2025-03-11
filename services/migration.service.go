package services

import (
	"go-gerbang/handlers"
	"go-gerbang/models"
	"log"

	"go-gerbang/database"

	"github.com/gofiber/fiber/v2"
)

func CheckMigrationStatus(c *fiber.Ctx) error {
	user := new(models.User)
	err := models.FindUserByIdentity(user, "admin", "admin", "admin", "admin")
	if err != nil {
		return handlers.NotFoundErrorResponse(c, err)
	}
	return c.JSON(fiber.Map{"success": true, "message": "Admin found"})
}

func MigrationService(c *fiber.Ctx) error {
	err := database.GDB.AutoMigrate(&models.User{}, &models.UserLogIp{}, &models.AuthRule{}, &models.UserAssignment{})

	if err != nil {
		return c.JSON(fiber.Map{"message": "failed to migration", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "success migration"})
}

func MigrateAdminUser(c *fiber.Ctx) error {
	log.Println("here")
	if database.GDB.Migrator().HasTable(&models.User{}) {
		errCreate := database.GDB.Create(&models.User{
			Username:      "admin",
			FullName:      "Admin",
			StatusAccount: 10,
			PasswordHash:  handlers.GeneratePasswordHash("@dmin9192"),
		})

		if errCreate.Error != nil {
			// IF ERROR CACHE, RUN THIS IN SQL:
			// DISCARD ALL;
			return c.JSON(fiber.Map{"message": "failed to create admin", "error": errCreate.Error})
		}
		return c.JSON(fiber.Map{"message": "success create admin"})
	} else {
		return c.JSON(fiber.Map{"message": "table user is not found"})
	}
}
