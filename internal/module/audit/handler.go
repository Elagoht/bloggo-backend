package audit

import (
	"bloggo/internal/utils/handlers"
	"bloggo/internal/utils/pagination"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type AuditHandler struct {
	service AuditService
}

func NewAuditHandler(service AuditService) AuditHandler {
	return AuditHandler{
		service,
	}
}

// GetAuditLogs handles GET /audit-logs
func (handler *AuditHandler) GetAuditLogs(
	writer http.ResponseWriter,
	request *http.Request,
) {
	// Get pagination options
	paginationOptions, ok := pagination.GetPaginationOptions(
		writer,
		request,
		[]string{"created_at"},
	)
	if !ok {
		return
	}

	// Set defaults and calculate offset
	page := 1
	take := 20
	if paginationOptions.Page != nil {
		page = *paginationOptions.Page
	}
	if paginationOptions.Take != nil {
		take = *paginationOptions.Take
	}
	offset := (page - 1) * take

	// Get filter parameters
	query := request.URL.Query()

	// Parse comma-separated user IDs
	var userIDs []int64
	if userIDStr := query.Get("userId"); userIDStr != "" {
		for _, idStr := range strings.Split(userIDStr, ",") {
			idStr = strings.TrimSpace(idStr)
			if idStr != "" {
				if parsed, err := strconv.ParseInt(idStr, 10, 64); err == nil {
					userIDs = append(userIDs, parsed)
				}
			}
		}
	}

	// Parse comma-separated entity types
	var entityTypes []string
	if entityTypeStr := query.Get("entityType"); entityTypeStr != "" {
		for _, typeStr := range strings.Split(entityTypeStr, ",") {
			typeStr = strings.TrimSpace(typeStr)
			if typeStr != "" {
				entityTypes = append(entityTypes, typeStr)
			}
		}
	}

	// Parse comma-separated actions
	var actions []string
	if actionStr := query.Get("action"); actionStr != "" {
		for _, actStr := range strings.Split(actionStr, ",") {
			actStr = strings.TrimSpace(actStr)
			if actStr != "" {
				actions = append(actions, actStr)
			}
		}
	}

	logs, err := handler.service.GetAuditLogsWithFilters(take, offset, userIDs, entityTypes, actions)
	if err != nil {
		http.Error(writer, "Failed to retrieve audit logs", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(writer).Encode(logs)
}

// GetAuditLogsByEntity handles GET /audit-logs/entity/{type}/{id}
func (handler *AuditHandler) GetAuditLogsByEntity(
	writer http.ResponseWriter,
	request *http.Request,
) {
	entityType, ok := handlers.GetParam[string](writer, request, "type")
	if !ok {
		return
	}

	entityID, ok := handlers.GetParam[int64](writer, request, "id")
	if !ok {
		return
	}

	// Get pagination options
	paginationOptions, ok := pagination.GetPaginationOptions(
		writer,
		request,
		[]string{"created_at"},
	)
	if !ok {
		return
	}

	// Set defaults and calculate offset
	page := 1
	take := 20
	if paginationOptions.Page != nil {
		page = *paginationOptions.Page
	}
	if paginationOptions.Take != nil {
		take = *paginationOptions.Take
	}
	offset := (page - 1) * take

	logs, err := handler.service.GetAuditLogsByEntity(entityType, entityID, *paginationOptions.Take, offset)
	if err != nil {
		http.Error(writer, "Failed to retrieve audit logs", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(writer).Encode(logs)
}

// GetAuditLogsByUser handles GET /audit-logs/user/{id}
func (handler *AuditHandler) GetAuditLogsByUser(
	writer http.ResponseWriter,
	request *http.Request,
) {
	userID, ok := handlers.GetParam[int64](writer, request, "id")
	if !ok {
		return
	}

	// Get pagination options
	paginationOptions, ok := pagination.GetPaginationOptions(
		writer,
		request,
		[]string{"created_at"},
	)
	if !ok {
		return
	}

	// Set defaults and calculate offset
	page := 1
	take := 20
	if paginationOptions.Page != nil {
		page = *paginationOptions.Page
	}
	if paginationOptions.Take != nil {
		take = *paginationOptions.Take
	}
	offset := (page - 1) * take

	logs, err := handler.service.GetAuditLogsByUser(userID, *paginationOptions.Take, offset)
	if err != nil {
		http.Error(writer, "Failed to retrieve audit logs", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(writer).Encode(logs)
}
