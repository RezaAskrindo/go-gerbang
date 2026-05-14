package services

import (
	"go-gerbang/handlers"
	"go-gerbang/models"

	"go-gerbang/database"

	"github.com/gofiber/fiber/v3"
)

func CheckMigrationStatus(c fiber.Ctx) error {
	tables := []string{
		"users",
		"user_assignments",
		"auth_rules",
		"loggers",
		// "configurations",
		"configs",
	}

	missing := []string{}

	for _, table := range tables {
		if !database.GDB.Migrator().HasTable(table) {
			missing = append(missing, table)
		}
	}

	if len(missing) > 0 {
		count := int64(len(missing))
		// getAllTables, _ := database.GDB.Migrator().GetTables()
		result := map[string]interface{}{
			"tables":  tables,
			"missing": missing,
			// "getAllTables": getAllTables,
		}
		return handlers.SuccessResponse(c, true, "success get missing table", result, &count)
	}

	return handlers.SuccessResponse(c, true, "no missing table", nil, nil)
}

func MigrationService(c fiber.Ctx) error {
	err := database.GDB.AutoMigrate(&models.User{}, &models.UserAssignment{}, &models.AuthRule{}, &models.Logger{}, &models.Config{})

	if err != nil {
		return c.JSON(fiber.Map{"message": "failed to migration", "error": err.Error()})
	}

	return handlers.SuccessResponse(c, true, "success migration", nil, nil)
}

func MigrateAdminUser(c fiber.Ctx) error {
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
