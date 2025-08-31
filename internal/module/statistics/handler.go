package statistics

import (
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/handlers"
	"encoding/json"
	"net/http"
	"strconv"
)

type StatisticsHandler struct {
	service StatisticsService
}

func NewStatisticsHandler(service StatisticsService) StatisticsHandler {
	return StatisticsHandler{
		service,
	}
}

func (handler *StatisticsHandler) GetAllStatistics(
	writer http.ResponseWriter,
	request *http.Request,
) {
	roleId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenRoleId)
	if !ok {
		return
	}

	userId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenUserId)
	if !ok {
		return
	}

	// Check if userId query parameter is provided
	userIdParam := request.URL.Query().Get("userId")

	if userIdParam == "" {
		// No userId param - return all statistics
		response, err := handler.service.GetAllStatistics(roleId)
		if err != nil {
			apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
				apierrors.ErrForbidden: {
					Message: "You need permissions to view statistics.",
					Status:  http.StatusForbidden,
				},
			})
			return
		}
		json.NewEncoder(writer).Encode(response)
	} else {
		// userId param provided - return user-specific statistics
		targetUserId, err := strconv.ParseInt(userIdParam, 10, 64)
		if err != nil {
			apierrors.MapErrors(apierrors.ErrBadRequest, writer, apierrors.HTTPErrorMapping{
				apierrors.ErrBadRequest: {
					Message: "Invalid userId parameter.",
					Status:  http.StatusBadRequest,
				},
			})
			return
		}

		response, err := handler.service.GetAuthorStatistics(targetUserId, roleId, userId)
		if err != nil {
			apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
				apierrors.ErrForbidden: {
					Message: "You can only view your own statistics or you need admin/editor permissions to view others'.",
					Status:  http.StatusForbidden,
				},
			})
			return
		}
		json.NewEncoder(writer).Encode(response)
	}
}

func (handler *StatisticsHandler) GetUserOwnStatistics(
	writer http.ResponseWriter,
	request *http.Request,
) {
	roleId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenRoleId)
	if !ok {
		return
	}

	userId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenUserId)
	if !ok {
		return
	}

	response, err := handler.service.GetUserOwnStatistics(roleId, userId)
	if err != nil {
		apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
			apierrors.ErrForbidden: {
				Message: "You need permission to view your own statistics.",
				Status:  http.StatusForbidden,
			},
		})
		return
	}

	json.NewEncoder(writer).Encode(response)
}

func (handler *StatisticsHandler) GetAuthorStatistics(
	writer http.ResponseWriter,
	request *http.Request,
) {
	roleId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenRoleId)
	if !ok {
		return
	}

	userId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenUserId)
	if !ok {
		return
	}

	authorIdStr, ok := handlers.GetParam[string](writer, request, "authorId")
	if !ok {
		return
	}

	authorId, err := strconv.ParseInt(authorIdStr, 10, 64)
	if err != nil {
		apierrors.MapErrors(apierrors.ErrBadRequest, writer, apierrors.HTTPErrorMapping{
			apierrors.ErrBadRequest: {
				Message: "Invalid author ID.",
				Status:  http.StatusBadRequest,
			},
		})
		return
	}

	response, err := handler.service.GetAuthorStatistics(authorId, roleId, userId)
	if err != nil {
		apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
			apierrors.ErrForbidden: {
				Message: "You can only view your own statistics or you need editor/admin permissions.",
				Status:  http.StatusForbidden,
			},
		})
		return
	}

	json.NewEncoder(writer).Encode(response)
}
