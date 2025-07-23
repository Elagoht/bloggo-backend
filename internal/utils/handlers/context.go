package handlers

import (
	"bloggo/internal/utils/apierrors"
	"net/http"
	"reflect"
	"strconv"
)

// GetContextValue retrieves a value from the request context and converts it to the specified type.
func GetContextValue[T any](
	writer http.ResponseWriter,
	request *http.Request,
	key string,
) (T, bool) {
	var zeroValue T

	// Retrieve value from context
	ctxValue := request.Context().Value(key)
	if ctxValue == nil {
		WriteError(
			writer,
			apierrors.NewAPIError("Field \""+key+"\" not found in context", nil),
			http.StatusInternalServerError,
		)
		return zeroValue, false
	}

	// Get the type of T
	typeOfType := reflect.TypeOf(zeroValue)
	var result any
	success := true

	// Handle type conversion based on the type of T
	switch typeOfType.Kind() {
	case reflect.String:
		if str, ok := ctxValue.(string); ok {
			result = str
		} else {
			success = false
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch v := ctxValue.(type) {
		case int:
			result = reflect.ValueOf(v).Convert(typeOfType).Interface()
		case int64:
			result = reflect.ValueOf(v).Convert(typeOfType).Interface()
		case string:
			parsed, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				success = false
			}
			result = reflect.ValueOf(parsed).Convert(typeOfType).Interface()
		default:
			success = false
		}
	case reflect.Float32, reflect.Float64:
		switch v := ctxValue.(type) {
		case float64:
			result = reflect.ValueOf(v).Convert(typeOfType).Interface()
		case string:
			parsed, err := strconv.ParseFloat(v, 64)
			if err != nil {
				success = false
			}
			result = reflect.ValueOf(parsed).Convert(typeOfType).Interface()
		default:
			success = false
		}
	case reflect.Bool:
		switch v := ctxValue.(type) {
		case bool:
			result = v
		case string:
			parsed, err := strconv.ParseBool(v)
			if err != nil {
				success = false
			}
			result = parsed
		default:
			success = false
		}
	default:
		success = false
	}

	if !success {
		WriteError(
			writer,
			apierrors.NewAPIError("Field \""+key+"\" could not be converted to requested type", nil),
			http.StatusUnprocessableEntity,
		)
		return zeroValue, false
	}

	return result.(T), true
}
