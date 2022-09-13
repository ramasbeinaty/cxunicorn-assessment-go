package auth

import "time"

type UserRegister struct {
	UserDetails User   `json:"user_details"`
	RoleDetails string `json:"role_details"`
}

type User struct {
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	DOB         time.Time `json:"dob"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Role        string    `json:"role"`
}

type UserLogin struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}
