package handlers

import (
	"exuberance-backend/app/auditlog/repositories"
	"net/http"
	"sort"

	"github.com/labstack/echo/v4"
)

func GetAuditLog(c echo.Context) error {
	auditLog, err := repositories.GetAuditLog()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	sort.Slice(auditLog, func(i, j int) bool {
		return auditLog[i].Id > auditLog[j].Id
	})

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":        auditLog,
		"status_code": http.StatusOK,
	})
}
