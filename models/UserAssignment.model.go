package models

import (
	"go-gerbang/database"

	"gorm.io/gorm"
)

type UserAssignment struct {
	AccountId  string `gorm:"type:uuid;primaryKey" json:"account_id"`
	AuthRoleId int    `gorm:"primaryKey;not null;" json:"auth_role_id"`
}

type UserAssignmentJoin struct {
	AccountId  string   `gorm:"type:uuid;primaryKey" json:"account_id"`
	AuthRoleId int      `gorm:"primaryKey;not null;" json:"auth_role_id"`
	User       User     `gorm:"foreignKey:IdAccount;references:AccountId" json:"user"`
	AuthRule   AuthRule `gorm:"foreignKey:IdAuthRole;references:AuthRoleId" json:"auth_rule"`
}

func CreateUserAssignment(userAssignment *UserAssignment) *gorm.DB {
	return database.GDB.Create(userAssignment)
}

func CreateUserAssignments(assignments []UserAssignment) *gorm.DB {
	return database.GDB.Create(&assignments)
}

func UpdateUserAssignment(userAssignment *UserAssignment) *gorm.DB {
	return database.GDB.Save(userAssignment)
}

func UpdateUserAssignments(assignments []UserAssignment) error {
	tx := database.GDB.Begin()
	for _, assignment := range assignments {
		if err := tx.Save(&assignment).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func DeleteUserAssignment(accountId string, authRoleId int) *gorm.DB {
	return database.GDB.Delete(&UserAssignment{}, "account_id = ? AND auth_role_id = ?", accountId, authRoleId)
}

func DeleteUserAssignments(assignments []UserAssignment) *gorm.DB {
	return database.GDB.Delete(&assignments)
}

func FindUserAssignment(dest interface{}, conds ...interface{}) *gorm.DB {
	return database.GDB.Model(&UserAssignment{}).Preload("User").Preload("AuthRule").Find(dest, conds...)
}

func CountFindUserAssignment(count *int64, conds ...interface{}) error {
	tx := database.GDB.Model(&UserAssignment{})
	if len(conds) > 0 {
		tx = tx.Where(conds[0], conds[1:]...)
	}
	return tx.Count(count).Error
}
