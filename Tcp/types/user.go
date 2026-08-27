package types

// User represents a user in the system.
// This is the core domain model for a user.
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
