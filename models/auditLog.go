package models

import (
	"time"

	_ "gorm.io/gorm"
)

type AuditLog struct {
	Id            int32     `gorm:"primarykey" json:"id"`
	Created       time.Time `gorm:"default:current_timestamp" json:"created"`
	UserId        int32     `json:"user_id"`
	RemoteAddress string    `gorm:"type:varchar(255)" json:"remote_address"`
	SessionId     string    `gorm:"type:varchar" json:"session_id"`
}
