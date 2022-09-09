package editing

import (
	"clinicapp/pkg/storage/postgres"
	"errors"
)

var ErrAppointment = errors.New("failed to edit appointment")

// provide access to the appointment storage
type Repository interface {
	// edit the specified appointment from storage
	EditAppointment(int, postgres.AppointmentEdit) error

	// get the specified appointment from storage
	GetAppointment(int) postgres.Appointment
}

// provide listing operations for appointment
type Service interface {
	EditAppointment(int, Appointment) error
}

type service struct {
	repo Repository
}

// creates a listing service with the necessary dependencies
func NewService(repo Repository) Service {
	return &service{repo}
}

// implement service methods
func (s *service) EditAppointment(id int, appointment Appointment) error {

	// get all data fields of specified appointment from the database
	// _appointment := s.repo.GetAppointment(id)

	var _appointment postgres.AppointmentEdit

	_appointment.StartDatetime = appointment.StartDatetime
	_appointment.EndDatetime = appointment.EndDatetime
	_appointment.IsCanceled = appointment.IsCanceled

	err := s.repo.EditAppointment(id, _appointment)

	return err
}
