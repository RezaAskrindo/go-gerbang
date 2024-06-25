package services

import (
	"sika_apigateway/handlers"
	"sika_apigateway/models"

	"sika_apigateway/database"

	"github.com/gofiber/fiber/v2"
)

func MigrationService(c *fiber.Ctx) error {
	database.GDB.AutoMigrate(&models.User{})
	if database.GDB.Migrator().HasTable(&models.User{}) {
		database.GDB.Create(&models.User{
			Username:      "admin",
			FullName:      "Admin",
			StatusAccount: 10,
			PasswordHash:  handlers.GeneratePasswordHash("@dmin9192"),
		})
	}
	database.GDB.AutoMigrate(&models.UserLogIp{})
	database.GDB.AutoMigrate(&models.AuthRule{})
	database.GDB.AutoMigrate(&models.UserAssignment{})
	// database.GDB.Migrator().CreateTable(&models.User{})
	// return nil
	return c.JSON(fiber.Map{"message": "success migration"})
}
