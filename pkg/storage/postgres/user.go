package postgres

import "time"

// defines the storage format of a user
type User struct {
	ID          int       `json:"id" db:"id"`
	FirstName   string    `json:"first_name" db:"users.first_name"`
	LastName    string    `json:"last_name" db:"users.last_name"`
	DOB         time.Time `json:"dob" db:"dob"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	Email       string    `json:"email" db:"email"`
	Password    string    `json:"password" db:"password"`
	Role        string    `json:"role" db:"role"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	IsVerified  bool      `json:"is_verified" db:"is_verified"`
}

type UserCreate struct {
	FirstName   string    `json:"first_name" db:"users.first_name"`
	LastName    string    `json:"last_name" db:"users.last_name"`
	DOB         time.Time `json:"dob" db:"dob"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	Email       string    `json:"email" db:"email"`
	Password    string    `json:"password" db:"password"`
	Role        string    `json:"role" db:"role"`
}

type UserLogin struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}
