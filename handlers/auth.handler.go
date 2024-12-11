package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-gerbang/config"
	"go-gerbang/database"
	"go-gerbang/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/oauth2/v2"
)

var TimeNow = time.Now()
var ExpiredReset = int64(24 * 60 * 60 * 1000)
var ExpiredTOTP = int64(2 * 60 * 1000)
var HttpClient = &http.Client{}

func GeneratePasswordHash(raw string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(raw), 10)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateResetRandom(length int) string {
	var time = strconv.FormatInt(int64(TimeNow.UnixMilli()), 10)

	reset_random := RandomString(32) + "_" + time

	return string(reset_random)
}

func IsPasswordResetTokenValid(token string) bool {
	chunks := strings.Split(token, "_")
	if len(chunks) < 2 {
		return false
	}

	intUnix, err := strconv.ParseInt(chunks[1], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	return intUnix+ExpiredReset >= TimeNow.Unix()
}

func GenerateTOTPToken(token string) string {
	var time = strconv.FormatInt(int64(TimeNow.UnixMilli()), 10)

	reset_random := GeneratePasswordHash(token) + "_" + time

	return string(reset_random)
}

func IsTOTPValid(token string) bool {
	chunks := strings.Split(token, "_")
	if len(chunks) < 2 {
		return false
	}

	intUnix, err := strconv.ParseInt(chunks[1], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	return intUnix+ExpiredTOTP >= TimeNow.Unix()
}

// GOOGLE VALIDATION SIGN IN
func VerifyIdTokenGoogle(idToken string) (*oauth2.Tokeninfo, error) {
	oauth2Service, err := oauth2.New(HttpClient)
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(idToken)
	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		return nil, err
	}
	return tokenInfo, nil
}

// SUPPORT FUNCTION
func SendSafeUserData(user *models.User, randString string) models.UserData {
	return models.UserData{
		IdAccount:      user.IdAccount.String(),
		IdentityNumber: user.IdentityNumber,
		Username:       user.Username,
		FullName:       user.FullName,
		Email:          user.Email,
		PhoneNumber:    user.PhoneNumber,
		AuthKey:        randString,
		StatusAccount:  user.StatusAccount,
		UsedPin:        user.UsedPin,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
		LoginIp:        user.LoginIp,
		LoginAttempts:  user.LoginAttempts,
		LoginTime:      user.LoginTime,
	}
}

func GenerateTokenJWT(user_data models.UserData, c *fiber.Ctx) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id_account"] = user_data.IdAccount
	claims["identity_number"] = user_data.IdentityNumber
	claims["username"] = user_data.Username
	claims["full_name"] = user_data.FullName
	claims["email"] = user_data.Email
	claims["phone_number"] = user_data.PhoneNumber
	claims["date_of_birth"] = user_data.DateOfBirth
	claims["auth_key"] = user_data.AuthKey
	claims["used_pin"] = user_data.UsedPin
	claims["is_google_account"] = user_data.IsGoogleAccount
	claims["status_account"] = user_data.StatusAccount
	claims["login_ip"] = user_data.LoginIp
	claims["login_attempts"] = user_data.LoginAttempts
	claims["login_time"] = user_data.LoginTime
	claims["created_at"] = user_data.CreatedAt
	claims["updated_at"] = user_data.UpdatedAt
	claims["exp"] = time.Now().Add(config.AuthTimeCache).Unix()

	t, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", err
	}

	return t, nil
}

func ValidateUserLoginIp(user_data models.UserData, c *fiber.Ctx) error {
	var count int64

	database.GDB.Table("user_log_ip").Where("user_id = ? AND ip_login = ?", user_data.IdAccount, user_data.LoginIp).Count(&count)

	if count == 0 && user_data.LoginIp != "127.0.0.1" {
		return errors.New("you're login in new IP, please identify yourself first")
	}

	return nil
}
