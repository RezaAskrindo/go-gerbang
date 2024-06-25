package models

import (
	"go-gerbang/database"

	"gorm.io/gorm"
)

type UserAssignment struct {
	AccountId  string   `gorm:"type:uuid;primaryKey" json:"account_id"`
	AuthRoleId int      `gorm:"primaryKey" json:"auth_role_id"`
	User       User     `gorm:"foreignKey:IdAccount;references:AccountId" json:"user"`
	AuthRule   AuthRule `gorm:"foreignKey:IdAuthRole;references:AuthRoleId" json:"auth_rule"`
}

func CreateUserAssignment(userAssignment *UserAssignment) *gorm.DB {
	return database.GDB.Create(userAssignment)
}
