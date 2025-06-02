package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// FOR WINDOWS
var BasePath = "."
var PathEnv = BasePath + "env"

// FOR LINUX
// var BasePath = "/home/siasura1/golangapp/" //DOMAINESIA
// var PathEnv = BasePath + ".env"

func Config(key string) string {
	if err := godotenv.Load(PathEnv); err != nil {
		log.Printf("Error loading .env file or not define env:%s", key)
	}
	return os.Getenv(key)
}

var ConfigPath = Config("CONFIG_PATH_JSON")

var APP_PORT = Config("PORT_APIGATEWAY")

var AuthTimeCache = 86400 * time.Second
var CsrfTimeCache = 3600 * time.Second

var RedisTimeCache = 10800 * time.Second

var SecretKey = Config("SECRET_KEY_APIGATEWAY")
var CookieKey = Config("KEY_COOKIE_APIGATEWAY")

var SecureCookiesString = Config("SECURE_COOKIES")
var SecureCookies bool
var CookieSameSite = Config("COOKIES_SAME_SITE")

// DEV
// var SecureCookies = false  //change true to prod false to dev
// var CookieSameSite = "Lax" //change None to prod Lax to dev

// PROD
// var SecureCookies = true    //change true to prod false to dev
// var CookieSameSite = "None" //change None to prod Lax to dev

// SIKA REPOSITORY
var SikaRepoURL = Config("SIKA_REPO_URL")
