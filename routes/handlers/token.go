package handlers

import (
	"net/http"
	"strings"

	token "exuberance-backend/app/token/repositories"

	"github.com/labstack/echo/v4"
)

func TokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error":       "Token is missing",
				"status_code": http.StatusUnauthorized,
			})
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		newToken, err := token.VerifyToken(tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error":       "Invalid token",
				"status_code": http.StatusUnauthorized,
			})
		}

		if newToken != tokenString {
			c.Set("newToken", newToken)
		}

		return next(c)
	}
}