package database

import (
	"fmt"
	"sika_apigateway/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var RedisAddrs = config.Config("REDIS_ADDRES")

var GormConn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.Config("DB_USER_POSTGRESQL"), config.Config("DB_PASSWORD_POSTGRESQL"), config.Config("DB_HOST_POSTGRESQL"), config.Config("DB_PORT_POSTGRESQL"), config.Config("DB_NAME_POSTGRESQL_APIGATEWAY"))

var GDB, _ = gorm.Open(postgres.Open(GormConn), &gorm.Config{
	PrepareStmt:            true,
	SkipDefaultTransaction: true,
})
