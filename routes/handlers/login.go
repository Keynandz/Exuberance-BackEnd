package handlers

import (
	"context"
	audit "exuberance-backend/app/auditlog/repositories"
	"exuberance-backend/app/login/repositories"
	token "exuberance-backend/app/token/repositories"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/crypto/bcrypt"
)

func LoginAccount(c echo.Context) error {
	loginData := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	if err := c.Bind(&loginData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"response": map[string]interface{}{
				"message":     "Invalid request payload",
				"status_code": http.StatusBadRequest,
			},
		})
	}

	user, err := repositories.GetUserByEmail(loginData.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"response": map[string]interface{}{
				"message":     "Your account is not registered, please register first.",
				"status_code": http.StatusUnauthorized,
			},
		})
	}

	if user.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"response": map[string]interface{}{
				"message":     "Your account is registed by Google, try login with Google.",
				"status_code": http.StatusBadRequest,
			},
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"response": map[string]interface{}{
				"message":     "Invalid password, try again!",
				"status_code": http.StatusUnauthorized,
			},
		})
	}

	if user.Isactive == 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"response": map[string]interface{}{
				"message":     "Your account is not verified. Please verify with otp",
				"status_code": http.StatusBadRequest,
			},
		})
	}

	userID := uint(user.Id)
	token, err := token.GenerateToken(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"response": map[string]interface{}{
				"message":     "Failed to generate token",
				"status_code": http.StatusInternalServerError,
			},
		})
	}

	if c.Get("newToken") != nil {
		newToken, ok := c.Get("newToken").(string)
		if ok && newToken != token {
			token = newToken
		}
	}

	ip, err := audit.GetIp(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	logAudit, err := audit.CreateAuditLog(user, token, ip)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	loadErr := godotenv.Load()
	if loadErr != nil {
		log.Fatal("error loading file .env")
	}

	ssl, _ := strconv.ParseBool(os.Getenv("MINIO_SSL"))
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := ssl
	bucketName := os.Getenv("MINIO_BUCKET")

	minioClient, _ := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	fileName := user.ImageProfile
	reqParams := make(url.Values)
	object, _ := minioClient.PresignedGetObject(context.Background(), bucketName, fileName, time.Hour*168, reqParams)

	imageProfileURL := object.String()

	response := map[string]interface{}{
		"status_code": http.StatusOK,
		"token":       token,
		"data": map[string]interface{}{
			"message":  "Login successful",
			"Id":       user.Id,
			"email":    user.Email,
			"name":     user.Name,
			"isactive": user.Isactive,
		},
		"image_profile": imageProfileURL,
		"audit":         logAudit.Id,
	}

	return c.JSON(http.StatusOK, SuccessResponse(response))
}

type errorResponse struct {
	Error string `json:"error"`
}

type successResponse struct {
	Response interface{} `json:"response"`
}

func ErrorResponse(errorMessage string) errorResponse {
	return errorResponse{Error: errorMessage}
}

func SuccessResponse(data interface{}) successResponse {
	return successResponse{Response: data}
}
