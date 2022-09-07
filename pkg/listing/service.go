package listing

import (
	"clinicapp/pkg/storage/postgres"
	"errors"
)

var ErrIdNotFound = errors.New("doctor with given id not found")
var ErrEmpty = errors.New("no doctor was found")

// provide access to the doctor storage
type Repository interface {
	// returns a doctor with given id
	GetDoctor(int) (postgres.Doctor, error)

	// returns all doctors in storage
	// GetAllDoctors() []Doctor
}

// provide listing operations for struct doctor
type Service interface {
	GetDoctor(int) (Doctor, error)
	// GetAllDoctors() []Doctor
}

type service struct {
	repo Repository
}

// creates a listing service with the necessary dependencies
func NewService(repo Repository) Service {
	return &service{repo}
}

// implement service methods
func (s *service) GetDoctor(id int) (Doctor, error) {
	var d postgres.Doctor
	var doctor Doctor
	var err error

	d, err = s.repo.GetDoctor(id)

	doctor.ID = d.ID
	doctor.Email = d.Email
	doctor.FirstName = d.FirstName
	doctor.LastName = d.LastName
	doctor.Specialization = d.Specialization

	if err != nil {
		return doctor, errors.New("GetDoctor - " + err.Error())
	}

	return doctor, nil
}

// func (s *service) GetAllDoctors() []Doctor {
// 	return s.repo.GetAllDoctors()
// }
