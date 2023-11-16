package repositories

import (
	dto "exuberance-backend/app/url"
	"exuberance-backend/models"
	"exuberance-backend/pkg/database"
	"fmt"
	"net"
	"net/url"
	"time"
)

func SaveUrl(url models.StatusResponse) error {
	db := database.GetDB()

	ipAddress, err := getIPAddress(url.URL)
	if err != nil {
		return fmt.Errorf("error getting IP address: %w", err)
	}

	postUrl := models.LogUrl{
		Created:   time.Now(),
		Url:       url.URL,
		IpAddress: ipAddress,
		Status:    url.Disposition,
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

func getIPAddress(urlString string) (string, error) {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return "", fmt.Errorf("error parsing URL: %w", err)
	}

	ips, err := net.LookupIP(parsedURL.Hostname())
	if err != nil {
		return "", fmt.Errorf("error looking up IP: %w", err)
	}

	if len(ips) > 0 {
		return ips[0].String(), nil
	}

	return "", fmt.Errorf("no IP address found for URL")
}
