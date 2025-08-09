package models

// -- New Post Created -- //
type ResponsePostCreated struct {
	Id int64 `json:"id"`
}

// -- Post Details -- //
type ResponsePostDetails struct {
	PostId    int64 `json:"postId"`
	VersionId int64 `json:"versionId"`
	Author    struct {
		Name   string  `json:"name"`
		Email  string  `json:"email"`
		Avatar *string `json:"avatar,omitempty"`
	} `json:"author"`
	Title            *string `json:"title"`
	Slug             *string `json:"slug"`
	Content          *string `json:"content"`
	CoverImage       *string `json:"coverImage"`
	Description      *string `json:"description"`
	Spot             *string `json:"spot"`
	Status           int64   `json:"status"`
	StatusChangedAt  *string `json:"statusChangedAt"`
	StatusChangedBy  *string `json:"statusChangedBy"`
	StatusChangeNote *string `json:"statusChangeNote"`
	CreatedBy        string  `json:"createdBy"`
	CreatedAt        *string `json:"createdAt"`
	UpdatedAt        *string `json:"updatedAt"`
	Category         struct {
		Slug *string `json:"slug"`
		Id   *string `json:"id"`
		Name *string `json:"name"`
	} `json:"category"`
}

// -- Post Card -- //
type ResponsePostCard struct {
	PostId int64 `json:"postId"`
	Author struct {
		Name   string  `json:"name"`
		Avatar *string `json:"avatar"`
	} `json:"author"`
	Title      *string `json:"title"`
	Slug       *string `json:"slug"`
	CoverImage *string `json:"coverImage"`
	Spot       *string `json:"spot"`
	Status     uint8   `json:"status"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
	Category   struct {
		Slug *string `json:"slug"`
		Id   *string `json:"id"`
		Name *string `json:"name"`
	} `json:"category"`
}

// -- Post Versions List Item -- //
type PostVersionsCard struct {
	VersionId     string `json:"id"`
	VersionAuthor struct {
		Id     int64   `json:"id"`
		Name   string  `json:"name"`
		Avatar *string `json:"avatar"`
	} `json:"versionAuthor"`
	Title     *string `json:"title"`
	Status    uint8   `json:"status"`
	UpdatedAt string  `json:"updatedAt"`
}

// -- Post Versions List Item -- //
type ResponseVersionsOfPost struct {
	CurrentVersionId int64  `json:"currentVersionId"`
	CreatedAt        string `json:"createdAt"`
	OriginalAuthor   struct {
		Id     int64   `json:"id"`
		Name   string  `json:"name"`
		Avatar *string `json:"avatar"`
	} `json:"originalAuthor"`
	Versions []PostVersionsCard `json:"versions"`
}
