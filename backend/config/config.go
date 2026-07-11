package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// FOR WINDOWS
var BasePath = "."
var PathEnv = BasePath + "env"

// FOR LINUX
// var BasePath = "/home/siasura1/golangapp/" //DOMAINESIA
// var BasePath = "/var/www/services/go-gerbang/" //BIZNET
// var PathEnv = BasePath + ".env"

func Config(key string) string {
	// if err := godotenv.Load(PathEnv); err != nil {
	// For Linux we don't use PathEnv
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file or not define env:%s", key)
	}
	return os.Getenv(key)
}

var ConfigPath = Config("CONFIG_PATH_JSON")

var APP_PORT = Config("PORT_APIGATEWAY")

var AuthTimeCache = 1 * time.Hour
var RefreshAuthTimeCache = 24 * 14 * time.Hour
var CsrfTimeCache = 1 * time.Second

var RedisTimeCache = 12 * time.Hour

var SecretKey = Config("SECRET_KEY_APIGATEWAY")
var CookieKey = Config("KEY_COOKIE_APIGATEWAY")

// var SecureCookies bool
var SecureCookiesString = Config("SECURE_COOKIES")
var CookieSameSite = Config("COOKIES_SAME_SITE")

// DEV
var SecureCookies = false //change true to prod false to dev

// PROD
// var SecureCookies = true    //change true to prod false to dev

// SIKA REPOSITORY
// var SikaRepoURL = Config("SIKA_REPO_URL")

func GetTrustedOrigins() []string {
	originsStr := Config("TRUSTED_ORIGINS")
	if originsStr == "" {
		return []string{} // or default origins
	}

	origins := strings.Split(originsStr, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}
	return origins
}
