package database

import (
	"time"
)

func ConnectGormDB() {
	var err error

	sqlDB, err := GDB.DB()

	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetMaxOpenConns(200)
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetConnMaxLifetime(24 * time.Hour)

	if err = sqlDB.Ping(); err != nil {
		panic("failed to ping database")
	}
}
