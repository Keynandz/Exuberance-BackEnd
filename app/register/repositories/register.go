package register

import (
	"exuberance-backend/models"
	"exuberance-backend/pkg/database"
	"fmt"
	"time"
)

func CreateUser(user models.DataUser) (models.DataUser, error) {
	db := database.GetDB()

	var google models.DataUser
	db.Where("email = ? AND password IS NULL", user.Email).First(&google)

	if google.OauthUid != "" {
		data := models.DataUser{
			Password: user.Password,
		}

		if err := db.Save(&data).Error; err != nil {
			return user, fmt.Errorf("error creating user from Google data: %w", err)
		}

		return user, nil
	}

	user.Created = time.Now()

	result := db.Create(&user)
	if result.Error != nil {
		return user, fmt.Errorf("error creating user: %w", result.Error)
	}

	return user, nil
}

func CreateUserByGoogle(googleUserInfo struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Id      string `json:"id"`
	Picture string `json:"picture"`
}, user models.DataUser) (models.DataUser, error) {
	db := database.GetDB()

	var google models.DataUser
	result := db.Where("email = ?", googleUserInfo.Email).First(&google)
	if result.Error == nil {

		google.OauthUid = googleUserInfo.Id
		google.ImageProfile = googleUserInfo.Picture
		google.Isactive = 1

		resultUpdate := db.Save(&google)
		if resultUpdate.Error != nil {
			return google, fmt.Errorf("error updating data: %w", resultUpdate.Error)
		}

		return google, nil
	}

	data := models.DataUser{
		Name:         googleUserInfo.Name,
		Email:        googleUserInfo.Email,
		OauthUid:     googleUserInfo.Id,
		Isactive:     user.Isactive,
		ImageProfile: googleUserInfo.Picture,
	}

	if err := db.Create(&data).Error; err != nil {
		return data, fmt.Errorf("error creating user from Google data: %w", err)
	}

	return data, nil
}

func DefaultPicture(id int, defaultImageData string) error {
	db := database.GetDB()

	var user models.DataUser
	err := db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	user.ImageProfile = defaultImageData

	updateResult := db.Save(&user)
	if updateResult.Error != nil {
		return fmt.Errorf("failed to upload image: %w", updateResult.Error)
	}

	return nil
}
