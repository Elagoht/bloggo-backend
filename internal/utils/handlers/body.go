package handlers

import (
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/validate"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var bindValidater = validate.GetValidator()

func BindAndValidate[T any](
	writer http.ResponseWriter,
	request *http.Request,
) (T, bool) {
	var body T
	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		WriteError(
			writer,
			apierrors.NewAPIError("Invalid JSON format", err),
			http.StatusBadRequest,
		)
		return body, false
	}

	if err := bindValidater.Struct(body); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errors []apierrors.ValidationError
			for _, validationError := range validationErrors {
				errors = append(errors, apierrors.ValidationError{
					Field:   validationError.Field(),
					Message: getValidationErrorMessage(validationError),
				})
			}
			WriteError(
				writer,
				apierrors.NewValidationAPIError(errors),
				http.StatusBadRequest,
			)
			return body, false
		}
		WriteError(
			writer,
			apierrors.NewAPIError("Validation failed", err),
			http.StatusBadRequest,
		)
		return body, false
	}

	return body, true
}
