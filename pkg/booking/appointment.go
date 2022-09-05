package booking

import "time"

// defines the properties required to book an appointment
type Appointment struct {
	PatientID int `json:"patient_id"`
	DoctorID  int `json:"doctor_id"`

	CreatedBy     int       `json:"created_by"`
	StartDatetime time.Time `json:"start_datetime"`
	EndDatetime   time.Time `json:"end_datetime"`
}
