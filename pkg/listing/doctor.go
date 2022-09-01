package listing

// defines the doctor properties to be listed
type Doctor struct {
	ID             int32  `json:"id" db:"id"`
	Specialization string `json:"specialization" db:"specialization"`

	// properties stored in users struct
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Email     string `json:"email" db:"email"`
	WorkShift string `json:"workshift" db:"work_shift"`
}
