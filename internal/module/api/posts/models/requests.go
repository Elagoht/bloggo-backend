package models

// Request body for tracking post view
type APITrackViewRequest struct {
	UserAgent string `json:"userAgent" validate:"required,max=500"`
}
