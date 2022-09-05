package booking

import (
	"errors"
)

var Err = errors.New("failed to book an appointment")

// provide access to the appointment storage
type Repository interface {
	// creates an appointment in storage
	CreateAppointment(Appointment) error
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
	return s.repo.CreateAppointment(a)
}
