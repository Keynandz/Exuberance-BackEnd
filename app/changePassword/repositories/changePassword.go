package repositories

import (
	dto "exuberance-backend/app/changePassword"
	"exuberance-backend/models"
	"exuberance-backend/pkg/database"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func CheckPassowrd(id int) (*models.DataUser, error) {
	db := database.GetDB()
	akun := &models.DataUser{}

	result := db.First(&akun, id)
	if result.Error != nil {
		return nil, fmt.Errorf("error fetching account: %w", result.Error)
	}

	return akun, nil
}

func UpdateNewPassword(pass dto.NewPassword, id int) error {
	db := database.GetDB()

	var user models.DataUser
	result := db.First(&user, id)
	if result.Error != nil {
		return fmt.Errorf("error retrieving user: %w", result.Error)
	}

	if pass.OldPassword == "" {
		return fmt.Errorf("old password is required")
	}

	if pass.NewPassword == "" {
		return fmt.Errorf("new password is required")
	}

	if err := EncryptNewPass(&pass); err != nil {
		return fmt.Errorf("password is required")
	}

	user.Password = pass.NewPassword

	update := db.Save(&user)
	if update.Error != nil {
		return fmt.Errorf("failed to update user: %w", update.Error)
	}

	return nil
}

func EncryptNewPass(pass *dto.NewPassword) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	pass.NewPassword = string(hashedPassword)
	return nil
}
