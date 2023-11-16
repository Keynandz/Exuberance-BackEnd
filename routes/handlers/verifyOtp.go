package handlers

import (
	otp "exuberance-backend/app/otp/repositories"
	"exuberance-backend/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func VerifyOtp(c echo.Context) error {
	var user models.MasterOtp
	if err := c.Bind(&user); err != nil {
		return err
	}

	otpRepo := otp.NewOtpVerify()

	otpVerificationResult, err := otpRepo.VerifyOtp(user.Createdby, user.OtpCode)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Error verifying OTP",
		})
	}

	if otpVerificationResult {
		err := otp.UpdateUserIsactive(user.Createdby)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "Error updating user activation",
			})
		}
		otp.UpdateOtp(user.Createdby)
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":     "OTP verification successful. User activated.",
			"status_code": http.StatusOK,
		})
	}

	return c.JSON(http.StatusUnauthorized, map[string]interface{}{
		"error": "Invalid OTP",
	})
}
