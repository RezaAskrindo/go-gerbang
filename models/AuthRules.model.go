package models

import (
	"sika_apigateway/database"

	"gorm.io/gorm"
)

type AuthRule struct {
	IdAuthRole   int    `gorm:"type:autoIncrement;primaryKey" json:"id_auth_role"`
	NameAuthRole string `gorm:"not null;size:64" json:"name_auth_role"`
	DescAuthRole string `gorm:"default:null;size:512" json:"desc_auth_role"`
	CreatedAt    int    `gorm:"autoCreateTime:true" json:"created_at"`
	UpdatedAt    int    `gorm:"default:0;autoCreateTime:false" json:"updated_at"`
}

func CreateAuthRule(authRule *AuthRule) *gorm.DB {
	return database.GDB.Create(authRule)
}
