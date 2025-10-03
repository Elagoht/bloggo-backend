package models

// ResponseConfig represents the webhook config response
type ResponseConfig struct {
	URL       string `json:"url"`
	UpdatedAt string `json:"updatedAt"`
}

// ResponseMessage is a generic message response
type ResponseMessage struct {
	Message string `json:"message"`
}
