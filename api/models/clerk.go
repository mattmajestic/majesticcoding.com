package models

type ClerkSession struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	SessionID string `json:"session_id"`
}

type ClerkUser struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}
