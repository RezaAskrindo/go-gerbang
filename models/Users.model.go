package models

import (
	"go-gerbang/database"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	// IdAccount          uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id_account"`
	IdAccount          string     `gorm:"type:uuid;primaryKey" json:"id_account"`
	IdentityNumber     string     `gorm:"default:null;size:64" json:"identity_number"`
	Username           string     `gorm:"not null;size:128;index:idx_username,unique" json:"username"`
	FullName           string     `gorm:"not null;size:128" json:"full_name"`
	Email              string     `gorm:"default:null;size:128" json:"email"`
	PhoneNumber        string     `gorm:"default:null;size:13" json:"phone_number"`
	DateOfBirth        *time.Time `gorm:"default:null" json:"date_of_birth"`
	StatusAccount      int8       `gorm:"default:0" json:"status_account"`
	AuthKey            string     `gorm:"default:null;size:32" json:"auth_key"`
	PasswordHash       string     `gorm:"default:null;size:256"`
	PasswordResetToken *string    `gorm:"default:null;size:256"`
	AccessToken        *string    `gorm:"default:null;size:256"`
	PinHash            *string    `gorm:"default:null;size:256"`
	UsedPin            int8       `gorm:"default:0" json:"used_pin"`
	IsGoogleAccount    int8       `gorm:"default:0" json:"is_google_account"`
	LoginIp            string     `gorm:"default:null;size:32" json:"login_ip"`
	LoginAttempts      int8       `gorm:"default:0" json:"login_attempts"`
	LoginTime          int64      `gorm:"default:0" json:"login_time"`
	CreatedAt          int        `gorm:"autoCreateTime:true" json:"created_at"`
	UpdatedAt          int        `gorm:"default:0;autoCreateTime:false" json:"updated_at"`
}

type UserData struct {
	// IdAccount       uuid.UUID  `json:"id_account"`
	IdAccount       string     `json:"id_account"`
	IdentityNumber  string     `json:"identity_number"`
	Username        string     `json:"username"`
	FullName        string     `json:"full_name"`
	Email           string     `json:"email"`
	PhoneNumber     string     `json:"phone_number"`
	DateOfBirth     *time.Time `json:"date_of_birth"`
	StatusAccount   int8       `json:"status"`
	AuthKey         string     `json:"auth_key"`
	UsedPin         int8       `json:"used_pin"`
	IsGoogleAccount int8       `json:"is_google_account"`
	LoginIp         string     `json:"login_ip"`
	LoginAttempts   int8       `json:"login_attempts"`
	LoginTime       int64      `json:"login_time"`
	CreatedAt       int        `json:"created_at"`
	UpdatedAt       int        `json:"updated_at"`
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

// func (User) TableName() string {
// 	return "user"
// }

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.IdAccount = uuid.New().String()
	// user.IdAccount = uuid.New()

	return nil
}

func CreateUser(user *User) *gorm.DB {
	return database.GDB.Create(user)
}

func CeneratePasswordResetToken(email interface{}, randomReset string) *gorm.DB {
	return database.GDB.Model(&User{}).Where("email = ?", email).Update("password_reset_token", randomReset)
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
	return database.GDB.Model(&User{}).Select("status").Where("id_account = ?", accountId).Update("status", 0)
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

func FindUserByIdRaw(dest interface{}, body int) *gorm.DB {
	return database.GDB.Raw("SELECT * FROM users WHERE id_account = ?", body).First(dest)
}

func FindUserById(accountId int) (*User, error) {
	User := new(User)
	if response := database.GDB.Where("id_account = ? AND users.status_account = ?", accountId, STATUS_ACTIVE).First(&User); response.Error != nil {
		return nil, response.Error
	}
	// if User.IdAccount == "" {
	// 	return User, errors.New("user not found")
	// }

	return User, nil
}

func FindUserByIdentity(identity string) (*User, error) {
	User := new(User)
	if response := database.GDB.Where("LOWER(username) = LOWER(?) OR LOWER(email) = LOWER(?) OR phone_number = ? AND users.status_account = ?", identity, identity, identity, STATUS_ACTIVE).First(&User); response.Error != nil {
		return nil, response.Error
	}
	// if User.IdAccount == "" {
	// 	return User, errors.New("user not found")
	// }

	return User, nil
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
