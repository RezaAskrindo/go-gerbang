package middleware

import (
	"errors"
	"fmt"

	"go-gerbang/config"
	"go-gerbang/handlers"
	"go-gerbang/models"

	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/golang-jwt/jwt"
	"github.com/steambap/captcha"
)

// DOMAINESIA NOT SUPPORT
// var StorageRedis = redis.New(redis.Config{
// 	URL: config.Config("REDIS_ADDRESS_FULL"),
// })

var SessionStore = session.New(session.Config{
	KeyLookup:      "cookie:__SecureGatewayS",
	CookieHTTPOnly: true,
	CookieSecure:   config.SecureCookies,
	CookieSameSite: config.CookieSameSite,
	KeyGenerator: func() string {
		return handlers.RandomString(8)
	},
	// Storage: StorageRedis,
})

var CsrfActivated = false

var CsrfStore = session.New(session.Config{
	Expiration:     config.CsrfTimeCache,       // Expire sessions after 30 minutes of inactivity
	KeyLookup:      "cookie:__SecureGatewayC3", // Recommended to use the __Host- prefix when serving the app over TLS
	CookieSecure:   config.SecureCookies,
	CookieHTTPOnly: true,
	CookieSameSite: "Lax",
})

var CsrfProtection = csrf.New(csrf.Config{
	Session: CsrfStore,
	Next: func(c *fiber.Ctx) bool {
		return CsrfActivated
	},
	KeyLookup:      "cookie:__SecureGatewayC",
	CookieName:     "__SecureGatewayC",
	CookieHTTPOnly: true,
	CookieSameSite: "Lax",
	Expiration:     config.CsrfTimeCache,
	KeyGenerator: func() string {
		return handlers.RandomString(32)
	},
	ContextKey: "token_csrf",
})

var CaptchaStore = session.New(session.Config{
	KeyLookup:      "cookie:__SecureGatewayC2",
	CookieHTTPOnly: true,
	CookieSecure:   config.SecureCookies,
	CookieSameSite: config.CookieSameSite,
	KeyGenerator: func() string {
		return handlers.RandomString(8)
	},
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

	cookie := c.Cookies("__SecureGatewayJ")

	if h == "" && cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": false, "message": "You Don't Have Authorization"})
	}

	var chunks []string
	// Spliting the header
	if h != "" {
		chunks = strings.Split(h, " ")
	} else if cookie != "" {
		chunks = strings.Split(cookie, " ")
	}

	// If header signature is not like `Bearer <token>`, then throw
	// This is also required, otherwise chunks[1] will throw out of bound error
	if len(chunks) < 2 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": false, "message": "Missing or malformed JWT", "data": nil})
	}

	// Verify the token which is in the chunks
	user, err := Verify(chunks[1])

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": false, "message": err.Error()})
		// return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": false, "message": "Invalid or expired JWT"})
	}

	if err := RoleAccessChecking(c, user); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Your role don't have authorized"})
		// return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	c.Locals("id_account", user.IdAccount)
	c.Locals("identity_number", user.IdentityNumber)
	c.Locals("username", user.Username)
	c.Locals("full_name", user.FullName)
	c.Locals("email", user.Email)
	c.Locals("phone_number", user.PhoneNumber)
	c.Locals("date_of_birth", user.DateOfBirth)
	c.Locals("auth_key", user.AuthKey)
	c.Locals("used_pin", user.UsedPin)
	c.Locals("is_google_account", user.IsGoogleAccount)
	c.Locals("status_account", user.StatusAccount)
	c.Locals("login_ip", user.LoginIp)
	c.Locals("login_attempts", user.LoginAttempts)
	c.Locals("login_time", user.LoginTime)
	c.Locals("created_at", user.CreatedAt)
	c.Locals("updated_at", user.UpdatedAt)

	return c.Next()
}

func parse(token string) (*jwt.Token, error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(config.SecretKey), nil
	})
}

