package database

import (
	"fmt"
	"go-gerbang/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// var dsn = fmt.Sprintf("host=%s dbname=%s port=%s user=%s password=%s sslmode=disable TimeZone=Asia/Shanghai", config.Config("DB_HOST_POSTGRESQL"), config.Config("DB_NAME_POSTGRESQL_APIGATEWAY"), config.Config("DB_PORT_POSTGRESQL"), config.Config("DB_USER_POSTGRESQL"), config.Config("DB_PASSWORD_POSTGRESQL"))
var dsn = fmt.Sprintf("host=%s dbname=%s port=%s user=%s password=%s sslmode=disable", config.Config("DB_HOST_POSTGRESQL"), config.Config("DB_NAME_POSTGRESQL_APIGATEWAY"), config.Config("DB_PORT_POSTGRESQL"), config.Config("DB_USER_POSTGRESQL"), config.Config("DB_PASSWORD_POSTGRESQL"))

var GDB, _ = gorm.Open(postgres.New(postgres.Config{
	DSN:                  dsn,
	PreferSimpleProtocol: true,
}), &gorm.Config{
	PrepareStmt:                              true,
	SkipDefaultTransaction:                   true,
	DisableForeignKeyConstraintWhenMigrating: true,
})
