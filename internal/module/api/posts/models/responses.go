package models

// API Post Card - For list endpoints
type APIPostCard struct {
	Slug        string     `json:"slug"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	Spot        *string    `json:"spot,omitempty"`
	CoverImage  *string    `json:"coverImage,omitempty"`
	ReadCount   int64      `json:"readCount"`
	ReadTime    int        `json:"readTime"`
	PublishedAt string     `json:"publishedAt"`
	Author      APIAuthor  `json:"author"`
	Category    APICategory `json:"category"`
	Tags        []APITag   `json:"tags"`
}

// API Post Details - For single post endpoint
type APIPostDetails struct {
	Slug        string      `json:"slug"`
	Title       string      `json:"title"`
	Content     string      `json:"content"`
	Description *string     `json:"description,omitempty"`
	Spot        *string     `json:"spot,omitempty"`
	CoverImage  *string     `json:"coverImage,omitempty"`
	ReadCount   int64       `json:"readCount"`
	ReadTime    int         `json:"readTime"`
	PublishedAt string      `json:"publishedAt"`
	UpdatedAt   string      `json:"updatedAt"`
	Author      APIAuthor   `json:"author"`
	Category    APICategory `json:"category"`
	Tags        []APITag    `json:"tags"`
}

// API Author
type APIAuthor struct {
	ID     int64   `json:"id"`
	Name   string  `json:"name"`
	Avatar *string `json:"avatar,omitempty"`
}

// API Category
type APICategory struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

// API Tag
type APITag struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

// Paginated response for posts
type APIPostsResponse struct {
	Data  []APIPostCard `json:"data"`
	Page  int           `json:"page"`
	Take  int           `json:"take"`
	Total int64         `json:"total"`
}
