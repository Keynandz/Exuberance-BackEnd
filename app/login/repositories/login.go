package repositories

import (
	"exuberance-backend/models"
	"exuberance-backend/pkg/database"
	"fmt"
)

func GetUserByEmail(email string) (*models.DataUser, error) {
	db := database.GetDB()
	akun := &models.DataUser{}

	result := db.Where("email = ?", email).First(akun)
	if result.Error != nil {
		return nil, fmt.Errorf("error fetching akun: %w", result.Error)
	}

	return akun, nil
}

func GetUserByEmailGoogle(email string) (*models.DataUser, error) {
	db := database.GetDB()
	akun := &models.DataUser{}
	result := db.Where("email = ?", email).First(&akun)
	if result.Error != nil {
        return nil, fmt.Errorf("error fetching account: %w", result.Error)
    }

	if akun.OauthUid != "" {
		return akun, nil
	}

	return nil, fmt.Errorf("registered but not with google")
}
