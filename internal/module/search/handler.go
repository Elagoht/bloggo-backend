package search

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type SearchHandler struct {
	service SearchService
}

func NewSearchHandler(service SearchService) SearchHandler {
	return SearchHandler{
		service: service,
	}
}

// Search handles GET /search
func (handler *SearchHandler) Search(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query().Get("q")
	if query == "" {
		http.Error(writer, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	limitStr := request.URL.Query().Get("limit")
	limit := 10 // default limit

	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	response, err := handler.service.Search(query, limit)
	if err != nil {
		http.Error(writer, "Failed to perform search", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)
}