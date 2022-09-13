package auth

type PatientRegister struct {
	UserDetails    User   `json:"user_details"`
	MedicalHistory string `json:"medical_history"`
}
