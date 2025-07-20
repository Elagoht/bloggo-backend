package filter

import (
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/handlers"
	"fmt"
	"net/http"
	"strings"
)

// SearchOptions holds parameters for search queries
// Only Q: the search string
// If Q is nil, no search is performed
// The fields to search are now provided by the programmer, not the user
type SearchOptions struct {
	Q *string
}

// GetSearchOptions extracts only the search string (q) from the request
// Returns nil if no search string is provided
func GetSearchOptions(
	writer http.ResponseWriter,
	request *http.Request,
) (*SearchOptions, bool) {
	q, okQ := handlers.GetQuery[string](writer, request, "q")

	var qPtr *string

	if okQ {
		trimmed := strings.TrimSpace(q)
		if trimmed == "" {
			handlers.WriteError(
				writer,
				apierrors.NewAPIError("'q' (search query) must not be empty", nil),
				http.StatusBadRequest,
			)
			return nil, false
		}
		qPtr = &trimmed
	}

	return &SearchOptions{
		Q: qPtr,
	}, true
}

// BuildSearchClause builds a SQL WHERE clause for LIKE-based search
// The fields to search are provided by the programmer
// Example: WHERE (field1 LIKE ? OR field2 LIKE ?)
// Returns the clause and the args (with %q%)
func BuildSearchClause(
	opts *SearchOptions,
	fields []string,
) (clause string, args []any) {
	if opts == nil || opts.Q == nil || len(fields) == 0 {
		return "", nil
	}
	var parts []string
	for range fields {
		parts = append(parts, "%s LIKE ?")
	}
	clause = fmt.Sprintf("(%s)", strings.Join(parts, " OR "))
	args = make([]any, len(fields))
	for i := range fields {
		args[i] = "%" + *opts.Q + "%"
	}
	// Replace %s with field names
	for _, f := range fields {
		clause = strings.Replace(clause, "%s", f, 1)
	}
	if clause != "" {
		clause = "AND " + clause
	}
	return
}
