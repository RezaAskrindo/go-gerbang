package services

import (
	"errors"
	"go-gerbang/handlers"
	"go-gerbang/models"

	"go-gerbang/database"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func CheckMigrationStatus(c fiber.Ctx) error {
	tables := []string{
		"users",
		"user_assignments",
		"auth_rules",
		"loggers",
		"configurations",
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

	if database.GDB.Migrator().HasTable(&models.User{}) {
		err := database.GDB.Where("username = ?", "admin").First(&models.User{}).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result := map[string]interface{}{
				"tables":  []interface{}{"admin"},
				"missing": []interface{}{"admin"},
			}
			return handlers.SuccessResponse(c, true, "admin account is missing table", result, nil)
		}
	}

	return handlers.SuccessResponse(c, true, "no missing table", nil, nil)
}

func MigrationService(c fiber.Ctx) error {
	err := database.GDB.AutoMigrate(
		&models.User{},
		&models.UserAssignment{},
		&models.AuthRule{},
		&models.Logger{},
		&models.Configuration{},
	)

	if err != nil {
		return c.JSON(fiber.Map{"message": "failed to migration", "error": err.Error()})
	}

	return handlers.SuccessResponse(c, true, "success migration", nil, nil)
}

func MigrateAdminUser(c fiber.Ctx) error {
	u := new(struct {
		Password string `json:"password"`
	})

	if err := handlers.ParseBody(c, u); err != nil {
		return handlers.BadRequestErrorResponse(c, err)
	}

	if err := handlers.ValidateStruct(*u); err != nil {
		return handlers.SuccessResponse(c, false, "error validation user", err, nil)
	}

	if database.GDB.Migrator().HasTable(&models.User{}) {
		errCreate := database.GDB.Create(&models.User{
			Username:      "admin",
			FullName:      "Admin",
			StatusAccount: 10,
			// PasswordHash:  handlers.GeneratePasswordHash("@dmin9192"),
			PasswordHash: handlers.GeneratePasswordHash(u.Password),
		})

		if errCreate.Error != nil {
			// IF ERROR CACHE, RUN THIS IN SQL:
			// DISCARD ALL;
			return c.JSON(fiber.Map{"message": "failed to create admin", "error": errCreate.Error})
		}
		return handlers.SuccessResponse(c, true, "success create admin", nil, nil)
	} else {
		return c.JSON(fiber.Map{"message": "table user is not found"})
	}
}
