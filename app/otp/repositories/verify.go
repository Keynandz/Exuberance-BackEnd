package otp

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"exuberance-backend/models"
	"exuberance-backend/pkg/database"
)

type OtpVerify struct{}

func NewOtpVerify() *OtpVerify {
	return &OtpVerify{}
}

func (verify *OtpVerify) VerifyOtp(userID int32, otpCode int32) (bool, error) {
	db := database.GetDB()

	var existingOtp models.MasterOtp
	if err := db.Where("createdby = ?", userID).First(&existingOtp).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("OTP record not found for user:", userID)
			return false, nil
		}
		fmt.Println("Error fetching OTP record:", err)
		return false, err
	}

	if existingOtp.OtpCode == otpCode {
		expirationTime := time.Now().Add(-5 * time.Minute)
		if existingOtp.Created.Before(expirationTime) {
			fmt.Println("OTP expired for user:", userID)
			return false, nil
		}
		fmt.Println("OTP verified for user:", userID)
		return true, nil
	}

	fmt.Println("OTP does not match for user:", userID)
	return false, nil
}

func UpdateUserIsactive(userID int32) error {
	db := database.GetDB()

	result := db.Model(&models.DataUser{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"isactive": 1,
		})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func UpdateOtp(userID int32) error {
	db := database.GetDB()

	result := db.Model(&models.MasterOtp{}).
		Where("createdby = ?", userID).
		Update("status", "Verified")
	if result.Error != nil {
		return result.Error
	}

	return nil
}
