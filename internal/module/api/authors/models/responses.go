package models

// API Author
type APIAuthor struct {
	ID     int64   `json:"id"`
	Name   string  `json:"name"`
	Avatar *string `json:"avatar,omitempty"`
}

// API Author Details
type APIAuthorDetails struct {
	ID                  int64   `json:"id"`
	Name                string  `json:"name"`
	Avatar              *string `json:"avatar,omitempty"`
	PublishedPostCount  int     `json:"publishedPostCount"`
	MemberSince         string  `json:"memberSince"`
}

// Response for authors list
type APIAuthorsResponse struct {
	Authors []APIAuthorDetails `json:"authors"`
}
