package cache

import "time"

// defines the storage format of a user
type User struct {
	ID          int       `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	DOB         time.Time `json:"dob"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	IsActive    bool      `json:"is_active"`
	IsVerified  bool      `json:"is_verified"`
}
