package models

// -- New Category Created -- //
type ResponseCategoryCreated struct {
	Id int64 `json:"id"`
}

// -- Category Details -- //
type ResponseCategoryDetails struct {
	Id          int64   `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description string  `json:"description"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   *string `json:"updatedAt,omitempty"`
}
