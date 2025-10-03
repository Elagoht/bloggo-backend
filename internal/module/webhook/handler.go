package webhook

import (
	"bloggo/internal/module/audit"
	"bloggo/internal/module/webhook/models"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/handlers"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type WebhookHandler struct {
	service WebhookService
}

func NewWebhookHandler(service WebhookService) WebhookHandler {
	return WebhookHandler{
		service,
	}
}

// GetConfig handles GET /api/webhook/config
func (handler *WebhookHandler) GetConfig(
	writer http.ResponseWriter,
	request *http.Request,
) {
	roleID, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenRoleId)
	if !ok {
		return
	}

	config, err := handler.service.GetConfig(roleID)
	if err != nil {
		apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
			apierrors.ErrForbidden: {
				Message: "Only editors and admins can manage webhooks.",
				Status:  http.StatusForbidden,
			},
		})
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(config)
}

// UpdateConfig handles PUT /api/webhook/config
func (handler *WebhookHandler) UpdateConfig(
	writer http.ResponseWriter,
	request *http.Request,
) {
	roleID, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenRoleId)
	if !ok {
		return
	}

	userID, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenUserId)
	if !ok {
		return
	}

	body, ok := handlers.BindAndValidate[models.RequestUpdateConfig](writer, request)
	if !ok {
		return
	}

	err := handler.service.UpdateConfig(body.URL, roleID)
	if err != nil {
		apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
			apierrors.ErrForbidden: {
				Message: "Only editors and admins can manage webhooks.",
				Status:  http.StatusForbidden,
			},
		})
		return
	}

	// Log audit
	audit.LogWebhookAction(&userID, "config_updated")

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(models.ResponseMessage{
		Message: "Webhook config updated successfully",
	})
}

// GetHeaders handles GET /api/webhook/headers
func (handler *WebhookHandler) GetHeaders(
	writer http.ResponseWriter,
	request *http.Request,
) {
	roleID, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenRoleId)
	if !ok {
		return
	}

	headers, err := handler.service.GetHeaders(roleID)
	if err != nil {
		apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
			apierrors.ErrForbidden: {
				Message: "Only editors and admins can manage webhooks.",
				Status:  http.StatusForbidden,
			},
		})
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(headers)
}

// BulkUpsertHeaders handles PUT /api/webhook/headers
func (handler *WebhookHandler) BulkUpsertHeaders(
	writer http.ResponseWriter,
	request *http.Request,
) {
	roleID, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenRoleId)
	if !ok {
		return
	}

	userID, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenUserId)
	if !ok {
		return
	}

	body, ok := handlers.BindAndValidate[models.RequestBulkUpsertHeaders](writer, request)
	if !ok {
		return
	}

	err := handler.service.BulkUpsertHeaders(body.Items, roleID)
	if err != nil {
		apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
			apierrors.ErrForbidden: {
				Message: "Only editors and admins can manage webhooks.",
				Status:  http.StatusForbidden,
			},
		})
		return
	}

	// Log audit
	audit.LogWebhookAction(&userID, "headers_updated")

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(models.ResponseMessage{
		Message: "Webhook headers updated successfully",
	})
}

// ManualFire handles POST /api/webhook/fire
func (handler *WebhookHandler) ManualFire(
	writer http.ResponseWriter,
	request *http.Request,
) {
	roleID, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenRoleId)
	if !ok {
		return
	}

	userID, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenUserId)
	if !ok {
		return
	}

	err := handler.service.ManualFire(roleID)
	if err != nil {
		apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
			apierrors.ErrForbidden: {
				Message: "Only editors and admins can manage webhooks.",
				Status:  http.StatusForbidden,
			},
		})
		return
	}

	// Log audit
	audit.LogWebhookAction(&userID, "manual_fire")

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(models.ResponseMessage{
		Message: "Webhook fired successfully",
	})
}

// GetRequests handles GET /api/webhook/requests
func (handler *WebhookHandler) GetRequests(
	writer http.ResponseWriter,
	request *http.Request,
) {
	roleID, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenRoleId)
	if !ok {
		return
	}

	// Parse query params
	limitStr := request.URL.Query().Get("limit")
	offsetStr := request.URL.Query().Get("offset")
	search := request.URL.Query().Get("search")

	limit := 20
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	requests, total, err := handler.service.GetRequests(limit, offset, search, roleID)
	if err != nil {
		apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
			apierrors.ErrForbidden: {
				Message: "Only editors and admins can view webhook requests.",
				Status:  http.StatusForbidden,
			},
		})
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(map[string]interface{}{
		"data":  requests,
		"total": total,
		"limit": limit,
		"offset": offset,
	})
}

// GetRequestByID handles GET /api/webhook/requests/:id
func (handler *WebhookHandler) GetRequestByID(
	writer http.ResponseWriter,
	request *http.Request,
) {
	roleID, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenRoleId)
	if !ok {
		return
	}

	idStr := chi.URLParam(request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(map[string]string{
			"error": "Invalid request ID",
		})
		return
	}

	req, err := handler.service.GetRequestByID(id, roleID)
	if err != nil {
		apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
			apierrors.ErrForbidden: {
				Message: "Only editors and admins can view webhook requests.",
				Status:  http.StatusForbidden,
			},
		})
		return
	}

	if req == nil {
		writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(writer).Encode(map[string]string{
			"error": "Request not found",
		})
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(req)
}
