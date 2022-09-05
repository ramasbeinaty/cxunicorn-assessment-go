package postgres

// defines the storage format of a patient
type Patient struct {
	User
	ID             int  `json:"id"`
	MedicalHistory string `json:"medical_history"`
}
