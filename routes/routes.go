package routes

import (
	"exuberance-backend/routes/handlers"

	"github.com/labstack/echo/v4"
)

func WebRoutes(e *echo.Echo) {
	// Check Connection
	e.GET("/connection/web/v1/", handlers.Home)

	// Middleware
	authGroup := e.Group("")
	authGroup.Use(handlers.TokenMiddleware)

	// Register and Login
	e.POST("/connection/exuberance/login", handlers.LoginAccount)
	e.POST("/connection/exuberance/register", handlers.CreateAccount)

	// Change Password
	authGroup.POST("/connection/exuberance/user/changepass/:id", handlers.ChangePassword)

	// Google
	e.GET("/connection/exuberance/google/callback", handlers.GoogleCallbackWeb)

	// OTP
	e.POST("/connection/exuberance/register/generateotp", handlers.GenerateOtp)
	e.POST("/connection/exuberance/register/resendotp", handlers.ResendOtp)
	e.POST("/connection/exuberance/register/verifotp", handlers.VerifyOtp)

	// Audit Log
	authGroup.GET("/connection/exuberance/auditlog", handlers.GetAuditLog)

	// Phising Checker
	authGroup.POST("/connection/exuberance/url/checker", handlers.ScanHandler)
	authGroup.GET("/connection/exuberance/url/checked", handlers.GetUrl)
}
