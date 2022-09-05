package postgres

import "time"

// defines the storage format of an appointment
type Appointment struct {
	ID            int       `json:"id"`
	PatientID     int       `json:"patient_id"`
	DoctorID      int       `json:"doctor_id"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     int       `json:"created_by"`
	StartDatetime time.Time `json:"start_datetime"`
	EndDatetime   time.Time `json:"end_datetime"`
	IsCanceled    bool      `json:"is_canceled"`
}

type AppointmentCreate struct {
	PatientID int `json:"patient_id"`
	DoctorID  int `json:"doctor_id"`

	CreatedBy     int       `json:"created_by"`
	StartDatetime time.Time `json:"start_datetime"`
	EndDatetime   time.Time `json:"end_datetime"`
}

// func NewAppointment(patientID int, doctorID int, createdBy int,
// 	startDatetime time.Time, endDatetime time.Time) appointment {
// 	_appointment := appointment{}
// 	_appointment.PatientID = patientID
// 	_appointment.DoctorID = doctorID

// 	_appointment.CreatedAt = time.Now()
// 	_appointment.CreatedBy = createdBy
// 	_appointment.StartDatetime = startDatetime
// 	_appointment.EndDatetime = endDatetime
// 	_appointment.IsCanceled = false

// 	return _appointment
// }
