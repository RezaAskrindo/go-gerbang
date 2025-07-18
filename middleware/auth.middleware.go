package middleware

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"go-gerbang/config"
	"go-gerbang/database"
	"go-gerbang/handlers"
	"go-gerbang/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/golang-jwt/jwt"
	"github.com/steambap/captcha"
)

// DOMAINESIA NOT SUPPORT
// var StorageRedisFiber = redis.New(redis.Config{
// 	URL: config.Config("REDIS_ADDRESS_FULL"),
// })

// var StorageRedisFiber = redis.New()

const (
	UserId           = "userId"
	AuthKey          = "authKey"
	Username         = "username"
	FullName         = "fullName"
	UserActive       = "user-active"
	CookieJWT        = "__SGJwt"
	CookieRefreshJWT = "__SGRefreshJwt"
	CookieSession    = "__SGSession"
)

var SessionStore = session.New(session.Config{
	// Expiration:     config.AuthTimeCache,
	KeyLookup:      "cookie:__SGSession",
	CookieHTTPOnly: true,
	CookieSecure:   config.SecureCookies,
	CookieSameSite: config.CookieSameSite,
	KeyGenerator: func() string {
		return handlers.RandomStringV1(64)
	},
	// Storage:        StorageRedisFiber,
})

var CsrfActivated = false
var CsrfContextKey = "token_csrf"
var CsrfHeaderName = "X-SGCsrf-Token"

var CsrfProtection = csrf.New(csrf.Config{
	Next: func(c *fiber.Ctx) bool {
		return CsrfActivated
	},
	KeyLookup:      "header:" + CsrfHeaderName,
	CookieName:     "csrf_header",
	CookieSameSite: "Lax",
	CookieSecure:   config.SecureCookies,
	Expiration:     config.CsrfTimeCache,
	ContextKey:     CsrfContextKey,
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		// ERROR ON REFERER CHECKING
		// https://docs.gofiber.io/api/middleware/csrf#referer-checking
		switch err {
		case csrf.ErrTokenNotFound:
			log.Println("Error CSRF: Token not found")
		case csrf.ErrTokenInvalid:
			log.Println("Error CSRF: Token invalid")
		case csrf.ErrNoReferer:
			log.Println("Error CSRF: No Referer header")
		case csrf.ErrBadReferer:
			log.Println("Error CSRF: Bad Referer header")
		default:
			log.Printf("Error CSRF: %v\n", err)
		}
		tokenHeader := c.Get("X-SGCsrf-Token")
		tokenCookie := c.Cookies("csrf_header")
		log.Printf("Debug — Header: %q, Cookie: %q\n", tokenHeader, tokenCookie)
		return handlers.ForbiddenErrorResponse(c, fmt.Errorf("CSRF validation failed"))
	},
	Extractor: csrf.CsrfFromHeader(CsrfHeaderName),
	KeyGenerator: func() string {
		return handlers.RandomStringV1(32)
	},
	Session: CsrfStore,
})

var CsrfStore = session.New(session.Config{
	Expiration:     config.CsrfTimeCache,     // Expire sessions after 30 minutes of inactivity
	KeyLookup:      "cookie:__SGCsrfSession", // Recommended to use the __Host- prefix when serving the app over TLS
	CookieSecure:   config.SecureCookies,
	CookieHTTPOnly: true,
	CookieSameSite: "Lax",
})

var CsrfProtectionCookies = csrf.New(csrf.Config{
	Session: CsrfStore,
	Next: func(c *fiber.Ctx) bool {
		return CsrfActivated
	},
	KeyLookup:      "cookie:__SGCsrf",
	CookieName:     "__SGCsrf",
	CookieHTTPOnly: false,
	CookieSameSite: "Lax",
	Expiration:     config.CsrfTimeCache,
	ContextKey:     CsrfContextKey,
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return handlers.ForbiddenErrorResponse(c, fmt.Errorf("CSRF Token is invalid"))
	},
	KeyGenerator: func() string {
		return handlers.RandomStringV1(32)
	},
})

var CaptchaStore = session.New(session.Config{
	KeyLookup:      "cookie:__SGCaptcha",
	CookieHTTPOnly: true,
	CookieSecure:   config.SecureCookies,
	CookieSameSite: config.CookieSameSite,
})

