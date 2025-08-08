package models

// -- New Post Created -- //
type ResponsePostCreated struct {
	Id int64 `json:"id"`
}

// -- Post Details -- //
type ResponsePostDetails struct {
	PostId    int64
	VersionId int64
	Author    struct {
		Name   string  `json:"author"`
		Email  string  `json:"email"`
		Avatar *string `json:"avatar,omitempty"`
	} `json:"author"`
	Title            string  `json:"title"`
	Slug             string  `json:"slug"`
	Content          string  `json:"content"`
	CoverImage       string  `json:"coverImage"`
	Description      string  `json:"description"`
	Spot             string  `json:"spot"`
	Status           int64   `json:"status"`
	StatusChangedAt  *string `json:"statusChangedAt"`
	StatusChangedBy  *string `json:"statusChangedBy"`
	StatusChangeNote *string `json:"statusChangeNote"`
	IsActive         bool    `json:"isActive"`
	CreatedBy        string  `json:"createdBy"`
	CreatedAt        string  `json:"createdAt"`
	UpdatedAt        string  `json:"updatedAt"`
	Category         struct {
		Slug string `json:"slug"`
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
}

// -- Post Card -- //
type ResponsePostCard struct {
	PostId int64 `json:"postId"`
	Author struct {
		Name   string  `json:"author"`
		Avatar *string `json:"avatar,omitempty"`
	} `json:"author"`
	Title      string `json:"title"`
	Slug       string `json:"slug"`
	CoverImage string `json:"coverImage"`
	Spot       string `json:"spot"`
	Status     string `json:"status"`
	IsActive   string `json:"isActive"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
	Category   struct {
		Slug string `json:"slug"`
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
}
