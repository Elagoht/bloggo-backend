package models

// -- Create new category -- //
type RequestCategoryCreate struct {
	Name        string `json:"name" validate:"required,max=100"`
	Spot        string `json:"spot" validate:"required,min=20,max=75"`
	Description string `json:"description" validate:"required,min=70,max=155"`
}

// -- Patch existing category with only given properties -- //
type RequestCategoryUpdate struct {
	Name        string `json:"name,omitempty" validate:"omitempty,max=100"`
	Spot        string `json:"spot,omitempty" validate:"omitempty,min=20,max=75"`
	Description string `json:"description,omitempty" validate:"omitempty,min=70,max=155"`
}
