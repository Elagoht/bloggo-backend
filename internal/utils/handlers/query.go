package handlers

import (
	"bloggo/internal/utils/apierrors"
	"net/http"
	"reflect"
	"strconv"
)

// Extracts a query parameter from the request
// and converts it to the specified type.
func GetQuery[T any](
	writer http.ResponseWriter,
	request *http.Request,
	key string,
	required ...bool,
) (T, bool) {
	var zeroValue T

	values := request.URL.Query()[key]

	// Make the required field optional
	isRequired := len(required) > 0 && required[0]

	if len(values) == 0 || values[0] == "" {
		if isRequired {
			WriteError(
				writer,
				apierrors.NewAPIError("Field \""+key+"\" is required", nil),
				http.StatusBadRequest,
			)
		}
		return zeroValue, false
	}
	value := values[0]

	var result any
	success := true

	typeOfType := reflect.TypeOf(zeroValue)
	switch typeOfType.Kind() {
	case reflect.String:
		result = value
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			success = false
		}
		result = reflect.ValueOf(parsed).Convert(typeOfType).Interface()
	case reflect.Float32, reflect.Float64:
		parsed, err := strconv.ParseFloat(value, 64)
		if err != nil {
			success = false
		}
		result = reflect.ValueOf(parsed).Convert(typeOfType).Interface()
	case reflect.Bool:
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			success = false
		}
		result = parsed
	default:
		success = false
	}

	if !success {
		WriteError(
			writer,
			apierrors.NewAPIError("Field \""+key+"\" could not validated", nil),
			http.StatusUnprocessableEntity,
		)
		return zeroValue, false
	}

	return result.(T), true
}
