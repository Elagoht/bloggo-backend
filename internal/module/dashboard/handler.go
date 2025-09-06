package dashboard

import (
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/handlers"
	"encoding/json"
	"net/http"
)

type DashboardHandler struct {
	service DashboardService
}

func NewDashboardHandler(service DashboardService) DashboardHandler {
	return DashboardHandler{
		service,
	}
}

func (handler *DashboardHandler) GetDashboardStats(
	writer http.ResponseWriter,
	request *http.Request,
) {
	roleId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenRoleId)
	if !ok {
		return
	}

	stats, err := handler.service.GetDashboardStats(roleId)
	if err != nil {
		apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
			apierrors.ErrForbidden: {
				Message: "You don't have permission to view dashboard statistics.",
				Status:  http.StatusForbidden,
			},
		})
		return
	}

	json.NewEncoder(writer).Encode(stats)
}