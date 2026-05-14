package models

import (
	"go-gerbang/database"

	"gorm.io/gorm"
)

type AuthRule struct {
	IdAuthRole   int    `gorm:"type:autoIncrement;primaryKey" json:"id_auth_role"`
	NameAuthRole string `gorm:"not null;size:64" json:"name_auth_role" validate:"required"`
	DescAuthRole string `gorm:"default:null;size:512" json:"desc_auth_role"`
	CreatedAt    int    `gorm:"autoCreateTime:true" json:"created_at"`
	UpdatedAt    int    `gorm:"default:0;autoCreateTime:false" json:"updated_at"`
}

func CreateAuthRule(authRule *AuthRule) *gorm.DB {
	return database.GDB.Create(authRule)
}

func UpdateAuthRule(authRoleId interface{}, data interface{}) *gorm.DB {
	return database.GDB.Model(&AuthRule{}).Where("id_auth_role = ?", authRoleId).Updates(data)
}

func DeleteAuthRule(authRoleId interface{}) *gorm.DB {
	return database.GDB.Unscoped().Delete(&User{}, "id_auth_role = ?", authRoleId)
}

func FindAuthRule(dest interface{}, conds ...interface{}) *gorm.DB {
	return database.GDB.Model(&AuthRule{}).Find(dest, conds...)
}

func CountFindAuthRule(count *int64, conds ...interface{}) error {
	tx := database.GDB.Model(&AuthRule{})
	if len(conds) > 0 {
		tx = tx.Where(conds[0], conds[1:]...)
	}
	return tx.Count(count).Error
}
