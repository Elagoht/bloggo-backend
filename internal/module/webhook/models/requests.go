package models

// RequestUpdateConfig represents the request to update webhook config
type RequestUpdateConfig struct {
	URL string `json:"url" validate:"required,url,max=2048"`
}

// RequestHeaderUpsert represents a single header upsert
type RequestHeaderUpsert struct {
	Key   string `json:"key" validate:"required,max=255"`
	Value string `json:"value" validate:"required"`
}

// RequestBulkUpsertHeaders represents bulk header updates
type RequestBulkUpsertHeaders struct {
	Items []RequestHeaderUpsert `json:"items" validate:"required,dive"`
}

// RequestGetWebhookRequests represents query params for getting webhook requests
type RequestGetWebhookRequests struct {
	Search string `json:"search"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}
