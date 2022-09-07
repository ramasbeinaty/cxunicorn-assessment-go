package booking

import (
	"errors"
)

var Err = errors.New("failed to book an appointment")

// provide access to the appointment storage
type Repository interface {
	// creates an appointment in storage
	CreateAppointment(Appointment) error

	// validate doctor exists
	DoctorExists(int) bool
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

	// validate the appointment end datetime is after the start datetime
	if a.EndDatetime.Before(a.StartDatetime) {
		return errors.New("ERROR: CreateAppointment - end datetime cannot be before start datetime")
	}

	// appointment must be within the work time of the doctor

	// appointment must be

	// if validations come through, add the appointment to storage

	return s.repo.CreateAppointment(a)
}
