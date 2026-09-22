package database

import (
	"go-gerbang/config"
	"log"
	"time"
)

func ConnectGormDB() {
	sqlDB, _ := GDB.DB()

	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetConnMaxLifetime(5 * time.Hour)

	log.Printf("Successfully connected to the %s, %s:%s", config.Config("DB_USER_POSTGRESQL"), config.Config("DB_NAME_POSTGRESQL_APIGATEWAY"), config.Config("DB_PORT_POSTGRESQL"))
}
