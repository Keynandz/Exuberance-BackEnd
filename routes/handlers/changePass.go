package handlers

import (
	dto "exuberance-backend/app/changePassword"
	"exuberance-backend/app/changePassword/repositories"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func ChangePassword(c echo.Context) error {
	user := c.Param("id")
	id, _ := strconv.Atoi(user)

	updatePass := dto.NewPassword{}
	if err := c.Bind(&updatePass); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	pass, _ := repositories.CheckPassowrd(id)

	if pass.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"response": map[string]interface{}{
				"message":     "Your account is registed by Google, try login with Google.",
				"status_code": http.StatusBadRequest,
			},
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(pass.Password), []byte(updatePass.OldPassword)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"response": map[string]interface{}{
				"message":     "Invalid old password, try again!",
				"status_code": http.StatusUnauthorized,
			},
		})
	}

	err := repositories.UpdateNewPassword(updatePass, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":     "Password Has Been Updated Successfully.",
		"status_code": http.StatusOK,
	})
}