func ValidateCaptcha(c *fiber.Ctx) error {
	method := c.Method()

	if method == "GET" {
		data, err := captcha.NewMathExpr(175, 60)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": err.Error(), "message": "Generate Captcha Error"})
		}

		c.Locals("captcha", data)

		sess, err := CaptchaStore.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": err.Error(), "message": "Get Captcha Error"})
		}

		sess.Set("captcha", data.Text)

		if err := sess.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": err.Error(), "message": "Save Captcha Error"})
		}
	}

	return c.Next()
}

func Auth(c *fiber.Ctx) error {
	h := c.Get("Authorization")

	cookie := c.Cookies(CookieJWT)

	if h == "" && cookie == "" {
		return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("you don't have authorization"))
	}

	var chunks []string
	if h != "" {
		chunks = strings.Split(h, " ")
	} else if cookie != "" {
		chunks = strings.Split(cookie, " ")
	}

	if len(chunks) < 2 {
		return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("missing or malformed JWT"))
	}

	user, err := Verify(chunks[1], "access")
	if err != nil {
		return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("invalid or expired JWT"))
	}

	c.Locals("user", user)
	// c.Locals("id_account", user.IdAccount)
	// c.Locals("identity_number", user.IdentityNumber)
	// c.Locals("username", user.Username)
	// c.Locals("full_name", user.FullName)
	// c.Locals("email", user.Email)
	// c.Locals("phone_number", user.PhoneNumber)
	// c.Locals("date_of_birth", user.DateOfBirth)
	// c.Locals("auth_key", user.AuthKey)
	// c.Locals("used_pin", user.UsedPin)
	// c.Locals("is_google_account", user.IsGoogleAccount)
	// c.Locals("status_account", user.StatusAccount)
	// c.Locals("login_ip", user.LoginIp)
	// c.Locals("login_attempts", user.LoginAttempts)
	// c.Locals("login_time", user.LoginTime)
	// c.Locals("created_at", user.CreatedAt)
	// c.Locals("updated_at", user.UpdatedAt)

	return c.Next()
}

func parse(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(config.SecretKey), nil
	})
}

func Verify(token string, expectedType string) (*models.UserData, error) {
	parsed, err := parse(token)

	if err != nil || !parsed.Valid {
		// return nil, err
		return nil, errors.New("something went wrong on parse or validation")
	}

	// Parsing token claims
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		// return nil, err
		return nil, errors.New("something went wrong on claims")
	}

	typ, ok := claims["typ"].(string)
	if !ok || typ != expectedType {
		return nil, fmt.Errorf("invalid token type: expected %s, got %v", expectedType, typ)
	}

	// Getting ID, it's an interface{} so I need to cast it to uint
	// JIKA TYPE INT GANTI KE FLOAT64
	// id_account, ok := claims["id_account"].(string)
	// if !ok {
	// 	return nil, errors.New("something went wrong on id account")
	// }
	// identity_number, ok := claims["identity_number"].(string)
	// if !ok {
	// 	return nil, errors.New("something went wrong on identity number")
	// }
	// username, ok := claims["username"].(string)
	// if !ok {
	// 	return nil, errors.New("something went wrong on username")
	// }
	// full_name, ok := claims["full_name"].(string)
	// if !ok {
	// 	return nil, errors.New("something went wrong on full name")
	// }
	// email, ok := claims["email"].(string)
	// if !ok {
	// 	return nil, errors.New("something went wrong on email")
	// }
	// phone_number, ok := claims["phone_number"].(string)
	// if !ok {
	// 	return nil, errors.New("something went wrong on phone number")
	// }
	// TIPE NULL MASIH BELUM KEDETEKSI
	// date_of_birth, ok := claims["date_of_birth"].(*time.Time)
	// if !ok {
	// 	return nil, errors.New("something went wrong on date of birth")
	// }
	// auth_key, ok := claims["auth_key"].(string)
	// if !ok {
	// 	return nil, errors.New("something went wrong on auth key")
	// }
	// used_pin, ok := claims["used_pin"].(float64)
	// if !ok {
	// 	return nil, errors.New("something went wrong on used_pin")
	// }
	// is_google_account, ok := claims["is_google_account"].(float64)
	// if !ok {
	// 	return nil, errors.New("something went wrong on is google account")
	// }
	// status_account, ok := claims["status_account"].(float64)
	// if !ok {
	// 	return nil, errors.New("something went wrong on status account")
	// }
	// login_ip, ok := claims["login_ip"].(string)
	// if !ok {
	// 	return nil, errors.New("something went wrong on login_ip")
	// }
	// login_attempts, ok := claims["login_attempts"].(float64)
	// if !ok {
	// 	return nil, errors.New("something went wrong on login_attempts")
	// }
	// login_time, ok := claims["login_time"].(float64)
	// if !ok {
	// 	return nil, errors.New("something went wrong on login_time")
	// }
	// created_at, ok := claims["created_at"].(float64)
	// if !ok {
	// 	return nil, errors.New("something went wrong on created_at")
	// }
	// updated_at, ok := claims["updated_at"].(float64)
	// if !ok {
	// 	return nil, errors.New("something went wrong on updated_at")
	// }

	id_account, _ := claims["id_account"].(string)
	identity_number, _ := claims["identity_number"].(string)
	username, _ := claims["username"].(string)
	full_name, _ := claims["full_name"].(string)
	email, _ := claims["email"].(string)
	phone_number, _ := claims["phone_number"].(string)
	auth_key, _ := claims["auth_key"].(string)
	used_pin, _ := toInt8(claims["used_pin"])
	is_google_account, _ := toInt8(claims["is_google_account"])
	status_account, _ := toInt8(claims["status_account"])
	login_ip, _ := claims["login_ip"].(string)
	login_attempts, _ := toInt8(claims["login_attempts"])
	login_time, _ := toInt64(claims["login_time"])
	created_at, _ := toInt(claims["created_at"])
	updated_at, _ := toInt(claims["updated_at"])
	jti, _ := claims["jti"].(string)

	return &models.UserData{
		IdAccount:      string(id_account),
		IdentityNumber: string(identity_number),
		Username:       string(username),
		FullName:       string(full_name),
		Email:          string(email),
		PhoneNumber:    string(phone_number),
		// DateOfBirth:     date_of_birth,
		AuthKey:         string(auth_key),
		UsedPin:         int8(used_pin),
		IsGoogleAccount: int8(is_google_account),
		StatusAccount:   int8(status_account),
		LoginIp:         string(login_ip),
		LoginAttempts:   int8(login_attempts),
		LoginTime:       int64(login_time),
		CreatedAt:       int(created_at),
		UpdatedAt:       int(updated_at),
		Jti:             &jti,
	}, nil
}

