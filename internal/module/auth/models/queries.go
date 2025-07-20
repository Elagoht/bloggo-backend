package models

// -- Login Related Data -- //
type UserLoginDetails struct {
	UserId         int64
	RoleId         int64
	PassphraseHash string
}
