package models

// -- Login Related Data -- //
type UserLoginDetails struct {
	UserId         int64
	RoleId         int64
	UserName       string
	RoleName       string
	PassphraseHash string
}
