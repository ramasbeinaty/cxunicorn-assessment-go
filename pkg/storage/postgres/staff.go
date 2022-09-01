package postgres

// defines the storage format of a staff
type Staff struct {
	User
	ID        int32  `json:"id" db:"id"`
	WorkShift string `json:"work_shift" db:"work_shift"`
}
