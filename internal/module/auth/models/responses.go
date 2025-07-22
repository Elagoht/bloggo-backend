package models

// -- Login or Refresh Response -- //
type ResponseSession struct {
	AccessToken string   `json:"accessToken"`
	Name        string   `json:"name"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}
