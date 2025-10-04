package models

// WebhookConfig represents the webhook configuration
type WebhookConfig struct {
	ID        int    `json:"id"`
	URL       string `json:"url"`
	UpdatedAt string `json:"updatedAt"`
}

// WebhookHeader represents a custom HTTP header
type WebhookHeader struct {
	ID        int    `json:"id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// WebhookRequest represents a webhook request log
type WebhookRequest struct {
	ID             int     `json:"id"`
	Event          string  `json:"event"`
	Entity         string  `json:"entity"`
	EntityID       *int64  `json:"entityId"`
	Slug           *string `json:"slug"`
	RequestBody    string  `json:"requestBody"`
	ResponseStatus *int    `json:"responseStatus"`
	ResponseBody   *string `json:"responseBody"`
	AttemptCount   int     `json:"attemptCount"`
	ErrorMessage   *string `json:"errorMessage"`
	WebhookURL     *string `json:"webhookUrl"`
	WebhookHeaders *string `json:"webhookHeaders"`
	CreatedAt      string  `json:"createdAt"`
}

// WebhookPayload represents the payload sent to webhook
type WebhookPayload struct {
	Event     string                 `json:"event"`
	Entity    string                 `json:"entity"`
	ID        *int64                 `json:"id"`
	Slug      *string                `json:"slug"`
	OldSlug   *string                `json:"oldSlug,omitempty"`
	Action    string                 `json:"action"`
	Timestamp string                 `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}