func Verify(token string) (*models.UserData, error) {
	parsed, err := parse(token)

	if err != nil {
		return nil, err
	}

	// Parsing token claims
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}

	// Getting ID, it's an interface{} so I need to cast it to uint
	// JIKA TYPE INT GANTI KE FLOAT64
	id_account, ok := claims["id_account"].(string)
	if !ok {
		return nil, errors.New("Something went wrong on id account")
	}
	identity_number, ok := claims["identity_number"].(string)
	if !ok {
		return nil, errors.New("Something went wrong on identity number")
	}
	username, ok := claims["username"].(string)
	if !ok {
		return nil, errors.New("Something went wrong on username")
	}
	full_name, ok := claims["full_name"].(string)
	if !ok {
		return nil, errors.New("Something went wrong on full name")
	}
	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("Something went wrong on email")
	}
	phone_number, ok := claims["phone_number"].(string)
	if !ok {
		return nil, errors.New("Something went wrong on phone number")
	}
	// TIPE NULL MASIH BELUM KEDETEKSI
	// date_of_birth, ok := claims["date_of_birth"].(*time.Time)
	// if !ok {
	// 	return nil, errors.New("Something went wrong on date of birth")
	// }
	auth_key, ok := claims["auth_key"].(string)
	if !ok {
		return nil, errors.New("Something went wrong on auth key")
	}
	used_pin, ok := claims["used_pin"].(float64)
	if !ok {
		return nil, errors.New("Something went wrong on used_pin")
	}
	is_google_account, ok := claims["is_google_account"].(float64)
	if !ok {
		return nil, errors.New("Something went wrong on is google account")
	}
	status_account, ok := claims["status_account"].(float64)
	if !ok {
		return nil, errors.New("Something went wrong on status account")
	}
	login_ip, ok := claims["login_ip"].(string)
	if !ok {
		return nil, errors.New("Something went wrong on login_ip")
	}
	login_attempts, ok := claims["login_attempts"].(float64)
	if !ok {
		return nil, errors.New("Something went wrong on login_attempts")
	}
	login_time, ok := claims["login_time"].(float64)
	if !ok {
		return nil, errors.New("Something went wrong on login_time")
	}
	created_at, ok := claims["created_at"].(float64)
	if !ok {
		return nil, errors.New("Something went wrong on created_at")
	}
	updated_at, ok := claims["updated_at"].(float64)
	if !ok {
		return nil, errors.New("Something went wrong on updated_at")
	}

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
	}, nil
}

func RequireLogin(c *fiber.Ctx) error {
	currSession, err := SessionStore.Get(c)
	if err != nil {
		return err
	}
	user := currSession.Get("UserID")
	AuthKey := fmt.Sprintf("%q", currSession.Get("auth_key"))
	defer currSession.Save()

	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "You're Session is Nil"})
	} else {
		user_auth_by_id := fmt.Sprintf("sika-user-%d", user)

		res, err := handlers.Cache.Get(handlers.Ctx, user_auth_by_id).Result()
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "You're Session is Invalid"})
		}

		if AuthKey != res {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "You're Login Other Device"})
		}

		return c.Next()
	}
}

func SaveUserSession(user_data models.UserData, c *fiber.Ctx) error {
	currSession, err := SessionStore.Get(c)

	defer currSession.Save()
	if err != nil {
		return err
	}
	err = currSession.Regenerate()
	if err != nil {
		return err
	}
	currSession.Set("User", fiber.Map{"AuthKey": user_data.AuthKey})

	go handlers.GenerateAuthkeyToRedis(user_data.IdAccount, user_data.AuthKey)

	return nil
}

func GenerateCaptcha(c *fiber.Ctx) error {
	c.Type("png")
	data := c.Locals("captcha").(*captcha.Data)
	output := data.WriteImage(c)

	return output
}

func RoleAccessChecking(c *fiber.Ctx, user *models.UserData) error {
	// RbacConfigs := &types.AuthRoleRouteResponses{}
	// val, err := handlers.Cache.Get(handlers.Ctx, "sika-auth-role-route").Bytes()
	// if err != nil {
	// 	rows, err := database.DBNOORM.Query("SELECT role_access, route, method FROM sika_auth_role_route ORDER BY route")
	// 	if err != nil {
	// 		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	// 	}
	// 	defer rows.Close()

	// 	for rows.Next() {
	// 		d := types.AuthRoleRoute{}
	// 		if err := rows.Scan(&d.RoleAccess, &d.Route, &d.Method); err != nil {
	// 			return err // Exit if we get an error
	// 		}

	// 		RbacConfigs.Items = append(RbacConfigs.Items, d)
	// 	}

	// 	data := handlers.ToMarshal(RbacConfigs)

	// 	cacheErr := handlers.Cache.Set(handlers.Ctx, "sika-auth-role-route", data, config.RedisTimeCache).Err()
	// 	if cacheErr != nil {
	// 		return cacheErr
	// 	}
	// } else {
	// 	if err := json.Unmarshal(val, RbacConfigs); err != nil {
	// 		return err
	// 	}
	// }
	// // val, _ := handlers.Cache.Get(handlers.Ctx, "sika-auth-role-route").Bytes()

	// originalURL := strings.ToLower(c.Path())

	// var RoleAccess []int = nil
	// var HadAccess bool

	// for _, auth_role := range RbacConfigs.Items {
	// 	if auth_role.Route == originalURL && auth_role.Method == c.Method() {
	// 		RoleAccess = append(RoleAccess, auth_role.RoleAccess)
	// 	}
	// }

	// HadAccess = true
	// if len(RoleAccess) > 0 {
	// 	HadAccess = false
	// 	for _, role := range RoleAccess {
	// 		if user.StatusRole == role {
	// 			HadAccess = true
	// 		}
	// 	}
	// }

	// if HadAccess != true {
	// 	return errors.New("You're No Authorized To Access This Route")
	// }

	return nil
}
