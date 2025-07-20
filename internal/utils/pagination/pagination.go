package pagination

import (
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/handlers"
	"fmt"
	"net/http"
	"slices"
	"strings"
)

// PaginationOptions holds optional parameters for pagination and ordering
// This struct is reusable across modules
type PaginationOptions struct {
	Page      *int
	Take      *int
	OrderBy   *string
	Direction *string // "asc" or "desc"
}

func GetPaginationOptions(
	writer http.ResponseWriter,
	request *http.Request,
	orderWhiteList []string,
) (*PaginationOptions, bool) {
	page, okPage := handlers.GetQuery[int](writer, request, "page")
	take, okTake := handlers.GetQuery[int](writer, request, "take")
	order, okOrder := handlers.GetQuery[string](writer, request, "order")
	direction, okDirection := handlers.GetQuery[string](writer, request, "dir")

	var pagePtr, takePtr *int
	var orderByPtr, orderDirectionPtr *string

	if okPage {
		if page < 1 {
			handlers.WriteError(
				writer,
				apierrors.NewAPIError("'page' must be a positive integer", nil),
				http.StatusBadRequest,
			)
			return nil, false
		} else {
			pagePtr = &page
		}
	}

	if okTake {
		if take < 1 {
			handlers.WriteError(
				writer,
				apierrors.NewAPIError("'take' must be a positive integer", nil),
				http.StatusBadRequest,
			)
			return nil, false
		} else {
			takePtr = &take
		}
	}

	if okOrder {
		found := slices.Contains(orderWhiteList, order)
		allowedFields := strings.Join(orderWhiteList, ", ")
		if !found {
			handlers.WriteError(
				writer,
				apierrors.NewAPIError(
					fmt.Sprintf("'order' must be one of the '%s'", allowedFields),
					nil,
				),
				http.StatusBadRequest,
			)
			return nil, false
		}
		orderByPtr = &order
	}

	if okDirection {
		direction := strings.ToLower(direction)
		if direction != "asc" && direction != "desc" {
			handlers.WriteError(
				writer,
				apierrors.NewAPIError("dir' must be 'asc' or 'desc'", nil),
				http.StatusBadRequest,
			)
			return nil, false
		} else {
			orderDirectionPtr = &direction
		}
	}

	return &PaginationOptions{
		Page:      pagePtr,
		Take:      takePtr,
		OrderBy:   orderByPtr,
		Direction: orderDirectionPtr,
	}, true
}

// BuildPaginationClauses returns the SQL clauses and args for LIMIT, OFFSET, and ORDER BY
func (opts PaginationOptions) BuildPaginationClauses() (
	orderByClause string,
	limitClause string,
	offsetClause string,
	parinationArgs []any,
) {
	orderByClause = ""
	if opts.OrderBy != nil {
		orderByClause = fmt.Sprintf("ORDER BY %s", *opts.OrderBy)
		if opts.Direction != nil {
			orderByClause += " " + *opts.Direction
		}
	}

	limitClause = ""
	if opts.Take != nil {
		limitClause = "LIMIT ?"
	}

	offsetClause = ""
	if opts.Page != nil && opts.Take != nil {
		offsetClause = "OFFSET ?"
	}

	parinationArgs = []any{}
	if opts.Take != nil {
		parinationArgs = append(parinationArgs, *opts.Take)
	}
	if opts.Page != nil && opts.Take != nil {
		offset := (*opts.Page - 1) * (*opts.Take)
		parinationArgs = append(parinationArgs, offset)
	}

	return
}
