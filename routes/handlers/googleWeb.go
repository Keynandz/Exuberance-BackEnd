package handlers

import (
	"context"
	"encoding/json"
	audit "exuberance-backend/app/auditlog/repositories"
	"exuberance-backend/app/login/repositories"
	register "exuberance-backend/app/register/repositories"
	jwt "exuberance-backend/app/token/repositories"
	"exuberance-backend/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GoogleCallbackWeb(c echo.Context) error {
	loadErr := godotenv.Load()
	if loadErr != nil {
		log.Fatal("error loading file .env")
	}

	formatWeb := "http://%s:%s"
	formatRedirect := "http://%s:%s/connection/exuberance/google/callback"

	redirectUrl := fmt.Sprintf(formatRedirect, os.Getenv("SERVER_ADDRESS"), os.Getenv("SERVER_PORT"))
	webUrl := fmt.Sprintf(formatWeb, os.Getenv("WEB_HOST"), os.Getenv("WEB_PORT"))

	var oauthConf = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
		RedirectURL:  redirectUrl,
	}

	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message":     "code is required",
			"status_code": http.StatusBadRequest,
		})
	}

	token, err := oauthConf.Exchange(context.Background(), code)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message":     err.Error(),
			"status_code": http.StatusBadRequest,
		})
	}

	client := oauthConf.Client(oauth2.NoContext, token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to fetch user info")
	}
	defer response.Body.Close()

	var googleUserInfo struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Id      string `json:"id"`
		Picture string `json:"picture"`
	}

	err = json.NewDecoder(response.Body).Decode(&googleUserInfo)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to decode user info")
	}

	email := googleUserInfo.Email

	user, err := repositories.GetUserByEmailGoogle(email)
	if err != nil {

		reuser := models.DataUser{}
		reuser.Isactive = 1
		result, _ := register.CreateUserByGoogle(googleUserInfo, reuser)

		userID := uint(result.Id)
		jwt, err := jwt.GenerateToken(userID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse("Failed to generate token"))
		}

		if c.Get("newToken") != nil {
			newToken, ok := c.Get("newToken").(string)
			if ok && newToken != jwt {
				jwt = newToken
			}
		}

		ip, err := audit.GetIp(c)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		_, err = audit.CreateAuditLog(&result, jwt, ip)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		cookie := http.Cookie{
			Name:   "token",
			Value:  jwt,
			Path:   "/",
			MaxAge: 86400,
		}
		http.SetCookie(c.Response().Writer, &cookie)

		id := int(result.Id)
		resultId := strconv.Itoa(id)

		cookies := http.Cookie{
			Name:   "id",
			Value:  resultId,
			Path:   "/",
			MaxAge: 86400,
		}
		http.SetCookie(c.Response().Writer, &cookies)

		cookiess := http.Cookie{
			Name:   "oauth",
			Value:  result.OauthUid,
			Path:   "/",
			MaxAge: 86400,
		}
		http.SetCookie(c.Response().Writer, &cookiess)

		profiles := http.Cookie{
			Name:   "profile",
			Value:  result.ImageProfile,
			Path:   "/",
			MaxAge: 86400,
		}
		http.SetCookie(c.Response().Writer, &profiles)

		return c.Redirect(http.StatusFound, webUrl)
	}

	if googleUserInfo.Id != user.OauthUid {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"response": map[string]interface{}{
				"message":     "Invalid oauth id, try again!",
				"status_code": http.StatusUnauthorized,
			},
		})
	}

	userID := uint(user.Id)
	jwt, err := jwt.GenerateToken(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("Failed to generate token"))
	}

	if c.Get("newToken") != nil {
		newToken, ok := c.Get("newToken").(string)
		if ok && newToken != jwt {
			jwt = newToken
		}
	}

	ip, err := audit.GetIp(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	_, err = audit.CreateAuditLog(user, jwt, ip)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	cookie := http.Cookie{
		Name:   "token",
		Value:  jwt,
		Path:   "/",
		MaxAge: 86400,
	}
	http.SetCookie(c.Response().Writer, &cookie)

	id := int(user.Id)
	resultId := strconv.Itoa(id)

	cookies := http.Cookie{
		Name:   "id",
		Value:  resultId,
		Path:   "/",
		MaxAge: 86400,
	}
	http.SetCookie(c.Response().Writer, &cookies)

	cookiess := http.Cookie{
		Name:   "oauth",
		Value:  user.OauthUid,
		Path:   "/",
		MaxAge: 86400,
	}
	http.SetCookie(c.Response().Writer, &cookiess)

	profiles := http.Cookie{
		Name:   "profile",
		Value:  user.ImageProfile,
		Path:   "/",
		MaxAge: 86400,
	}
	http.SetCookie(c.Response().Writer, &profiles)

	return c.Redirect(http.StatusFound, webUrl)
}
