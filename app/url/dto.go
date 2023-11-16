package url

import (
	"exuberance-backend/models"
	"time"
)

type GetUrl struct {
	Id        int       `json:"id"`
	Created   time.Time `json:"created"`
	Url       string    `json:"url"`
	IpAddress string    `json:"ip_address"`
	Status    string    `json:"status"`
}

func NewGetUrl(url models.LogUrl) GetUrl {

	status := ""
	if url.Status == "clean" {
		status = "Clean Websites"
	} else if url.Status == "mortgage" {
		status = "Clean Websites"
	} else if url.Status == "hacked_website" {
		status = "Hacked websites"
	} else if url.Status == "streaming" {
		status = "Illegal Streaming website"
	} else if url.Status == "cryptojacking" {
		status = "Cryptojacking"
	} else if url.Status == "likely_phish" {
		status = "Likely Phish websites"
	} else if url.Status == "suspicious" {
		status = "Suspicious websites"
	} else if url.Status == "gambling" {
		status = "Gambling websites"
	} else if url.Status == "drug_spam" {
		status = "Drug Spam"
	} else if url.Status == "adult" {
		status = "Adult websitse"
	} else if url.Status == "scam" {
		status = "Tech support scams"
	} else if url.Status == "phish" {
		status = "Phishing Websites"
	}

	return GetUrl{
		Id:        url.Id,
		Created:   url.Created,
		Url:       url.Url,
		IpAddress: url.IpAddress,
		Status:    status,
	}
}
