package repositories

import (
	dto "exuberance-backend/app/url"
	"exuberance-backend/models"
	"exuberance-backend/pkg/database"
	"fmt"
	"time"
)

func SaveUrl(url models.StatusResponse) error {
	db := database.GetDB()

	postUrl := models.LogUrl{
		Created: time.Now(),
		Url:     url.URL,
		Status:  url.Disposition,
	}

	result := db.Create(&postUrl)
	if result.Error != nil {
		return fmt.Errorf("error creating audit log: %w", result.Error)
	}

	return nil
}

func GetUrl() ([]dto.GetUrl, error) {
	db := database.GetDB()

	var url []models.LogUrl
	result := db.Find(&url)
	if result.Error != nil {
		return nil, fmt.Errorf("error retrieving url: %w", result.Error)
	}

	getUrl := make([]dto.GetUrl, len(url))
	for i, url := range url {
		getUrl[i] = dto.NewGetUrl(url)
	}

	return getUrl, nil
}
