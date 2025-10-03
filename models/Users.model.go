package models

import (
	"errors"
	"fmt"
	"time"

	"go-gerbang/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	IdAccount          uuid.UUID        `gorm:"type:uuid;primaryKey" json:"idAccount"`
	IdentityNumber     string           `gorm:"default:null;size:64" json:"identityNumber"`
	Username           string           `gorm:"not null;size:128;unique" json:"username" validate:"required"`
	FullName           string           `gorm:"not null;size:128" json:"fullName" validate:"required"`
	Email              string           `gorm:"default:null;size:128" json:"email"`
	PhoneNumber        string           `gorm:"default:null;size:13" json:"phoneNumber"`
	DateOfBirth        *time.Time       `gorm:"default:null" json:"dateOfBirth"`
	StatusAccount      int8             `gorm:"default:0" json:"statusAccount"`
	AuthKey            string           `gorm:"default:null;size:32" json:"authKey"`
	Password           string           `gorm:"-" json:"password"`
	PasswordHash       string           `gorm:"default:null;size:256" json:"-"`
	PasswordResetToken *string          `gorm:"default:null;size:256" json:"-"`
	AccessToken        *string          `gorm:"default:null;size:256" json:"-"`
	PinHash            *string          `gorm:"default:null;size:256" json:"-"`
	UsedPin            int8             `gorm:"default:0" json:"usedPin"`
	IsGoogleAccount    int8             `gorm:"default:0" json:"isGoogleAccount"`
	LoginIp            string           `gorm:"default:null;size:32" json:"loginIp"`
	LoginAttempts      int8             `gorm:"default:0" json:"loginAttempts"`
	LoginTime          int64            `gorm:"default:0" json:"loginTime"`
	CreatedAt          int              `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt          int              `gorm:"default:0;autoUpdateTime" json:"updatedAt"`
	UserAssignments    []UserAssignment `gorm:"foreignKey:AccountId;references:IdAccount"`
}

type UserData struct {
	// IdAccount       uuid.UUID  `json:"idAccount"`
	IdAccount       string           `json:"idAccount"`
	IdentityNumber  string           `json:"identityNumber"`
	Username        string           `json:"username"`
	FullName        string           `json:"fullName"`
	Email           string           `json:"email"`
	PhoneNumber     string           `json:"phoneNumber"`
	DateOfBirth     *time.Time       `json:"dateOfBirth"`
	StatusAccount   int8             `json:"statusAccount"`
	AuthKey         string           `json:"authKey"`
	UsedPin         int8             `json:"usedPin"`
	IsGoogleAccount int8             `json:"isGoogleAccount"`
	LoginIp         string           `json:"loginIp"`
	LoginAttempts   int8             `json:"loginAttempts"`
	LoginTime       int64            `json:"loginTime"`
	UserAssignments []UserAssignment `gorm:"foreignKey:AccountId;references:IdAccount" json:"user_assignments"`
	CreatedAt       int              `json:"-"`
	UpdatedAt       int              `json:"-"`
	Jti             *string
}

type UserDataResponse struct {
	Item *UserData `json:"item"`
}

type UserDataResponses struct {
	Items *[]UserData `json:"items"`
}

var STATUS_DELETED = 0
var STATUS_INACTIVE = 9
var STATUS_ACTIVE = 10

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.IdAccount = uuid.New()

	return nil
}

func CreateUser(user *User) *gorm.DB {
	return database.GDB.Create(user)
}

func CeneratePasswordResetToken(idAccount interface{}, randomReset string) *gorm.DB {
	return database.GDB.Model(&User{}).Where("id_account = ?", idAccount).Update("password_reset_token", randomReset)
}

func CenerateTimeBasedOneTimeTokenByEmail(email interface{}, TOTP string) *gorm.DB {
	return database.GDB.Model(&User{}).Where("email = ?", email).Update("access_token", TOTP)
}

func CenerateTimeBasedOneTimeTokenByPhone(phone_number interface{}, TOTP string) *gorm.DB {
	return database.GDB.Model(&User{}).Where("phone_number = ?", phone_number).Update("access_token", TOTP)
}

func SetNullAccessTokenUser(accountId interface{}) *gorm.DB {
	return database.GDB.Model(&User{}).Select("AccessToken").Where("id_account = ?", accountId).Update("access_token", nil)
}

func UpdateUser(accountId interface{}, data interface{}) *gorm.DB {
	return database.GDB.Model(&User{}).Where("id_account = ?", accountId).Updates(data)
}

func BlockUser(accountId interface{}) *gorm.DB {
	return database.GDB.Model(&User{}).Select("status_account").Where("id_account = ?", accountId).Update("status_account", 0)
}

func UpdateUserPassword(accountId interface{}, data interface{}) *gorm.DB {
	return database.GDB.Model(&User{}).Select("password_hash", "password_reset_token").Where("id_account = ?", accountId).Updates(data)
}

func GenerateAuthKeyUser(accountId interface{}, data interface{}) *gorm.DB {
	return database.GDB.Model(&User{}).Select("AuthKey").Where("id_account = ?", accountId).Update("auth_key", data)
}

func FindUserByAuthKey(dest interface{}, body string) *gorm.DB {
	return database.GDB.Raw("SELECT * FROM users WHERE auth_key = ?", body).First(dest)
}

func FindUserByPasswordReset(dest interface{}, body string) *gorm.DB {
	return database.GDB.Raw("SELECT * FROM users WHERE password_reset_token = ?", body).First(dest)
}

func FindUserById(dest interface{}, idAccount interface{}) error {
	err := database.GDB.Raw("SELECT * FROM users WHERE id_account = ?", idAccount).First(dest).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("user is not found")
	}

	return nil
}

func FindUserByIdentity(dest interface{}, username interface{}, email interface{}, phoneNumber interface{}, identityNumber interface{}) error {
	// err := database.GDB.Raw("SELECT * FROM users WHERE LOWER(username) = LOWER(?) OR LOWER(email) = LOWER(?) OR phone_number = ? OR identity_number = ?", username, email, phoneNumber, identityNumber).InnerJoins("UserAssignment").First(dest).Error
	err := database.GDB.
		Where("LOWER(users.username) = LOWER(?) OR LOWER(users.email) = LOWER(?) OR users.phone_number = ? OR users.identity_number = ?", username, email, phoneNumber, identityNumber).
		Preload("UserAssignments").
		First(dest).Error
	if err != nil {
		return err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("user is not found")
	}

	return nil
}

func FindUser(dest interface{}, conds ...interface{}) *gorm.DB {
	return database.GDB.Model(&User{}).Select("id_account", "identity_number", "username", "full_name", "email", "phone_number", "date_of_birth", "status_account", "login_ip", "login_attempts", "login_time", "created_at", "updated_at").Find(dest, conds...)
}

func FindAllUser(dest interface{}) *gorm.DB {
	return FindUser(dest)
}

func FindUserDataById(dest interface{}, accountId interface{}) *gorm.DB {
	return FindUser(dest, "id_account = ?", accountId)
}

// func HardDeleteUser(idAccount interface{}) *gorm.DB {
// 	return database.GDB.Unscoped().Where("id_account = ?", idAccount).Delete(&User{})
// }

func HardDeleteUser(idAccount string) error {
	tx := database.GDB.Unscoped().Where("id_account = ?", idAccount).Delete(&User{})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return fmt.Errorf("no user found with id_account = %s", idAccount)
	}

	return nil
}
