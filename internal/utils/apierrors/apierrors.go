package apierrors

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"strings"

	"github.com/mattn/go-sqlite3"
)

type APIError struct {
	Message string            `json:"message,omitempty"`
	Stack   error             `json:"-"`
	Details []ValidationError `json:"details,omitempty"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type HTTPErrorMapItem struct {
	Message string
	Status  int
}

type HTTPErrorMapping map[error]HTTPErrorMapItem

func NewAPIError(message string, err error) *APIError {
	return &APIError{
		Message: message,
		Stack:   err,
	}
}

func NewValidationAPIError(errors []ValidationError) *APIError {
	return &APIError{
		Message: "Validation failed",
		Details: errors,
	}
}

func extractConflictFields(err error) []ValidationError {
	var details []ValidationError
	if err != nil {
		msg := err.Error()
		if idx := strings.Index(msg, ": "); idx != -1 {
			fields := strings.Split(msg[idx+2:], ", ")
			for _, field := range fields {
				// Only take the last part after the dot, e.g., categories.slug -> slug
				parts := strings.Split(field, ".")
				fieldName := parts[len(parts)-1]
				details = append(details, ValidationError{
					Field:   fieldName,
					Message: "conflict on this field",
				})
			}
		}
	}
	return details
}

func MapErrors(
	err error,
	writer http.ResponseWriter,
	customErrors HTTPErrorMapping,
) {
	log.Default().Println(err)

	var item HTTPErrorMapItem
	found := false

	// Handle sqlite3 errors
	if sqliteErr, ok := err.(sqlite3.Error); ok {
		// Temporarily mark as found, default case will set false
		found = true

		switch sqliteErr.Code {
		case sqlite3.ErrConstraint:
			item = (*DefaultErrorMapping)[ErrConflict]

			details := extractConflictFields(err)
			apiErr := NewAPIError(item.Message, err)
			apiErr.Details = details
			writer.WriteHeader(item.Status)
			json.NewEncoder(writer).Encode(apiErr)
			return
		case sqlite3.ErrNotFound:
			item = (*DefaultErrorMapping)[ErrNotFound]
		case sqlite3.ErrAuth:
			item = (*DefaultErrorMapping)[ErrUnauthorized]
		case sqlite3.ErrBusy:
			item = (*DefaultErrorMapping)[ErrServiceUnavailable]
		case sqlite3.ErrEmpty:
		case sqlite3.ErrNotFound:
			item = (*DefaultErrorMapping)[ErrNotFound]
		case sqlite3.ErrError:
			item = (*DefaultErrorMapping)[ErrBadRequest]
		default:
			// Re-set as false if not matched
			found = false
		}
	}

	if !found {
		for key, val := range customErrors {
			if errors.Is(err, key) {
				item = val
				found = true
				break
			}
		}
	}

	if !found {
		for key, val := range *DefaultErrorMapping {
			if errors.Is(err, key) {
				item = val
				found = true
				break
			}
		}
	}

	if found {
		apiErr := NewAPIError(item.Message, err)
		writer.WriteHeader(item.Status)
		json.NewEncoder(writer).Encode(apiErr)
		return
	}

	// Fallback: unknown error
	apiErr := NewAPIError(ErrInternalServer.Error(), err)
	writer.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(writer).Encode(apiErr)
}

var (
	ErrConflict             = errors.New("resource conflict")
	ErrNotFound             = errors.New("resource not found")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrForbidden            = errors.New("forbidden")
	ErrBadRequest           = errors.New("bad request")
	ErrInternalServer       = errors.New("internal server error")
	ErrRequestTimeout       = errors.New("request timeout")
	ErrTooManyRequests      = errors.New("too many requests")
	ErrUnprocessableEntity  = errors.New("unprocessable entity")
	ErrNotImplemented       = errors.New("not implemented")
	ErrServiceUnavailable   = errors.New("service unavailable")
	ErrUnsupportedMediaType = errors.New("unsupported media type")
	ErrPreconditionFailed   = errors.New("precondition failed")
	ErrEncryptionError      = errors.New("encyption failed")
)

var DefaultErrorMapping = &HTTPErrorMapping{
	// Standard Errors
	ErrConflict:             {ErrConflict.Error(), http.StatusConflict},
	ErrNotFound:             {ErrNotFound.Error(), http.StatusNotFound},
	ErrUnauthorized:         {ErrUnauthorized.Error(), http.StatusUnauthorized},
	ErrForbidden:            {ErrForbidden.Error(), http.StatusForbidden},
	ErrBadRequest:           {ErrBadRequest.Error(), http.StatusBadRequest},
	ErrInternalServer:       {ErrInternalServer.Error(), http.StatusInternalServerError},
	ErrRequestTimeout:       {ErrRequestTimeout.Error(), http.StatusRequestTimeout},
	ErrTooManyRequests:      {ErrTooManyRequests.Error(), http.StatusTooManyRequests},
	ErrUnprocessableEntity:  {ErrUnprocessableEntity.Error(), http.StatusUnprocessableEntity},
	ErrNotImplemented:       {ErrNotImplemented.Error(), http.StatusNotImplemented},
	ErrServiceUnavailable:   {ErrServiceUnavailable.Error(), http.StatusServiceUnavailable},
	ErrUnsupportedMediaType: {ErrUnsupportedMediaType.Error(), http.StatusUnsupportedMediaType},
	ErrPreconditionFailed:   {ErrPreconditionFailed.Error(), http.StatusPreconditionFailed},
	// Standard SQL errors
	sql.ErrNoRows: {ErrNotFound.Error(), http.StatusNotFound},
}

// To implement default Error interface
func (err *APIError) Error() string {
	return err.Message
}
