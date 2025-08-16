package removal_request

import (
	"bloggo/internal/module/removal_request/models"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/handlers"
	"encoding/json"
	"net/http"
)

type RemovalRequestHandler struct {
	service RemovalRequestService
}

func NewRemovalRequestHandler(service RemovalRequestService) RemovalRequestHandler {
	return RemovalRequestHandler{
		service,
	}
}

func (handler *RemovalRequestHandler) CreateRemovalRequest(
	writer http.ResponseWriter,
	request *http.Request,
) {
	userId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenUserId)
	if !ok {
		return
	}

	body, ok := handlers.BindAndValidate[*models.RequestCreateRemovalRequest](
		writer,
		request,
	)
	if !ok {
		return
	}

	var note *string
	if body.Note != "" {
		note = &body.Note
	}

	createdId, err := handler.service.CreateRemovalRequest(
		body.PostVersionId,
		userId,
		note,
	)
	if err != nil {
		apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
			apierrors.ErrForbidden: {
				Message: "You can only request removal of your own published or approved versions.",
				Status:  http.StatusForbidden,
			},
			apierrors.ErrPreconditionFailed: {
				Message: "Only published or approved versions can be requested for removal.",
				Status:  http.StatusPreconditionFailed,
			},
			apierrors.ErrConflict: {
				Message: "You already have a pending removal request for this version.",
				Status:  http.StatusConflict,
			},
		})
		return
	}

	json.NewEncoder(writer).Encode(createdId)
}

func (handler *RemovalRequestHandler) GetRemovalRequestList(
	writer http.ResponseWriter,
	request *http.Request,
) {
	roleId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenRoleId)
	if !ok {
		return
	}

	requests, err := handler.service.GetRemovalRequestList(roleId)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	json.NewEncoder(writer).Encode(requests)
}

func (handler *RemovalRequestHandler) GetRemovalRequestById(
	writer http.ResponseWriter,
	request *http.Request,
) {
	userId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenUserId)
	if !ok {
		return
	}

	roleId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenRoleId)
	if !ok {
		return
	}

	id, ok := handlers.GetParam[int64](writer, request, "id")
	if !ok {
		return
	}

	removalRequest, err := handler.service.GetRemovalRequestById(id, userId, roleId)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	json.NewEncoder(writer).Encode(removalRequest)
}

func (handler *RemovalRequestHandler) GetUserRemovalRequests(
	writer http.ResponseWriter,
	request *http.Request,
) {
	userId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenUserId)
	if !ok {
		return
	}

	requests, err := handler.service.GetUserRemovalRequests(userId)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	json.NewEncoder(writer).Encode(requests)
}

func (handler *RemovalRequestHandler) ApproveRemovalRequest(
	writer http.ResponseWriter,
	request *http.Request,
) {
	userId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenUserId)
	if !ok {
		return
	}

	roleId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenRoleId)
	if !ok {
		return
	}

	id, ok := handlers.GetParam[int64](writer, request, "id")
	if !ok {
		return
	}

	if err := handler.service.ApproveRemovalRequest(id, userId, roleId); err != nil {
		apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
			apierrors.ErrForbidden: {
				Message: "You don't have permission to approve removal requests.",
				Status:  http.StatusForbidden,
			},
			apierrors.ErrPreconditionFailed: {
				Message: "Only pending removal requests can be approved.",
				Status:  http.StatusPreconditionFailed,
			},
		})
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func (handler *RemovalRequestHandler) RejectRemovalRequest(
	writer http.ResponseWriter,
	request *http.Request,
) {
	userId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenUserId)
	if !ok {
		return
	}

	roleId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenRoleId)
	if !ok {
		return
	}

	id, ok := handlers.GetParam[int64](writer, request, "id")
	if !ok {
		return
	}

	if err := handler.service.RejectRemovalRequest(id, userId, roleId); err != nil {
		apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
			apierrors.ErrForbidden: {
				Message: "You don't have permission to reject removal requests.",
				Status:  http.StatusForbidden,
			},
			apierrors.ErrPreconditionFailed: {
				Message: "Only pending removal requests can be rejected.",
				Status:  http.StatusPreconditionFailed,
			},
		})
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
