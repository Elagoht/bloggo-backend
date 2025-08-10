package handlers

import (
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"bloggo/internal/utils/apierrors"

	"github.com/go-playground/validator/v10"
)

// BindAndValidateMultipart parses multipart/form-data into a struct (T can be struct or *struct),
// fills *multipart.FileHeader fields, converts basic types (int, bool, float, string, []string),
// performs validation (bindValidater) and writes errors to the provided writer.
// Returns the populated body and a bool success flag.
func BindAndValidateMultipart[T any](
	writer http.ResponseWriter,
	request *http.Request,
	maxFileSize int64,
) (T, bool) {
	var body T

	// Limit request size to avoid huge uploads
	request.Body = http.MaxBytesReader(writer, request.Body, maxFileSize)

	// Parse multipart form
	if err := request.ParseMultipartForm(maxFileSize); err != nil {
		WriteError(
			writer,
			apierrors.NewAPIError("Failed to parse multipart form: "+err.Error(), err),
			http.StatusBadRequest,
		)
		return body, false
	}

	// Reflect types
	typeOf := reflect.TypeOf(body)
	if typeOf == nil {
		// defensive - get type via pointer trick
		typeOf = reflect.TypeOf((*T)(nil)).Elem()
	}
	// handle pointer-to-struct generics, get struct type
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}
	if typeOf.Kind() != reflect.Struct {
		WriteError(
			writer,
			apierrors.NewAPIError(
				"BindAndValidateMultipart: destination must be a struct or pointer to struct",
				nil,
			),
			http.StatusInternalServerError,
		)
		return body, false
	}

	// Prepare an addressable value for setting fields.
	valueOf := reflect.ValueOf(&body).Elem()
	// If T is a pointer to struct (e.g. *MyReq), allocate a new struct and set it.
	if valueOf.Kind() == reflect.Ptr {
		valueOf.Set(reflect.New(valueOf.Type().Elem()))
		valueOf = valueOf.Elem()
	}

	// Helper types
	fileHeaderType := reflect.TypeOf(&multipart.FileHeader{})

	// Iterate fields and fill from form
	for index := 0; index < typeOf.NumField(); index++ {
		field := typeOf.Field(index)
		fieldValue := valueOf.Field(index)

		// Skip unexported fields
		if !fieldValue.CanSet() {
			continue
		}

		// determine form key: form tag first, then json tag, then lowercase field name
		formKey := field.Tag.Get("form")
		if formKey == "" {
			formKey = field.Tag.Get("json")
		}
		if formKey == "" {
			formKey = strings.ToLower(field.Name)
		}

		// FILE FIELD
		if field.Type == fileHeaderType ||
			field.Type.AssignableTo(fileHeaderType) {
			multipartFile, fileHeader, err := request.FormFile(formKey)
			if err != nil {
				// Missing file is not necessarily an error — let validator detect if required
				// Only surface real parse errors
				if err == http.ErrMissingFile ||
					strings.Contains(err.Error(), "no such file") {
					continue
				}
				WriteError(
					writer,
					apierrors.NewAPIError(
						"Error retrieving file '"+formKey+"': "+err.Error(),
						err,
					),
					http.StatusBadRequest,
				)
				return body, false
			}
			// Close the temporary multipart.File — FileHeader.Open() can be used later by caller/service.
			_ = multipartFile.Close()

			// Size check (extra guard)
			if fileHeader.Size > maxFileSize {
				WriteError(
					writer,
					apierrors.NewAPIError("File size exceeds maximum limit", nil),
					http.StatusRequestEntityTooLarge,
				)
				return body, false
			}

			fieldValue.Set(reflect.ValueOf(fileHeader))
			continue
		}

		// NON-FILE FIELD
		// If the field can be provided multiple times (slice of strings), use MultipartForm.Value
		if fieldValue.Kind() == reflect.Slice &&
			fieldValue.Type().Elem().Kind() == reflect.String {
			if request.MultipartForm != nil {
				if valSlice, ok := request.MultipartForm.Value[formKey]; ok {
					sliceVal := reflect.MakeSlice(
						fieldValue.Type(),
						len(valSlice),
						len(valSlice),
					)
					for idx, sliceString := range valSlice {
						sliceVal.Index(idx).SetString(sliceString)
					}
					fieldValue.Set(sliceVal)
				}
			}
			continue
		}

		// Single value
		val := request.FormValue(formKey)
		if val == "" {
			continue
		}

		switch fieldValue.Kind() {
		case reflect.String:
			fieldValue.SetString(val)
		case
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			bits := fieldValue.Type().Bits()

			parsed, err := strconv.ParseInt(val, 10, bits)
			if err != nil {
				WriteError(
					writer,
					apierrors.NewAPIError("Field '"+formKey+"' must be an integer", err),
					http.StatusUnprocessableEntity,
				)
				return body, false
			}
			fieldValue.SetInt(parsed)
		case
			reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64:

			bits := fieldValue.Type().Bits()
			parsed, err := strconv.ParseUint(val, 10, bits)
			if err != nil {
				WriteError(
					writer,
					apierrors.NewAPIError(
						"Field '"+formKey+"' must be an unsigned integer",
						err,
					),
					http.StatusUnprocessableEntity,
				)
				return body, false
			}
			fieldValue.SetUint(parsed)
		case reflect.Float32, reflect.Float64:
			bits := fieldValue.Type().Bits()
			parsed, err := strconv.ParseFloat(val, bits)
			if err != nil {
				WriteError(
					writer,
					apierrors.NewAPIError("Field '"+formKey+"' must be a float", err),
					http.StatusUnprocessableEntity,
				)
				return body, false
			}
			fieldValue.SetFloat(parsed)
		case reflect.Bool:
			parsed, err := strconv.ParseBool(val)
			if err != nil {
				WriteError(
					writer,
					apierrors.NewAPIError("Field '"+formKey+"' must be a boolean", err),
					http.StatusUnprocessableEntity,
				)
				return body, false
			}
			fieldValue.SetBool(parsed)
		default:
			// unsupported kind: ignore (validator can catch if needed)
			continue
		}
	}

	// run validation (uses package-level bindValidater from your handlers package)
	if err := bindValidater.Struct(body); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var vErrs []apierrors.ValidationError
			for _, validationError := range validationErrors {
				vErrs = append(vErrs, apierrors.ValidationError{
					Field:   validationError.Field(),
					Message: getValidationErrorMessage(validationError),
				})
			}
			WriteError(
				writer,
				apierrors.NewValidationAPIError(vErrs),
				http.StatusBadRequest,
			)
			return body, false
		}
		WriteError(
			writer,
			apierrors.NewAPIError("Validation failed: "+err.Error(), err),
			http.StatusBadRequest,
		)
		return body, false
	}

	return body, true
}
