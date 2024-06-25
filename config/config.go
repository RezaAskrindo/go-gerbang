package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var BasePath = ""
var ConfigPath = "./config/config.json"

// FOR LINUX USE FULL PATH
// var BasePath = "/home/adminfra/sika/"

// FOR DOMAINESIA
// var BasePath = "/home/siasura1/golangapp/"
// var ConfigPath = "config/config.json"

var PathEnv = BasePath + "../.env"

func Config(key string) string {
	err := godotenv.Load(PathEnv)

	if err != nil {
		fmt.Print("Error loading .env file")
	}
	return os.Getenv(key)
}

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
// var SecureCookies = false  //change true to prod false to dev
// var CookieSameSite = "Lax" //change None to prod Lax to dev

// PROD
var SecureCookies = true    //change true to prod false to dev
var CookieSameSite = "None" //change None to prod Lax to dev

// SIKA REPOSITORY
var SikaRepoURL = Config("SIKA_REPO_URL")
