package models

// API Tag
type APITag struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

// API Tag Details
type APITagDetails struct {
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	PostCount int    `json:"postCount"`
}

// Response for tags list
type APITagsResponse struct {
	Tags []APITagDetails `json:"tags"`
}
