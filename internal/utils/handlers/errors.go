package handlers

import (
	"bloggo/internal/utils/apierrors"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func WriteError(
	writer http.ResponseWriter,
	err *apierrors.APIError,
	status int,
) {
	writer.WriteHeader(status)
	json.NewEncoder(writer).Encode(err)
}

func WriteValidationError(
	writer http.ResponseWriter,
	err *apierrors.APIError,
) {
	writer.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(writer).Encode(err)
}

func getValidationErrorMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return fieldError.Field() + " is required"
	case "min":
		return fieldError.Field() + " must be at least " +
			fieldError.Param() + " characters"
	case "max":
		return fieldError.Field() + " must be at most " +
			fieldError.Param() + " characters"
	case "email":
		return fieldError.Field() + " must be a valid email address"
	case "url":
		return fieldError.Field() + " must be a valid URL"
	case "port":
		return fieldError.Field() + " must be a valid port (1025-65535)"
	case "safePath":
		return fieldError.Field() + " must be a valid path without traversal"
	case "sqlnull":
		return fieldError.Field() + " must be a valid SQL null type"
	default:
		return fieldError.Field() + " is invalid: " + fieldError.Tag()
	}
}
