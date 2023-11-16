package repositories

import (
	dto "exuberance-backend/app/auditlog"
	"exuberance-backend/models"
	"exuberance-backend/pkg/database"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func GetIp(c echo.Context) (string, error) {
	r := c.Request()

	ips := r.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")

	if len(splitIps) > 0 {
		netIP := net.ParseIP(splitIps[len(splitIps)-1])
		if netIP != nil {
			return netIP.String(), nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	netIP := net.ParseIP(ip)
	if netIP != nil {
		ip := netIP.String()
		if ip == "::1" {
			return "127.0.0.1", nil
		}
		return ip, nil
	}

	return "", fmt.Errorf("ip not found")
}

func CreateAuditLog(user *models.DataUser, token, ip string) (models.AuditLog, error) {
	db := database.GetDB()

	postAudit := models.AuditLog{
		Created:       time.Now(),
		UserId:        user.Id,
		RemoteAddress: ip,
		SessionId:     token,
	}

	result := db.Create(&postAudit)
	if result.Error != nil {
		return postAudit, fmt.Errorf("error creating audit log: %w", result.Error)
	}

	return postAudit, nil
}

func GetAuditLog() ([]dto.GetAuditLog, error) {
	db := database.GetDB()

	var log []models.AuditLog
	result := db.Find(&log)
	if result.Error != nil {
		return nil, fmt.Errorf("error retrieving log: %w", result.Error)
	}

	getLog := make([]dto.GetAuditLog, len(log))
	for i, log := range log {
		getLog[i] = dto.NewGetAuditLog(log)
	}

	return getLog, nil
}
