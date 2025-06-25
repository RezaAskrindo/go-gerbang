package handlers

import (
	"context"
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
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/idtoken"
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
	createdTime := time.Now()
	expiryTime := createdTime.Add(24 * time.Hour)
	expiryUnixMilli := expiryTime.UnixMilli()
	expiryTimeStr := strconv.FormatInt(expiryUnixMilli, 10)

	resetRandom := RandomStringV1(length) + "_" + expiryTimeStr
	return resetRandom
}

func IsPasswordResetTokenValid(token string) bool {
	chunks := strings.Split(token, "_")
	if len(chunks) < 2 {
		return false
	}

	intUnixMilli, err := strconv.ParseInt(chunks[1], 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	intUnixSec := intUnixMilli / 1000

	nowSec := time.Now().Unix()

	return nowSec <= intUnixSec+ExpiredReset
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
func VerifyIdTokenGoogle(ctx context.Context, idToken string, audience string) (*idtoken.Payload, error) {
	payload, err := idtoken.Validate(ctx, idToken, audience)
	if err != nil {
		return nil, err
	}

	return payload, nil
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

// token := jwt.New(jwt.SigningMethodHS256)
func GenerateTokenJWT(user_data models.UserData, isRefresh bool) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)

	claims := token.Claims.(jwt.MapClaims)
	claims["id_account"] = user_data.IdAccount
	claims["full_name"] = user_data.FullName
	claims["email"] = user_data.Email

	if isRefresh {
		claims["identity_number"] = user_data.IdentityNumber
		claims["username"] = user_data.Username
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
	}
	// claims["exp"] = time.Now().Add(config.AuthTimeCache).Unix()

	if isRefresh {
		claims["exp"] = time.Now().Add(config.RefreshAuthTimeCache).Unix()
		claims["typ"] = "refresh"
		claims["jti"] = uuid.New().String()
	} else {
		claims["exp"] = time.Now().Add(config.AuthTimeCache).Unix()
		claims["typ"] = "access"
	}

	t, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", err
	}

	if isRefresh {
		err = database.RedisDb.Set(database.RedisCtx, "refresh:"+user_data.IdAccount, t, config.RefreshAuthTimeCache).Err()
		if err != nil {
			return "", err
		}
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

func GenerateRefreshToken(user_data models.UserData) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)

	claims := token.Claims.(jwt.MapClaims)
	claims["id_account"] = user_data.IdAccount
	claims["username"] = user_data.Username
	claims["email"] = user_data.Email
	claims["typ"] = "refresh"
	claims["exp"] = time.Now().Add(7 * 24 * time.Hour).Unix()

	claims["jti"] = uuid.New().String()

	t, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", err
	}

	err = database.RedisDb.Set(database.RedisCtx, "refresh:"+user_data.IdAccount, t, 7*24*time.Hour).Err()
	if err != nil {
		return "", err
	}

	return t, nil
}

func IsTokenBlacklisted(jti string) bool {
	// Check if the jti exists in the blacklist
	val, err := database.RedisDb.Get(database.RedisCtx, "blacklist:"+jti).Result()
	if err == redis.Nil {
		// Token is not blacklisted
		return false
	}
	if err != nil {
		// Error with Redis
		return true
	}

	// If a value is found, it's blacklisted
	return val == "1"
}

// Add a token to the blacklist
func BlacklistToken(jti string) error {
	err := database.RedisDb.Set(database.RedisCtx, "blacklist:"+jti, "1", time.Hour*24).Err()
	if err != nil {
		return err
	}
	return nil
}
