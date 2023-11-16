package models

import (
	"time"

	_ "gorm.io/gorm"
)

type DataUser struct {
	Id           int32     `gorm:"primarykey" json:"id"`
	Created      time.Time `gorm:"current_timestamp" json:"created"`
	Isactive     int32     `json:"isactive"`
	Name         string    `gorm:"type:varchar" json:"name"`
	Email        string    `gorm:"type:varchar" json:"email"`
	Password     string    `gorm:"type:varchar" json:"password"`
	OauthUid     string    `gorm:"type:varchar" json:"oauth_uid"`
	ImageProfile string    `gorm:"type:varchar" json:"image_profile"`
	AuditLog     AuditLog  `gorm:"foreignKey:UserId; constraint:OnDelete:SET NULL;"`
}
