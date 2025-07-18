package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"go-gerbang/config"
	"go-gerbang/handlers"
	"go-gerbang/middleware"
	"go-gerbang/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func LogoutWeb(c *fiber.Ctx) error {
	currSession, err := middleware.SessionStore.Get(c)
	if err != nil {
		return handlers.UnauthorizedErrorResponse(c, err)
	}

	if err := currSession.Destroy(); err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	currSession.SetExpiry(time.Second * -60)

	token := c.Get("Authorization")
	if token == "" {
		token = c.Cookies(middleware.CookieRefreshJWT)
	}

	if token == "" {
		return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("missing refresh token"))
	}

	var refreshToken string
	if strings.HasPrefix(token, "Bearer ") {
		refreshToken = token[len("Bearer "):]
	} else {
		refreshToken = token
	}

	user, err := middleware.Verify(refreshToken, "refresh")
	if err != nil {
		return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("invalid refresh token"))
	}

	jti := user.Jti
	if jti == nil || *jti == "" || handlers.IsTokenBlacklisted(*jti) {
		return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("this refresh token is no longer valid"))
	}
	err = handlers.BlacklistToken(*jti)
	if err != nil {
		return handlers.InternalServerErrorResponse(c, fmt.Errorf("failed to blacklist token"))
	}

	domain := c.Query("domain")

	if domain == "" {
		domain = "siskor.web.id"
	}

	c.Cookie(&fiber.Cookie{
		Name:     middleware.CookieRefreshJWT,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Domain:   domain,
	})

	cookieJWT := new(fiber.Cookie)
	cookieJWT.Name = middleware.CookieJWT
	cookieJWT.Expires = time.Now().Add(-(time.Hour * 2))
	cookieJWT.HTTPOnly = true
	cookieJWT.Domain = domain
	cookieJWT.Secure = config.SecureCookies
	cookieJWT.SameSite = config.CookieSameSite
	cookieJWT.SessionOnly = false
	c.Cookie(cookieJWT)

	redirectUrl := c.Query("redirectUrl")
	if redirectUrl != "" {
		return c.Redirect(redirectUrl)
	}

	return handlers.SuccessResponse(c, true, "success logout", currSession, nil)
}

func AuthByJWT(c *fiber.Ctx) error {
	AuthKey := c.Params("token")
	url := c.Query("url")

	if url == "" {
		return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("need url params"))
	}

	user := new(models.User)

	err := models.FindUserByAuthKey(user, AuthKey).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return handlers.NotFoundErrorResponse(c, fmt.Errorf("auth key is not found"))
	}

	currSession, err := middleware.SessionStore.Get(c)
	if err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}
	defer currSession.Save()

	err = currSession.Regenerate()
	if err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	userSession := currSession.Get("User")

	if userSession == nil {
		currSession.Set(middleware.UserId, user.IdAccount)
		currSession.Set(middleware.AuthKey, user.AuthKey)
		currSession.Set(middleware.Username, user.Username)

		err := currSession.Save()
		if err != nil {
			return handlers.InternalServerErrorResponse(c, err)
		}
	}

	redirectUrl := c.Query("redirectUrl")
	if redirectUrl != "" {
		return c.Redirect(redirectUrl)
	}

	return handlers.SuccessResponse(c, true, "current session", userSession, nil)
}

func GetSessionJWT(c *fiber.Ctx) error {
	data := new(models.UserData)
	// data = &models.UserData{
	// 	IdAccount:       c.Locals("id_account").(string),
	// 	IdentityNumber:  c.Locals("identity_number").(string),
	// 	Username:        c.Locals("username").(string),
	// 	FullName:        c.Locals("full_name").(string),
	// 	Email:           c.Locals("email").(string),
	// 	PhoneNumber:     c.Locals("phone_number").(string),
	// 	DateOfBirth:     c.Locals("date_of_birth").(*time.Time),
	// 	AuthKey:         c.Locals("auth_key").(string),
	// 	UsedPin:         c.Locals("used_pin").(int8),
	// 	IsGoogleAccount: c.Locals("is_google_account").(int8),
	// 	StatusAccount:   c.Locals("status_account").(int8),
	// 	LoginIp:         c.Locals("login_ip").(string),
	// 	LoginAttempts:   c.Locals("login_attempts").(int8),
	// 	LoginTime:       c.Locals("login_time").(int64),
	// 	CreatedAt:       c.Locals("created_at").(int),
	// 	UpdatedAt:       c.Locals("updated_at").(int),
	// }
	userData, ok := c.Locals("user").(*models.UserData)
	if ok {
		data = userData
	} else {
		data = nil
	}

	return handlers.SuccessResponse(c, true, "Success Get JWT", data, nil)
}

func GetSession(c *fiber.Ctx) error {
	// currSession, err := middleware.SessionStore.Get(c)
	// if err != nil {
	// 	return handlers.UnauthorizedErrorResponse(c, err)
	// }

	// if len(currSession.Keys()) > 0 {

	// 	userId, ok := currSession.Get(middleware.UserId).(string)
	// 	if !ok {
	// 		return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("session userId is null"))
	// 	}
	// 	authKey, _ := currSession.Get(middleware.AuthKey).(string)
	// 	username, _ := currSession.Get(middleware.Username).(string)
	// 	fullName, _ := currSession.Get(middleware.FullName).(string)

	// 	data := new(models.UserData)
	// 	data = &models.UserData{
	// 		IdAccount: userId,
	// 		AuthKey:   authKey,
	// 		Username:  username,
	// 		FullName:  fullName,
	// 	}

	// 	return handlers.SuccessResponse(c, true, "success get session", data, nil)
	// }

	return handlers.SuccessResponse(c, true, "no session", nil, nil)
}
