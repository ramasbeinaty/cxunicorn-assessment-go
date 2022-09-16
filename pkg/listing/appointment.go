package listing

import "time"

// defines the properties required to book an appointment
type Appointment struct {
	ID        int `json:"id"`
	PatientID int `json:"patient_id"`
	DoctorID  int `json:"doctor_id"`

	CreatedBy     int       `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
	StartDatetime time.Time `json:"start_datetime"`
	EndDatetime   time.Time `json:"end_datetime"`
}

type AppointmentsRequest struct {
	Date time.Time `json:"date"`
}
