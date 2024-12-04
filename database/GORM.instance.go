package database

import (
	"fmt"
	"go-gerbang/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// var GormConn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.Config("DB_USER"), config.Config("DB_PASSWORD"), config.Config("DB_HOST"), config.Config("DB_PORT"), config.Config("DB_NAME_APIGATEWAY"))

// var GDB, _ = gorm.Open(postgres.Open(GormConn), &gorm.Config{

var dsn = fmt.Sprintf("host=%s dbname=%s port=%s user=%s password=%s sslmode=disable TimeZone=Asia/Shanghai", config.Config("DB_HOST"), config.Config("DB_NAME"), config.Config("DB_PORT"), config.Config("DB_USER"), config.Config("DB_PASSWORD"))

var GDB, _ = gorm.Open(postgres.New(postgres.Config{
	DSN:                  dsn,
	PreferSimpleProtocol: true,
}), &gorm.Config{
	PrepareStmt:                              true,
	SkipDefaultTransaction:                   true,
	DisableForeignKeyConstraintWhenMigrating: true,
})
