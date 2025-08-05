package models

type ClerkUser struct {
	ID        string `json:"sub"`
	Email     string `json:"email"`
	ImageURL  string `json:"image_url"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
}
