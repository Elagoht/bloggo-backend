package requests

// -- Create new category -- //
type CategoryCreate struct {
	Name        string `json:"name" validate:"required,max=100"`
	Description string `json:"description" validate:"required,min=70,max=155"`
}

// -- Patch existing category with only given properties -- //
type CategoryUpdate struct {
	Name        string `json:"name" validate:"max=100"`
	Description string `json:"description" validate:"min=70,max=155"`
}
