package models

import "time"

type StatusRequest struct {
	Websites []string `json:"websites"`
}

type WebsiteStatus struct {
	URL         string
	Status      string
	LastChecked time.Time
}

type StatusResponse struct {
	StatusArray []WebsiteStatus
}