func toInt8(value interface{}) (int8, error) {
	if f, ok := value.(float64); ok {
		return int8(f), nil
	}
	return 0, nil
}

func toInt64(value interface{}) (int64, error) {
	if f, ok := value.(float64); ok {
		return int64(f), nil
	}
	return 0, nil
}

func toInt(value interface{}) (int, error) {
	if f, ok := value.(float64); ok {
		return int(f), nil
	}
	return 0, nil
}

func ValidateSession(c *fiber.Ctx) error {
	currSession, err := SessionStore.Get(c)
	if err != nil {
		return handlers.UnauthorizedErrorResponse(c, err)
	}

	user := currSession.Get(UserId)
	AuthKey := fmt.Sprintf("%q", currSession.Get(AuthKey))
	defer currSession.Save()

	if user == nil {
		return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("your session is null"))
	} else {
		user_auth_by_id := fmt.Sprintf("%s-%s", UserActive, user)

		res, err := database.RedisDb.Get(handlers.Ctx, user_auth_by_id).Result()
		if err != nil {
			return c.Next() // if not implement
		}

		if AuthKey != res {
			return handlers.UnauthorizedErrorResponse(c, fmt.Errorf("you are login other device"))
		}

		return c.Next()
	}
}

func SaveUserSession(c *fiber.Ctx, user models.UserData, single_login bool) error {
	currSession, err := SessionStore.Get(c)
	if err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	if !currSession.Fresh() {
		if err := currSession.Destroy(); err != nil {
			return handlers.InternalServerErrorResponse(c, err)
		}
	}

	currSession.Set(UserId, user.IdAccount)
	currSession.Set(AuthKey, user.AuthKey)
	currSession.Set(Username, user.Username)
	currSession.Set(FullName, user.FullName)

	if err := currSession.Save(); err != nil {
		return handlers.InternalServerErrorResponse(c, err)
	}

	if single_login {
		user_auth_by_id := fmt.Sprintf("%s-%s", UserActive, user.IdAccount)
		handlers.SaveToRedis(user_auth_by_id, user.AuthKey)
	}

	return nil
}

func GenerateCaptcha(c *fiber.Ctx) error {
	c.Type("png")
	data := c.Locals("captcha").(*captcha.Data)
	output := data.WriteImage(c)

	return output
}
