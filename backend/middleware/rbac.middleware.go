package middleware

import (
	"fmt"

	"go-gerbang/config"
	"go-gerbang/database"
	"go-gerbang/handlers"
	"go-gerbang/models"

	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	_ "github.com/jackc/pgx/v5"

	"github.com/gofiber/fiber/v3"
)

type CasbinRule struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Ptype string `gorm:"size:255;uniqueIndex:unique_index"`
	V0    string `gorm:"size:255;uniqueIndex:unique_index"`
	V1    string `gorm:"size:255;uniqueIndex:unique_index"`
	V2    string `gorm:"size:128;uniqueIndex:unique_index"`
	V3    string `gorm:"default:null;size:128;uniqueIndex:unique_index"`
	V4    string `gorm:"default:null;size:128;uniqueIndex:unique_index"`
	V5    string `gorm:"default:null;size:128;uniqueIndex:unique_index"`
}

func (CasbinRule) TableName() string {
	return "casbin_rule"
}

var enforcer *casbin.Enforcer

func InitCasbin() error {
	adapter, err := gormadapter.NewAdapterByDBWithCustomTable(database.GDB, &CasbinRule{})
	if err != nil {
		return err
	}
	enforcer, err = casbin.NewEnforcer(config.BasePath+config.Config("CONFIG_PATH_CASBIN_MODEL"), adapter)
	if err != nil {
		return err
	}
	if err := enforcer.LoadPolicy(); err != nil {
		return err
	}
	return nil
}

func AuthRBAC(c fiber.Ctx) error {
	if enforcer == nil {
		return handlers.InternalServerErrorResponse(c, fmt.Errorf("casbin not initialized"))
	}

	user, ok := c.Locals("user").(*models.UserData)
	if !ok || user == nil {
		return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("user not found"))
	}

	if len(user.UserAssignments) == 0 {
		return handlers.ForbiddenErrorResponse(c, fmt.Errorf("no roles assigned"))
	}

	obj := c.Path()
	act := c.Method()

	for _, ua := range user.UserAssignments {
		sub := fmt.Sprintf("role:%d", ua.AuthRoleId)

		allowed, err := enforcer.Enforce(sub, obj, act)
		if err != nil {
			return handlers.InternalServerErrorResponse(c, err)
		}

		if allowed {
			return c.Next()
		}
	}

	return handlers.ForbiddenErrorResponse(c, fmt.Errorf("access denied"))
}

// func AuthRBACBackup(c fiber.Ctx) error {
// 	user, ok := c.Locals("user").(*models.UserData)
// 	if !ok {
// 		return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("user not found"))
// 	}

// 	a, _ := gormadapter.NewAdapterByDBWithCustomTable(database.GDB, &CasbinRule{})
// 	authz, err := casbin.NewEnforcer(config.BasePath+config.Config("CONFIG_PATH_CASBIN_MODEL"), a)
// 	if err != nil {
// 		log.Fatalf("error: enforcer: %s", err)
// 	}

// 	for _, ua := range user.UserAssignments {
// 		role := fmt.Sprintf("role:%d", ua.AuthRoleId)
// 		fmt.Println(role)
// 		allowed, _ := authz.Enforce(role, c.Path(), c.Method())
// 		if allowed {
// 			return c.Next()
// 		}
// 	}

// 	return handlers.ForbiddenErrorResponse(c, fmt.Errorf("your role don't have access"))
// }
