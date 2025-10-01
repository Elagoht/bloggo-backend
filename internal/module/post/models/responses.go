package models

// -- Post Details -- //
type ResponsePostDetails struct {
	PostId    int64 `json:"postId"`
	VersionId int64 `json:"versionId"`
	Author    struct {
		Id     int64   `json:"id"`
		Name   string  `json:"name"`
		Avatar *string `json:"avatar,omitempty"`
	} `json:"author"`
	Title       *string `json:"title"`
	Slug        *string `json:"slug"`
	Content     *string `json:"content"`
	CoverImage  *string `json:"coverImage"`
	Description *string `json:"description"`
	Spot        *string `json:"spot"`
	Status      int64   `json:"status"`
	ReadCount   int64   `json:"readCount"`
	CreatedAt   *string `json:"createdAt"`
	UpdatedAt   *string `json:"updatedAt"`
	Category    struct {
		Id        *string `json:"id"`
		Name      *string `json:"name"`
		Slug      *string `json:"slug"`
		DeletedAt *string `json:"deletedAt"`
	} `json:"category"`
	Tags []TagCard `json:"tags"`
}

type TagCard struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// -- Post Card -- //
type ResponsePostCard struct {
	PostId int64 `json:"postId"`
	Author struct {
		Id     int64   `json:"id"`
		Name   string  `json:"name"`
		Avatar *string `json:"avatar"`
	} `json:"author"`
	Title      *string `json:"title"`
	Slug       *string `json:"slug"`
	CoverImage *string `json:"coverImage"`
	Spot       *string `json:"spot"`
	Status     int64   `json:"status"`
	ReadCount  int64   `json:"readCount"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
	Category   struct {
		Id        *string `json:"id"`
		Name      *string `json:"name"`
		Slug      *string `json:"slug"`
		DeletedAt *string `json:"deletedAt"`
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
	Title      *string `json:"title"`
	CoverImage *string `json:"coverImage"`
	Status     int64   `json:"status"`
	UpdatedAt  string  `json:"updatedAt"`
	Category   struct {
		Id        *string `json:"id"`
		Name      *string `json:"name"`
		Slug      *string `json:"slug"`
		DeletedAt *string `json:"deletedAt"`
	} `json:"category"`
}

// -- Post Versions List Item -- //
type ResponseVersionsOfPost struct {
	CurrentVersionId *int64 `json:"currentVersionId,omitempty"`
	CreatedAt        string `json:"createdAt"`
	OriginalAuthor   struct {
		Id     int64   `json:"id"`
		Name   string  `json:"name"`
		Avatar *string `json:"avatar"`
	} `json:"originalAuthor"`
	Versions []PostVersionsCard `json:"versions"`
}

// -- Post Specific Version List Item -- //
type ResponseVersionDetailsOfPost struct {
	VersionId      int64  `json:"versionId"`
	DuplicatedFrom *int64 `json:"duplicatedFrom"`
	VersionAuthor  struct {
		Id     int64   `json:"id"`
		Name   string  `json:"name"`
		Avatar *string `json:"avatar,omitempty"`
	} `json:"versionAuthor"`
	Title            *string `json:"title"`
	Slug             *string `json:"slug"`
	Content          *string `json:"content"`
	CoverImage       *string `json:"coverImage"`
	Description      *string `json:"description"`
	Spot             *string `json:"spot"`
	Status           int64   `json:"status"`
	StatusChangedAt  *string `json:"statusChangedAt"`
	StatusChangedBy  *struct {
		Id     int64   `json:"id"`
		Name   string  `json:"name"`
		Avatar *string `json:"avatar"`
	} `json:"statusChangedBy"`
	StatusChangeNote *string `json:"statusChangeNote"`
	CreatedAt        *string `json:"createdAt"`
	UpdatedAt        *string `json:"updatedAt"`
	Category         struct {
		Id        *string `json:"id"`
		Name      *string `json:"name"`
		Slug      *string `json:"slug"`
		DeletedAt *string `json:"deletedAt"`
	} `json:"category"`
}

// -- Version Deletion Response -- //
type ResponseVersionDeleted struct {
	PostDeleted bool `json:"postDeleted,omitempty"`
}
