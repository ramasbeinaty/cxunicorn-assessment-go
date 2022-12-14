package canceling

import (
	"errors"
)

var ErrIdNotFound = errors.New("doctor with given id not found")
var ErrEmpty = errors.New("no doctor was found")

// provide access to the appointment storage
type Repository interface {
	// modify the storage to set the appointment field is_canceled = True
	CancelAppointment(int) error
}

// provide listing operations for appointment
type Service interface {
	CancelAppointment(int) error
}

type service struct {
	repo Repository
}

// creates a listing service with the necessary dependencies
func NewService(repo Repository) Service {
	return &service{repo}
}

// implement service methods
func (s *service) CancelAppointment(id int) error {
	err := s.repo.CancelAppointment(id)

	return err
}
