package listing

import (
	"errors"
)

var ErrIdNotFound = errors.New("Doctor with given id not found")
var ErrEmpty = errors.New("No doctor was found")

// provide access to the doctor storage
type Repository interface {
	// returns a doctor with given id
	GetDoctor(string) (Doctor, error)

	// returns all doctors in storage
	GetAllDoctors() []Doctor
}

// provide listing operations for struct doctor
type Service interface {
	GetDoctor(string) (Doctor, error)
	GetAllDoctors() []Doctor
}

type service struct {
	repo Repository
}

// creates a listing service with the necessary dependencies
func NewService(repo Repository) Service {
	return &service{repo}
}

// implement service methods
func (s *service) GetDoctor(id string) (Doctor, error) {
	return s.repo.GetDoctor(id)
}

func (s *service) GetAllDoctors() []Doctor {
	return s.repo.GetAllDoctors()
}
