package audit

import (
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/handlers"
	"bloggo/internal/utils/pagination"
	"encoding/json"
	"net/http"
)

type AuditHandler struct {
	service AuditService
}

func NewAuditHandler(service AuditService) AuditHandler {
	return AuditHandler{
		service,
	}
}

func (handler *AuditHandler) GetAuditLogs(
	writer http.ResponseWriter,
	request *http.Request,
) {
	roleId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenRoleId)
	if !ok {
		return
	}

	paginate, ok := pagination.GetPaginationOptions(writer, request, []string{
		"created_at",
	})
	if !ok {
		return
	}

	auditLogs, err := handler.service.GetAuditLogs(paginate, roleId)
	if err != nil {
		apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
			apierrors.ErrForbidden: {
				Message: "Only admins can view audit logs.",
				Status:  http.StatusForbidden,
			},
		})
		return
	}

	json.NewEncoder(writer).Encode(auditLogs)
}