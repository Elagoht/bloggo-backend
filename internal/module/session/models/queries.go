package models

// -- Login Related Data -- //
type SessionCreateDetails struct {
	UserId         int64
	RoleId         int64
	UserName       string
	RoleName       string
	PassphraseHash string
}
