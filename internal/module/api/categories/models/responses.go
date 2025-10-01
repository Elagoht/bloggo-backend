package models

// API Category
type APICategory struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

// API Category Details
type APICategoryDetails struct {
	Slug        string  `json:"slug"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Spot        *string `json:"spot,omitempty"`
	PostCount   int     `json:"postCount"`
}

// Response for categories list
type APICategoriesResponse struct {
	Categories []APICategoryDetails `json:"categories"`
}
