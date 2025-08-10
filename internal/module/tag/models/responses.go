package models

// -- Tag Details -- //
type ResponseTagDetails struct {
	Id int64 `json:"id"`

	Name string `json:"name"`
	Slug string `json:"slug"`

	BlogCount int `json:"blogCount"`

	CreatedAt string  `json:"createdAt"`
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

// -- Tag Card -- //
type ResponseTagCard struct {
	Id int64 `json:"id"`

	Name string `json:"name"`
	Slug string `json:"slug"`

	BlogCount int `json:"blogCount"`
}
