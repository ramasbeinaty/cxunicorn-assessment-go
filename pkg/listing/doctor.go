package listing

import "time"

// defines the doctor properties to be listed
type Doctor struct {
	ID             int    `json:"id" db:"id"`
	Specialization string `json:"specialization" db:"specialization"`

	// properties stored in users struct
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Email     string `json:"email" db:"email"`
}

type DoctorSlots struct {
	SlotsDate time.Time `json:"slots_date"`
}
