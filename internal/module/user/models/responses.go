package models

// -- User Card -- //
type ResponseUserCard struct {
	Id                 int64  `json:"id"`
	Name               string `json:"name"`
	Email              string `json:"email"`
	Avatar             string `json:"avatar,omitempty"`
	RoleId             int64  `json:"roleId"`
	RoleName           string `json:"roleName"`
	WrittenPostCount   int64  `json:"writtenPostCount"`
	PublishedPostCount int64  `json:"publishedPostCount"`
}

// -- User Details -- //
type ResponseUserDetails struct {
	Id                 int64   `json:"id"`
	Name               string  `json:"name"`
	Email              string  `json:"email"`
	Avatar             *string `json:"avatar,omitempty"`
	CreatedAt          string  `json:"createdAt"`
	LastLogin          *string `json:"lastLogin,omitempty"`
	RoleId             int64   `json:"roleId"`
	RoleName           string  `json:"roleName"`
	WrittenPostCount   int64   `json:"writtenPostCount"`
	PublishedPostCount int64   `json:"publishedPostCount"`
}
