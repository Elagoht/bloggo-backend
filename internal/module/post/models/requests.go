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
