package handlers

import (
	otp "exuberance-backend/app/otp/repositories"
	register "exuberance-backend/app/register/repositories"
	jwt "exuberance-backend/app/token/repositories"
	"exuberance-backend/models"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func CreateAccount(c echo.Context) error {
	user := models.DataUser{}
	otpRepo := otp.NewOtpRepository()

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, ErrResponse("Invalid request payload"))
	}

	user.Isactive = 0

	if err := EncryptPassword(&user); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	createdUser, err := register.CreateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	go func() {

		otpCode := otpRepo.GenerateOTP()
		err = otpRepo.CreateOtp(createdUser.Id, otpCode)
		otpCodeStr := strconv.Itoa(int(otpCode))
		err = otp.Otp(createdUser.Email, otpCodeStr)

	}()

	userID := uint(user.Id)
	token, err := jwt.GenerateToken(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("Failed to generate token"))
	}

	response := map[string]interface{}{
		"status_code": http.StatusOK,
		"token":       token,
		"data": map[string]interface{}{
			"message":  "Register successful. An OTP has been sent to your email.",
			"id":       createdUser.Id,
			"email":    createdUser.Email,
			"isactive": createdUser.Isactive,
		},
	}

	response["status_code"] = http.StatusOK

	userRegis := uint(createdUser.Id)

	err = DefaultPicture(userRegis)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, SuccResponse(response))
}

func EncryptPassword(user *models.DataUser) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}

type errResponse struct {
	Error string `json:"error"`
}

type succResponse struct {
	Response interface{} `json:"response"`
}

func ErrResponse(errMessage string) errResponse {
	return errResponse{Error: errMessage}
}

func SuccResponse(data interface{}) succResponse {
	return succResponse{Response: data}
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
