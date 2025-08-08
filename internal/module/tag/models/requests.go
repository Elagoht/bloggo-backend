package models

// -- Create new Tag -- //
type RequestTagCreate struct {
	Name string `json:"name" validate:"required,max=100"`
}

// -- Patch existing Tag with only given properties -- //
type RequestTagUpdate struct {
	Name string `json:"name,omitempty" validate:"omitempty,max=100"`
}
