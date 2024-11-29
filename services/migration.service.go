package services

import (
	"go-gerbang/handlers"
	"go-gerbang/models"

	"go-gerbang/database"

	"github.com/gofiber/fiber/v2"
)

func CheckMigrationStatus(c *fiber.Ctx) error {
	_, err := models.FindUserByIdentity("admin")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Admin not found"})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Admin found"})
}

func MigrationService(c *fiber.Ctx) error {
	err := database.GDB.AutoMigrate(&models.User{}, &models.UserLogIp{}, &models.AuthRule{}, &models.UserAssignment{})

	if err != nil {
		return c.JSON(fiber.Map{"message": "failed to migration", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "success migration"})
	// database.GDB.AutoMigrate(&models.UserLogIp{})
	// database.GDB.AutoMigrate(&models.AuthRule{})
	// database.GDB.AutoMigrate(&models.UserAssignment{})
	// database.GDB.Migrator().CreateTable(&models.User{})
	// return nil
}

func MigrateAdminUser(c *fiber.Ctx) error {
	if database.GDB.Migrator().HasTable(&models.User{}) {
		errCreate := database.GDB.Create(&models.User{
			Username:      "admin",
			FullName:      "Admin",
			StatusAccount: 10,
			PasswordHash:  handlers.GeneratePasswordHash("@dmin9192"),
		})

		if errCreate != nil {
			return c.JSON(fiber.Map{"message": "failed to create admin", "error": errCreate})
		}
	}

	return c.JSON(fiber.Map{"message": "success create admin"})
}
