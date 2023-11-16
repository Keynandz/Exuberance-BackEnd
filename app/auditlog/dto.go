package auditlog

import (
	"exuberance-backend/models"
	"exuberance-backend/pkg/database"
	"time"
)

type PostAuditLog struct {
	Created       time.Time `json:"created"`
	UserId        int       `json:"user_Id"`
	RemoteAddress string    `json:"remote_address"`
	SessionId     string    `json:"session_id"`
}

type GetAuditLog struct {
	Id            int32  `json:"id"`
	Created       string `json:"created"`
	UserId        int32  `json:"user_id"`
	Name          string `json:"name"`
	RemoteAddress string `json:"remote_address"`
}

func NewGetAuditLog(log models.AuditLog) GetAuditLog {
	db := database.GetDB()

	var user models.DataUser
	db.First(&user, log.UserId)

	return GetAuditLog{
		Id:            log.Id,
		Created:       log.Created.Format("2006-01-02 15:04:05"),
		UserId:        log.UserId,
		Name:          user.Email,
		RemoteAddress: log.RemoteAddress,
	}
}
