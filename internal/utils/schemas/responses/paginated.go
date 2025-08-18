package responses

// PaginatedResponse represents a generic paginated response
type PaginatedResponse[T any] struct {
	Data  []T   `json:"data"`
	Page  int   `json:"page"`
	Take  int   `json:"take"`
	Total int64 `json:"total"`
}
