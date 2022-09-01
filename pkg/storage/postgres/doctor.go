package postgres

// defines the storage format of a doctor
type Doctor struct {
	Staff
	ID             int32  `json:"id" db:"id"`
	Specialization string `json:"specialization" db:"specialization"`
}
