package keyvalue

import (
	"bloggo/internal/module/keyvalue/models"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/handlers"
	"encoding/json"
	"net/http"
)

type KeyValueHandler struct {
	service KeyValueService
}

func NewKeyValueHandler(service KeyValueService) KeyValueHandler {
	return KeyValueHandler{
		service,
	}
}

func (handler *KeyValueHandler) GetAll(
	writer http.ResponseWriter,
	request *http.Request,
) {
	roleId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenRoleId)
	if !ok {
		return
	}

	items, err := handler.service.GetAll(roleId)
	if err != nil {
		apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
			apierrors.ErrForbidden: {
				Message: "Only editors and admins can manage key-values.",
				Status:  http.StatusForbidden,
			},
		})
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(items)
}

func (handler *KeyValueHandler) BulkUpsert(
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

	body, ok := handlers.BindAndValidate[models.RequestKeyValueBulkUpsert](writer, request)
	if !ok {
		return
	}

	err := handler.service.BulkUpsert(body.Items, roleId, userId)
	if err != nil {
		apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
			apierrors.ErrForbidden: {
				Message: "Only editors and admins can manage key-values.",
				Status:  http.StatusForbidden,
			},
		})
		return
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(map[string]string{
		"message": "Key-value pairs saved successfully",
	})
}
