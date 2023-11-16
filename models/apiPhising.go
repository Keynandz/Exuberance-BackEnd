package models

type ScanRequest struct {
	APIKey   string `json:"apiKey"`
	URLInfo  URLInfo `json:"urlInfo"`
	ScanType string `json:"scanType"`
}

type URLInfo struct {
	URL string `json:"url"`
}

type ScanResponse struct {
	JobID     string `json:"jobID"`
	Timestamp int64  `json:"timestamp"`
}

type StatusRequest struct {
	APIKey   string `json:"apiKey"`
	JobID    string `json:"jobID"`
	Insights bool   `json:"insights"`
}

type StatusResponse struct {
	JobID           string        `json:"job_id"`
	Status          string        `json:"status"`
	URL             string        `json:"url"`
	URLSHA256       string        `json:"url_sha256"`
	Disposition     string        `json:"disposition"`
	Brand           string        `json:"brand"`
	Insights        string        `json:"insights"`
	Resolved        bool          `json:"resolved"`
	ScreenshotPath  string        `json:"screenshot_path"`
	ScanStartTS     int64         `json:"scan_start_ts"`
	ScanEndTS       int64         `json:"scan_end_ts"`
	Error           bool          `json:"error"`
	ImageObjects    []ImageObject `json:"image_objects"`
}

type ImageObject struct {
	ObjectType string  `json:"object_type"`
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	Width      float64 `json:"width"`
	Height     float64 `json:"height"`
	Text       string  `json:"text"`
}