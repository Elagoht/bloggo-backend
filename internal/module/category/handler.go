package category

import (
	"bloggo/internal/module/category/models"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/handlers"
	"encoding/json"
	"net/http"
)

type CategoryHandler struct {
	service CategoryService
}

func NewCategoryHandler(service CategoryService) CategoryHandler {
	return CategoryHandler{
		service,
	}
}

func (handler *CategoryHandler) CategoryCreate(
	writer http.ResponseWriter,
	request *http.Request,
) {
	body, ok := handlers.BindAndValidate[models.RequestCategoryCreate](writer, request)
	if !ok {
		return
	}

	response, err := handler.service.CategoryCreate(&body)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(response)
}
