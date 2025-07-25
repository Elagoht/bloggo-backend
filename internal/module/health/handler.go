package health

import (
	"bloggo/internal/module/health/models"
	"encoding/json"
	"net/http"
	"time"
)

type HealthHandler struct{}

func NewHealthHandler() HealthHandler {
	return HealthHandler{}
}

func (handler *HealthHandler) CheckHealth(
	writer http.ResponseWriter,
	request *http.Request,
) {
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(&models.ResponseHealth{
		Status: true,
		Time:   time.Now().String(),
	})
}
