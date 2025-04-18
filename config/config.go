package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var BasePath = "."
var PathEnv = BasePath + "env"

// FOR LINUX USE FULL PATH
// var BasePath = "/home/adminfra/sika/"

// FOR DOMAINESIA
// var BasePath = "/home/siasura1/golangapp/"
// var PathEnv = BasePath + ".env"

func Config(key string) string {
	if err := godotenv.Load(PathEnv); err != nil {
		log.Printf("Error loading .env file or not define env:%s", key)
		// panic("Error loading .env file")
	}
	return os.Getenv(key)
}

var ConfigPath = Config("CONFIG_PATH_JSON")

// DOMAINESIA 8060 IT 8070
// var APP_PORT = ":8060"
var APP_PORT = Config("PORT_APIGATEWAY")

var AuthTimeCache = 86400 * time.Second
var CsrfTimeCache = 3600 * time.Second

var RedisTimeCache = 10800 * time.Second

var SecretKey = Config("SECRET_KEY_APIGATEWAY")
var CookieKey = Config("KEY_COOKIE_APIGATEWAY")

// var SecretKey = "QkIXbHKbQtIU80zSbiLDeGueLcwTh9X3"
// var CookieKey = "FBsMP5oHv3EZS74jW1XnOVJRDmecc9F8"

// DEV
var SecureCookies = false  //change true to prod false to dev
var CookieSameSite = "Lax" //change None to prod Lax to dev

// PROD
// var SecureCookies = true    //change true to prod false to dev
// var CookieSameSite = "None" //change None to prod Lax to dev

// SIKA REPOSITORY
var SikaRepoURL = Config("SIKA_REPO_URL")
