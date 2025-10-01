package categories

import (
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/handlers"
	"encoding/json"
	"net/http"
)

type CategoriesAPIHandler struct {
	service CategoriesAPIService
}

func NewCategoriesAPIHandler(service CategoriesAPIService) CategoriesAPIHandler {
	return CategoriesAPIHandler{service}
}

func (h *CategoriesAPIHandler) ListCategories(writer http.ResponseWriter, request *http.Request) {
	response, err := h.service.GetAllCategories()
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	json.NewEncoder(writer).Encode(response)
}

func (h *CategoriesAPIHandler) GetCategoryBySlug(writer http.ResponseWriter, request *http.Request) {
	slug, ok := handlers.GetParam[string](writer, request, "slug")
	if !ok {
		return
	}

	category, err := h.service.GetCategoryBySlug(slug)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	json.NewEncoder(writer).Encode(category)
}
