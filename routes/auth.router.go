package routes

import (
	"go-gerbang/middleware"
	"go-gerbang/services"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/auth")

	auth.Get("/get-google-client-id", services.LoadGoogleLoginClienId)

	auth.Post("/less-secure/login", middleware.ValidateCaptcha, services.Login)
	auth.Post("/less-secure/refresh-token", middleware.ValidateCaptcha, services.RefreshAuth)
	auth.Post("/less-secure/request-reset-password", middleware.ValidateCaptcha, services.RequestResetPassword)
	auth.Post("/less-secure/reset-password", middleware.ValidateCaptcha, services.ResetPassword)
	auth.Post("/less-secure/change-password", middleware.ValidateCaptcha, services.ChangePassword)
	auth.Post("/less-secure/sign-up", services.Signup)

	auth.Get("/logout", services.LogoutWeb)
	auth.Post("/login", middleware.ValidateCaptcha, middleware.CsrfProtection, services.Login)
	auth.Post("/refresh-token", middleware.ValidateCaptcha, middleware.CsrfProtection, services.RefreshAuth)
	auth.Post("/login-with-google", middleware.ValidateCaptcha, middleware.CsrfProtection, services.LoginWithGoogle)
	auth.Post("/request-reset-password", middleware.ValidateCaptcha, middleware.CsrfProtection, services.RequestResetPassword)
	auth.Post("/reset-password", middleware.ValidateCaptcha, middleware.CsrfProtection, services.ResetPassword)
	auth.Post("/change-password", middleware.ValidateCaptcha, middleware.CsrfProtection, services.ChangePassword)
	auth.Post("/sign-up", middleware.CsrfProtection, services.Signup)

	auth.Get("/get-captcha", middleware.ValidateCaptcha, middleware.GenerateCaptcha)
	auth.Get("/get-session", middleware.ValidateSession, services.GetSession)
	auth.Get("/get-jwt-info", middleware.Auth, services.GetSessionJWT)

	// auth.Post("/with-google", services.LoginWithGoogle)
	// authSession := auth.Group("/v1/auth") //.Use(middleware.Auth)
	// authSession.Get("/auth-key/:token", services.AuthByJWT)

	usersApi := app.Group("/users")
	usersApi.Get("/all", services.GetAllUser)
	usersApi.Get("/:userId", services.FindUserById)
	usersApi.Post("/", middleware.ValidateCaptcha, middleware.CsrfProtection, services.CreateUser)
	usersApi.Put("/:userId", middleware.ValidateCaptcha, middleware.CsrfProtection, services.UpdateUser)
	usersApi.Delete("/:userId", middleware.ValidateCaptcha, middleware.CsrfProtection, services.DeleteUser)
}
