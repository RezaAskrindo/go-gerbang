package middleware

import (
	"go-gerbang/config"
	"time"

	"github.com/gofiber/fiber/v3"
)

func authCookie(name, value, domain string, ttl time.Duration) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/", // always set explicitly — don't rely on default
		Domain:   domain,
		Expires:  time.Now().Add(ttl),
		HTTPOnly: true,
		Secure:   config.SecureCookies,
		SameSite: config.CookieSameSite,
	}
}

func SetAuthCookies(c fiber.Ctx, domain, refreshToken, accessToken string) {
	c.Cookie(authCookie(CookieRefreshJWT, refreshToken, domain, config.RefreshAuthTimeCache))
	c.Cookie(authCookie(CookieJWT, "Bearer "+accessToken, domain, config.AuthTimeCache))
}

func ClearAuthCookies(c fiber.Ctx, domain string) {
	c.Cookie(authCookie(CookieRefreshJWT, "", domain, -time.Hour))
	c.Cookie(authCookie(CookieJWT, "", domain, -time.Hour))
}
