package webhook

import (
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/module/webhook/models"
	"bloggo/internal/utils/apierrors"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
)

const (
	MaxRetryAttempts = 5
	BaseDelaySeconds = 2
	MaxDelaySeconds  = 60
)

type WebhookService struct {
	repository      WebhookRepository
	permissionStore permissions.Store
}

func NewWebhookService(repository WebhookRepository, permissionStore permissions.Store) WebhookService {
	return WebhookService{
		repository,
		permissionStore,
	}
}

// Config methods
func (service *WebhookService) GetConfig(roleID int64) (*models.ResponseConfig, error) {
	if !service.permissionStore.HasPermission(roleID, "webhook:manage") {
		return nil, apierrors.ErrForbidden
	}

	config, err := service.repository.GetConfig()
	if err != nil {
		return nil, err
	}

	if config == nil {
		return &models.ResponseConfig{URL: "", UpdatedAt: ""}, nil
	}

	return &models.ResponseConfig{
		URL:       config.URL,
		UpdatedAt: config.UpdatedAt,
	}, nil
}

func (service *WebhookService) UpdateConfig(url string, roleID int64) error {
	if !service.permissionStore.HasPermission(roleID, "webhook:manage") {
		return apierrors.ErrForbidden
	}

	return service.repository.UpsertConfig(url)
}

// Header methods
func (service *WebhookService) GetHeaders(roleID int64) ([]models.WebhookHeader, error) {
	if !service.permissionStore.HasPermission(roleID, "webhook:manage") {
		return nil, apierrors.ErrForbidden
	}

	return service.repository.GetAllHeaders()
}

func (service *WebhookService) BulkUpsertHeaders(items []models.RequestHeaderUpsert, roleID int64) error {
	if !service.permissionStore.HasPermission(roleID, "webhook:manage") {
		return apierrors.ErrForbidden
	}

	return service.repository.BulkUpsertHeaders(items)
}

// Request methods
func (service *WebhookService) GetRequests(limit, offset int, search string, roleID int64) ([]models.WebhookRequest, int, error) {
	if !service.permissionStore.HasPermission(roleID, "webhook:manage") {
		return nil, 0, apierrors.ErrForbidden
	}

	var requests []models.WebhookRequest
	var total int
	var err error

	if search != "" {
		requests, err = service.repository.GetRequestsBySearch(search, limit, offset)
		if err != nil {
			return nil, 0, err
		}
		total, err = service.repository.CountRequestsBySearch(search)
	} else {
		requests, err = service.repository.GetAllRequests(limit, offset)
		if err != nil {
			return nil, 0, err
		}
		total, err = service.repository.CountRequests()
	}

	if err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

func (service *WebhookService) GetRequestByID(id int, roleID int64) (*models.WebhookRequest, error) {
	if !service.permissionStore.HasPermission(roleID, "webhook:manage") {
		return nil, apierrors.ErrForbidden
	}

	return service.repository.GetRequestByID(id)
}

// Fire webhook
func (service *WebhookService) FireWebhook(payload models.WebhookPayload) {
	// Get config
	config, err := service.repository.GetConfig()
	if err != nil || config == nil || config.URL == "" {
		// No config, silently skip
		return
	}

	// Get headers
	headers, err := service.repository.GetAllHeaders()
	if err != nil {
		headers = []models.WebhookHeader{}
	}

	// Serialize payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return
	}

	// Serialize headers for storage
	headersJSON, err := json.Marshal(headers)
	if err != nil {
		headersJSON = []byte("[]")
	}
	headersStr := string(headersJSON)

	// Create initial request record
	requestRecord := &models.WebhookRequest{
		Event:          payload.Event,
		Entity:         payload.Entity,
		EntityID:       payload.ID,
		Slug:           payload.Slug,
		RequestBody:    string(payloadBytes),
		AttemptCount:   0,
		WebhookURL:     &config.URL,
		WebhookHeaders: &headersStr,
	}

	requestID, err := service.repository.InsertRequest(requestRecord)
	if err != nil {
		return
	}

	// Fire in background with retry
	go service.sendWithRetry(config.URL, headers, payloadBytes, requestID)
}

func (service *WebhookService) sendWithRetry(url string, headers []models.WebhookHeader, payload []byte, requestID int64) {
	var lastErr error
	var lastStatus *int
	var lastBody *string

	for attempt := 1; attempt <= MaxRetryAttempts; attempt++ {
		// Send request
		status, body, err := service.sendHTTPRequest(url, headers, payload)

		lastStatus = status
		lastBody = body
		lastErr = err

		// Success
		if err == nil && status != nil && *status >= 200 && *status < 300 {
			service.repository.UpdateRequest(requestID, status, body, attempt, nil)
			return
		}

		// Update with current attempt
		var errMsg *string
		if err != nil {
			msg := err.Error()
			errMsg = &msg
		}
		service.repository.UpdateRequest(requestID, status, body, attempt, errMsg)

		// If max attempts reached, stop
		if attempt >= MaxRetryAttempts {
			break
		}

		// Calculate exponential backoff delay
		delay := time.Duration(math.Min(
			float64(BaseDelaySeconds)*math.Pow(2, float64(attempt-1)),
			float64(MaxDelaySeconds),
		)) * time.Second

		time.Sleep(delay)
	}

	// Final update with last error
	var finalErrMsg *string
	if lastErr != nil {
		msg := fmt.Sprintf("Failed after %d attempts: %s", MaxRetryAttempts, lastErr.Error())
		finalErrMsg = &msg
	}
	service.repository.UpdateRequest(requestID, lastStatus, lastBody, MaxRetryAttempts, finalErrMsg)
}

func (service *WebhookService) sendHTTPRequest(url string, headers []models.WebhookHeader, payload []byte) (*int, *string, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for _, header := range headers {
		req.Header.Set(header.Key, header.Value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)

	status := resp.StatusCode
	return &status, &bodyStr, err
}

// Manual fire
func (service *WebhookService) ManualFire(roleID int64) error {
	if !service.permissionStore.HasPermission(roleID, "webhook:manage") {
		return apierrors.ErrForbidden
	}

	payload := models.WebhookPayload{
		Event:     "cms.sync",
		Entity:    "cms",
		ID:        nil,
		Slug:      nil,
		Action:    "sync",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data:      nil,
	}

	service.FireWebhook(payload)
	return nil
}
