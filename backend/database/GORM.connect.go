package database

import (
	"go-gerbang/config"
	"log"
	"time"
)

func ConnectGormDB() {
	sqlDB, _ := GDB.DB()

	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetMaxOpenConns(200)
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetConnMaxLifetime(24 * time.Hour)

	log.Printf("Successfully connected to the %s, %s:%s", config.Config("DB_USER_POSTGRESQL"), config.Config("DB_NAME_POSTGRESQL_APIGATEWAY"), config.Config("DB_PORT_POSTGRESQL"))
}
