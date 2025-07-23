package models

// -- Login To System -- //
type RequestSessionCreate struct {
	Email      string `json:"email" validate:"required,email"`
	Passphrase string `json:"passphrase" validate:"required"`
}
