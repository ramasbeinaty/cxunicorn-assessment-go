package booking

import (
	"clinicapp/pkg/storage/postgres"
	"errors"
	"time"
)

var Err = errors.New("failed to book an appointment")

// provide access to the appointment storage
type Repository interface {
	// creates an appointment in storage
	CreateAppointment(postgres.AppointmentCreate) error

	// validate doctor exists
	DoctorExists(int) bool

	// to validate doctor does not have more than 12 appoointments with different patients per day
	GetNumberOfAppointmentsWithDistinctPatient(int, time.Time) int

	// to validate doctors have less than 8 hours of appointments per day
	GetAppointmentHoursPerDay(int, time.Time) int
}

// provide booking operations for struct appointment
type Service interface {
	CreateAppointment(Appointment) error
}

type service struct {
	repo Repository
}

// creates a booking service with the necessary dependencies
func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) CreateAppointment(a Appointment) error {

	// json.Unmarshal([]byte(Appointment))

	// validate that doctor id exists
	doctorExists := s.repo.DoctorExists(a.DoctorID)

	if !doctorExists {
		return errors.New("ERROR: CreateAppointment - doctor id does not exist. ")
	}

	// validate start date time is after now's date time
	if a.StartDatetime.Before(time.Now()) {
		return errors.New("ERROR: CreateAppointment - start date time cannot be before now date time")
	}

	// validate the appointment end date time is after the start datetime
	if a.EndDatetime.Before(a.StartDatetime) {
		return errors.New("ERROR: CreateAppointment - end date time cannot be before start date time")
	}

	// calculate the start and end date time difference
	datetime_difference := a.EndDatetime.Sub(a.StartDatetime)

	// validate the appointment is at least 15 minutes
	if datetime_difference.Minutes() < 15 {
		return errors.New("ERROR: CreateAppointment - appointment duration cannot less than 15 minutes")
	}

	// validate the appointment is at most 2 hours
	if datetime_difference.Hours() > 2 {
		return errors.New("ERROR: CreateAppointment - appointment duration cannot greater than 2 hours")
	}

	// validate doctor does not have more than 12 different patients per day
	number_of_appointments := s.repo.GetNumberOfAppointmentsWithDistinctPatient(a.DoctorID, a.StartDatetime)
	if number_of_appointments > 12 {
		return errors.New("ERROR: CreateAppointment - doctor cannot have more than 12 different patients in a day")
	}

	// validate doctor does not have more than 8 hours per day
	appointment_hours := s.repo.GetAppointmentHoursPerDay(a.DoctorID, a.StartDatetime)
	if appointment_hours > 8 {
		return errors.New("ERROR: CreateAppointment - doctor cannot have more than 8 hours in a day")
	}

	// appointment must be within the work time of the doctor
	

	// appointment can not be during doctor's break time

	// appointment should not conflict with others

	// if validations come through, add the appointment to storage
	var appointments postgres.AppointmentCreate
	appointments.PatientID = a.PatientID
	appointments.DoctorID = a.DoctorID
	appointments.CreatedBy = a.CreatedBy
	appointments.StartDatetime = a.StartDatetime
	appointments.EndDatetime = a.EndDatetime

	return s.repo.CreateAppointment(appointments)
}
