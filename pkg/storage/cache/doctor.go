package cache

const cachedDoctor = "doctor"

// defines the storage format of a doctor
type Doctor struct {
	Staff
	ID             int    `json:"id"`
	Specialization string `json:"specialization"`
}