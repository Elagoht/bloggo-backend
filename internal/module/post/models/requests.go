package models

import "mime/multipart"

// -- Create New Version -- //
type RequestPostUpsert struct {
	// Everything is optional to create a draft
	// Will be validated while publishing
	Title       string                `json:"title" validate:"max=100"`
	Content     *string               `json:"content"`
	Cover       *multipart.FileHeader `form:"cover"`
	Description *string               `json:"description" validate:"omitempty,max=155"`
	Spot        *string               `json:"spot" validate:"omitempty,max=75"`
	CategoryId  *int64                `json:"categoryId"`
}

// -- Change Status With Commit Note -- //
type RequestPostStatusModerate struct {
	Id   int64  `json:"versionId" validate:"required"`
	Note string `json:"note" validate:"max=255,required"`
}

// -- Track Post View -- //
type RequestTrackView struct {
	PostId    int64  `json:"postId" validate:"required"`
	UserAgent string `json:"userAgent" validate:"required"`
}

// -- List Posts With Filters -- //
type RequestPostFilters struct {
	Page       *int    `form:"page" validate:"omitempty,min=1"`
	Take       *int    `form:"take" validate:"omitempty,min=1,max=100"`
	Order      *string `form:"order" validate:"omitempty,oneof=title created_at updated_at read_count"`
	Dir        *string `form:"dir" validate:"omitempty,oneof=asc desc"`
	Q          *string `form:"q" validate:"omitempty,max=100"`
	Status     *int    `form:"status" validate:"omitempty,min=0,max=5"`
	CategoryId *int64  `form:"categoryId"`
	AuthorId   *int64  `form:"authorId"`
}

// -- Validation for Post Submission -- //
type PostSubmissionValidation struct {
	Title       string `validate:"required,max=100"`
	Content     string `validate:"required,min=100"`
	Description string `validate:"required,max=155"`
	Spot        string `validate:"required,max=75"`
	CategoryId  int64  `validate:"required"`
}

// -- Assign Tags to Post -- //
type RequestAssignTagsToPost struct {
	TagIds []int64 `json:"tagIds" validate:"required"`
}
