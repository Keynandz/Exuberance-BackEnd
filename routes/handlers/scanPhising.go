package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"exuberance-backend/app/url/repositories"
	"exuberance-backend/models"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func ScanHandler(c echo.Context) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("API_KEY")
	scanType := os.Getenv("SCAN_TYPE")
	insights := os.Getenv("INSIGHTS")

	var scanRequest models.ScanRequest
	if err := c.Bind(&scanRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	scanRequest.APIKey = apiKey
	scanRequest.ScanType = scanType

	checkphishURL := "https://developers.checkphish.ai/api/neo/scan"
	requestBody, err := json.Marshal(scanRequest)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	resp, err := http.Post(checkphishURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to make request"})
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read response"})
	}

	var scanResponse models.ScanResponse
	if err := json.Unmarshal(responseBody, &scanResponse); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse response"})
	}

	statusResponse, err := sendStatusRequest(apiKey, scanResponse.JobID, insights)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get scan status"})
	}

	if statusResponse.Status == "" {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":      "Server has reached daily limit, please try again later",
			"status_code": http.StatusInternalServerError,
		})
	}

	for statusResponse.Status == "PENDING" {
		time.Sleep(5 * time.Second)

		statusResponse, err = sendStatusRequest(apiKey, scanResponse.JobID, insights)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"status":      "Failed to get scan status",
				"status_code": http.StatusInternalServerError,
			})
		}
	}

	repositories.SaveUrl(statusResponse)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"url":         statusResponse.URL,
		"status":      statusResponse.Disposition,
		"status_code": http.StatusOK,
	})

}

func sendStatusRequest(apiKey, jobID, insights string) (models.StatusResponse, error) {
	statusRequest := models.StatusRequest{
		APIKey:   apiKey,
		JobID:    jobID,
		Insights: insights == "true",
	}

	checkphishStatusURL := "https://developers.checkphish.ai/api/neo/scan/status"
	statusRequestBody, err := json.Marshal(statusRequest)
	if err != nil {
		return models.StatusResponse{}, err
	}

	resp, err := http.Post(checkphishStatusURL, "application/json", bytes.NewBuffer(statusRequestBody))
	if err != nil {
		return models.StatusResponse{}, err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.StatusResponse{}, err
	}

	var statusResponse models.StatusResponse
	if err := json.Unmarshal(responseBody, &statusResponse); err != nil {
		return models.StatusResponse{}, err
	}

	return statusResponse, nil
}

func GetUrl(c echo.Context) error {
	auditLog, err := repositories.GetUrl()
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
