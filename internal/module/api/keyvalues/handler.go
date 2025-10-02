package keyvalues

import (
	"bloggo/internal/utils/apierrors"
	"encoding/json"
	"net/http"
)

type KeyValuesAPIHandler struct {
	service KeyValuesAPIService
}

func NewKeyValuesAPIHandler(service KeyValuesAPIService) KeyValuesAPIHandler {
	return KeyValuesAPIHandler{service}
}

func (h *KeyValuesAPIHandler) ListKeyValues(writer http.ResponseWriter, request *http.Request) {
	// Get query parameters
	key := request.URL.Query().Get("key")
	starting := request.URL.Query().Get("starting")

	response, err := h.service.GetKeyValues(key, starting)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	// Return simple array
	json.NewEncoder(writer).Encode(response)
}
