package database

import (
	"fmt"
	"go-gerbang/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var RedisAddrs = config.Config("REDIS_ADDRES")

var GormConn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.Config("DB_USER"), config.Config("DB_PASSWORD"), config.Config("DB_HOST"), config.Config("DB_PORT"), config.Config("DB_NAME_APIGATEWAY"))

var GDB, _ = gorm.Open(postgres.Open(GormConn), &gorm.Config{
	PrepareStmt:                              true,
	SkipDefaultTransaction:                   true,
	DisableForeignKeyConstraintWhenMigrating: true,
})
