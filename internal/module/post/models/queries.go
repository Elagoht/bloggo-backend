package models

type QueryGetPostVersionDuplicateData struct {
	PostId      int64
	Title       *string
	Slug        *string
	Content     *string
	CoverImage  *string
	Description *string
	Spot        *string
	CategoryId  *int64
	CreatedBy   int64
}
