package models

// -- Create new Tag -- //
type RequestTagCreate struct {
	Name string `json:"name" validate:"required,max=100"`
}

// -- Patch existing Tag with only given properties -- //
type RequestTagUpdate struct {
	Name string `json:"name,omitempty" validate:"omitempty,max=100"`
}

// -- Assign tags to a post -- //
type RequestAssignTagsToPost struct {
	TagIds []int64 `json:"tagIds" validate:"required,min=1"`
}

// -- Remove tag from a post -- //
type RequestRemoveTagFromPost struct {
	TagId int64 `json:"tagId" validate:"required"`
}
