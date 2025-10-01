package authors

import (
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/handlers"
	"encoding/json"
	"net/http"
)

type AuthorsAPIHandler struct {
	service AuthorsAPIService
}

func NewAuthorsAPIHandler(service AuthorsAPIService) AuthorsAPIHandler {
	return AuthorsAPIHandler{service}
}

func (h *AuthorsAPIHandler) ListAuthors(writer http.ResponseWriter, request *http.Request) {
	response, err := h.service.GetAllAuthors()
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	json.NewEncoder(writer).Encode(response)
}

func (h *AuthorsAPIHandler) GetAuthorById(writer http.ResponseWriter, request *http.Request) {
	id, ok := handlers.GetParam[int64](writer, request, "id")
	if !ok {
		return
	}

	author, err := h.service.GetAuthorById(id)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	json.NewEncoder(writer).Encode(author)
}
