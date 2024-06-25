package models

import (
	"sika_apigateway/database"

	"gorm.io/gorm"
)

type UserLogIp struct {
	AccountId string `gorm:"type:uuid;primaryKey" json:"account_id"`
	IpLogin   string `gorm:"not null;size:32" json:"ip_login"`
	TimeLogin int    `gorm:"not null" json:"time_login"`
	User      User   `gorm:"foreignKey:IdAccount;references:AccountId" json:"user"`
}

func CreateUserLogIp(userLogIp *UserLogIp) *gorm.DB {
	return database.GDB.Create(userLogIp)
}
