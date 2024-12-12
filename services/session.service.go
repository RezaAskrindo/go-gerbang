package services

import (
	"errors"
	"fmt"
	"time"

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

	cookieJWT := new(fiber.Cookie)
	cookieJWT.Name = middleware.CookieJWT
	// cookieJWT.Expires = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	cookieJWT.Expires = time.Now().Add(-(time.Hour * 2))
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
	data = &models.UserData{
		IdAccount:       c.Locals("id_account").(string),
		IdentityNumber:  c.Locals("identity_number").(string),
		Username:        c.Locals("username").(string),
		FullName:        c.Locals("full_name").(string),
		Email:           c.Locals("email").(string),
		PhoneNumber:     c.Locals("phone_number").(string),
		DateOfBirth:     c.Locals("date_of_birth").(*time.Time),
		AuthKey:         c.Locals("auth_key").(string),
		UsedPin:         c.Locals("used_pin").(int8),
		IsGoogleAccount: c.Locals("is_google_account").(int8),
		StatusAccount:   c.Locals("status_account").(int8),
		LoginIp:         c.Locals("login_ip").(string),
		LoginAttempts:   c.Locals("login_attempts").(int8),
		LoginTime:       c.Locals("login_time").(int64),
		CreatedAt:       c.Locals("created_at").(int),
		UpdatedAt:       c.Locals("updated_at").(int),
	}

	return handlers.SuccessResponse(c, true, "Success Get JWT", data, nil)
}

func GetSession(c *fiber.Ctx) error {
	currSession, err := middleware.SessionStore.Get(c)
	if err != nil {
		return handlers.UnauthorizedErrorResponse(c, err)
	}

	if len(currSession.Keys()) > 0 {

		userId, ok := currSession.Get(middleware.UserId).(string)
		if !ok {
			return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("session userId is null"))
		}
		authKey, _ := currSession.Get(middleware.AuthKey).(string)
		username, _ := currSession.Get(middleware.Username).(string)
		fullName, _ := currSession.Get(middleware.FullName).(string)

		data := new(models.UserData)
		data = &models.UserData{
			IdAccount: userId,
			AuthKey:   authKey,
			Username:  username,
			FullName:  fullName,
		}

		return handlers.SuccessResponse(c, true, "success get session", data, nil)
	}

	return handlers.SuccessResponse(c, true, "no session", nil, nil)
}
