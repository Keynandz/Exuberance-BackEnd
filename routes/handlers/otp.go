package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	otp "exuberance-backend/app/otp/repositories"
	"exuberance-backend/models"

	"github.com/labstack/echo/v4"
)

func GenerateOtp(c echo.Context) error {
	var user models.MasterOtp
	if err := c.Bind(&user); err != nil {
		return err
	}

	repo := otp.NewOtpRepository()

	isVerified, verifyErr := repo.IsUserVerified(user.Createdby)
	if verifyErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to check user verification status",
		})
	}

	if isVerified {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "User is already verified",
		})
	}

	email, err := repo.GetUserEmailByID(user.Createdby)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to retrieve user's email, try to resend otp!",
		})
	}

	createOtp := repo.GenerateOTP()
	createErr := repo.CreateOtp(user.Createdby, createOtp)
	if createErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":       "Failed to generate and store OTP, Try to resend otp!",
			"status_code": http.StatusInternalServerError,
		})
	}

	createErr = otp.Otp(email, strconv.Itoa(int(createOtp)))
	if createErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to send OTP email",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":     "OTP sent successfully",
		"status_code": http.StatusOK,
	})
}

func ResendOtp(c echo.Context) error {
	var user models.MasterOtp
	if err := c.Bind(&user); err != nil {
		return err
	}

	repo := otp.NewOtpRepository()

	email, err := repo.GetUserEmailByID(user.Createdby)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to retrieve user's email",
		})
	}

	isVerified, err := repo.IsUserVerified(user.Createdby)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to check user verification status",
		})
	}

	if isVerified {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":     "User is already verified",
			"status_code": http.StatusOK,
		})
	}

	updatedOtp := repo.GenerateOTP()
	createErr := repo.UpdateOrCreateOtp(user.Createdby, updatedOtp)
	if createErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":       "Failed to resend OTP",
			"status_code": http.StatusInternalServerError,
		})
	}

	createErr = otp.Otp(email, strconv.Itoa(int(updatedOtp)))
	if createErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to send OTP email",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":     "OTP resent successfully",
		"status_code": http.StatusOK,
	})
}

func ScheduleDeleteExpiredOtps() {
	for {
		repo := otp.NewOtpRepository()
		err := repo.DeleteExpiredOtps()
		if err != nil {
			fmt.Println("Error deleting expired OTPs:", err)
		}
		time.Sleep(5 * time.Minute)
	}
}
