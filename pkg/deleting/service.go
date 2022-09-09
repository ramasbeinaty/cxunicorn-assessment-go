package deleting

import (
	"errors"
)

var ErrAppointment = errors.New("failed to delete appointment")

// provide access to the appointment storage
type Repository interface {
	// delete the specified appointment from storage
	DeleteAppointment(int) error
}

// provide listing operations for appointment
type Service interface {
	DeleteAppointment(int) error
}

type service struct {
	repo Repository
}

// creates a listing service with the necessary dependencies
func NewService(repo Repository) Service {
	return &service{repo}
}

// implement service methods
func (s *service) DeleteAppointment(id int) error {
	err := s.repo.DeleteAppointment(id)

	return err
}
