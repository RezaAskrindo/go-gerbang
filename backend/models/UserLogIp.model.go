package models

import (
	"go-gerbang/database"

	"gorm.io/gorm"
)

type UserLogIp struct {
	AccountId string `gorm:"type:uuid;primaryKey" json:"account_id"`
	IpLogin   string `gorm:"primaryKey;not null;size:32" json:"ip_login"`
	TimeLogin int    `gorm:"not null" json:"time_login"`
	Users     []User `gorm:"foreignKey:IdAccount;references:AccountId" json:"users"`
}

func CreateUserLogIp(userLogIp *UserLogIp) *gorm.DB {
	return database.GDB.Create(userLogIp)
}
