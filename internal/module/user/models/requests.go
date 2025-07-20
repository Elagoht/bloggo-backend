package models

// -- Create New User -- //
type RequestUserCreate struct {
	Name       string `json:"name" validate:"required,min=3,max=100"`
	Email      string `json:"email" validate:"required,email,min=5,max=255"`
	Avatar     string `json:"avatar,omitempty" validate:"omitempty,max=100"`
	Passphrase string `json:"passphrase" validate:"omitempty,min=12,max=100"`
	RoleId     int64  `json:"roleId" validate:"required"`
}
