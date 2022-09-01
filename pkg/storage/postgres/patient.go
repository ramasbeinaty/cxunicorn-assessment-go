package postgres

// defines the storage format of a patient
type Patient struct {
	User
	ID             int32  `json:"id"`
	MedicalHistory string `json:"medical_history"`
}
