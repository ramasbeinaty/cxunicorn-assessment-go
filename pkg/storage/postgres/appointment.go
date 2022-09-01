package postgres

// defines the storage format of an appointment
type Appointment struct {
	Event
	ID        int32 `json:"id"`
	PatientID int32 `json:"patient_id"`
	DoctorID  int32 `json:"doctor_id"`
}
