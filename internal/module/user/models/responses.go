package models

// -- New User Created -- //
type ResponseUserCreated struct {
	Id int64 `json:"id"`
}

// -- User Card -- //
type ResponseUserCard struct {
	Id                 int64  `json:"id"`
	Name               string `json:"name"`
	Email              string `json:"email"`
	Avatar             string `json:"avatar,omitempty"`
	RoleId             string `json:"roleId"`
	RoleName           string `json:"roleName"`
	WrittenPostCount   string `json:"writtenPostCount"`
	PublishedPostCount string `json:"publishedPostCount"`
}

// -- User Details -- //
type ResponseUserDetails struct {
	Id                 int64   `json:"id"`
	Name               string  `json:"name"`
	Email              string  `json:"email"`
	Avatar             *string `json:"avatar,omitempty"`
	CreatedAt          string  `json:"createdAt"`
	LastLogin          *string `json:"lastLogin,omitempty"`
	RoleId             string  `json:"roleId"`
	RoleName           string  `json:"roleName"`
	WrittenPostCount   string  `json:"writtenPostCount"`
	PublishedPostCount string  `json:"publishedPostCount"`
}
