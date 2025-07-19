package category

import (
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

func (h CategoryHandler) List(
	writer http.ResponseWriter,
	request *http.Request,
) {
	categories, err := h.service.ListCategories()
	if err != nil {
		return
	}

	json.NewEncoder(writer).Encode(categories)
}

func (h CategoryHandler) Create(
	writer http.ResponseWriter,
	request *http.Request,
) {
	var input CreateCategoryRequest
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		return
	}

	category, err := h.service.CreateCategory(input)
	if err != nil {
		return
	}

	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(category)
}

func (h CategoryHandler) Update(
	writer http.ResponseWriter,
	request *http.Request,
) {
	var input UpdateCategoryRequest
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		return
	}

	err := h.service.UpdateCategory(input)
	if err != nil {
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func (h CategoryHandler) Delete(
	writer http.ResponseWriter,
	request *http.Request,
) {
	id := request.URL.Query().Get("id")
	if id == "" {
		return
	}

	err := h.service.DeleteCategory(id)
	if err != nil {
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
