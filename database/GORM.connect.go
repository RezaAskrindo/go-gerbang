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

	log.Printf("Successfully connected to the %s database: %s", config.Config("DB_CONNECTION"), config.Config("DB_NAME"))
}
