package models

import "time"

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  *string   `json:"username,omitempty"`
	FirstName *string   `json:"firstName,omitempty"`
	LastName  *string   `json:"lastName,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// UserSession represents an authenticated user session
type UserSession struct {
	User      User   `json:"user"`
	LoggedIn  bool   `json:"loggedIn"`
	SessionID string `json:"sessionId,omitempty"`
}

// AuthResponse represents the response from authentication endpoints
type AuthResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message,omitempty"`
	User    *UserSession `json:"user,omitempty"`
}
