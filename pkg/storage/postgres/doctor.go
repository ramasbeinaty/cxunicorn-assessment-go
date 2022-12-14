package postgres

// defines the storage format of a doctor
type Doctor struct {
	Staff
	ID             int    `json:"id" db:"id"`
	Specialization string `json:"specialization" db:"specialization"`
}

type DoctorCreate struct {
	ID             int    `json:"id" db:"id"`
	Specialization string `json:"specialization" db:"specialization"`
}
