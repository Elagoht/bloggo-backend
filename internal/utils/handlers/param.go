package handlers

import (
	"bloggo/internal/utils/apierrors"
	"net/http"
	"reflect"
	"strconv"

	"github.com/go-chi/chi"
)

// Extracts a chi URL parameter and converts it to the specified type.
func GetParam[T any](
	writer http.ResponseWriter,
	request *http.Request,
	key string,
) (T, bool) {
	var zeroValue T

	value := chi.URLParam(request, key)

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
