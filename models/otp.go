package models

import (
	"time"

	_ "gorm.io/gorm"
)

type MasterOtp struct {
	Id        int32     `gorm:"primarykey" json:"id"`
	Created   time.Time `gorm:"default:current_timestamp" json:"created"`
	Createdby int32     `gorm:"unique" json:"createdby"`
	OtpCode   int32     `json:"otp_code"`
	Status    string    `gorm:"type:varchar(255)" json:"status"`
}