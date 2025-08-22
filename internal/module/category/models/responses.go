package models

// -- Category Details -- //
type ResponseCategoryDetails struct {
	Id          int64   `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Spot        string  `json:"spot"`
	Description string  `json:"description"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   *string `json:"updatedAt,omitempty"`
	BlogCount   int     `json:"blogCount"`
}

// -- Category Card -- //
type ResponseCategoryCard struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Spot      string `json:"spot"`
	Slug      string `json:"slug"`
	BlogCount int    `json:"blogCount"`
}

// -- Category List Item -- //
type ResponseCategoryListItem struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
