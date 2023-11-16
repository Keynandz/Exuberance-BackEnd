package main

import (
	"exuberance-backend/pkg/database"
	"exuberance-backend/pkg/log"
	"exuberance-backend/routes"
	"exuberance-backend/routes/handlers"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	database.InitDB()
	database.Migrate()

	go handlers.ScheduleDeleteExpiredOtps()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "X-Auth-Token"},
	}))
	
	routes.WebRoutes(e)

	e.Use(log.LogRequest)
	e.Logger.Fatal(e.Start(":8800"))
}
