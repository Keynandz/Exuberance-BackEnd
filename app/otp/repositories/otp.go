package otp

import (
	"math/rand"
	"time"

	"gorm.io/gorm"

	"exuberance-backend/models"
	"exuberance-backend/pkg/database"
)

type OtpRepository struct{}

func NewOtpRepository() *OtpRepository {
	return &OtpRepository{}
}

func (repo *OtpRepository) GenerateOTP() int32 {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	min := 100000
	max := 999999
	return int32(rng.Intn(max-min+1) + min)
}

func (repo *OtpRepository) CreateOtp(userID int32, otp int32) error {
	db := database.GetDB()

	newOtp := models.MasterOtp{
		OtpCode:   otp,
		Createdby: userID,
		Status:    "unverif",
	}

	return db.Create(&newOtp).Error
}

func (repo *OtpRepository) UpdateOtp(userID int32) error {
	db := database.GetDB()

	updatedOtp := repo.GenerateOTP()

	return db.Transaction(func(tx *gorm.DB) error {
		var existingOtp models.MasterOtp
		if err := tx.Where("createdby = ?", userID).First(&existingOtp).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return repo.CreateOtp(userID, updatedOtp)
			}
			return err
		}

		return tx.Model(&existingOtp).Update("otp_code", updatedOtp).Error
	})
}

func (repo *OtpRepository) UpdateOrCreateOtp(userID int32, otp int32) error {
	db := database.GetDB()

	return db.Transaction(func(tx *gorm.DB) error {
		var existingOtp models.MasterOtp
		if err := tx.Where("createdby = ?", userID).First(&existingOtp).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return repo.CreateOtp(userID, otp)
			}
			return err
		}

		return tx.Model(&existingOtp).Update("otp_code", otp).Error
	})
}

func (repo *OtpRepository) DeleteExpiredOtps() error {
	db := database.GetDB()

	expirationTime := time.Now().Add(-5 * time.Minute)
	err := db.Where("created < ?", expirationTime).Delete(&models.MasterOtp{}).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return err
}

func (repo *OtpRepository) IsUserVerified(userID int32) (bool, error) {
	db := database.GetDB()

	var user models.DataUser
	if err := db.Select("isactive").Where("id = ?", userID).First(&user).Error; err != nil {
		return false, err
	}

	return user.Isactive == 1, nil
}

func (repo *OtpRepository) GetUserEmailByID(userID int32) (string, error) {
	db := database.GetDB()

	var user models.DataUser
	if err := db.Select("email").Where("id = ?", userID).First(&user).Error; err != nil {
		return "", err
	}

	return user.Email, nil
}
