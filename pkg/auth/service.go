package auth

import (
	"clinicapp/pkg/storage/postgres"
	"encoding/json"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrAppointment = errors.New("failed to edit appointment")

// provide access to the users, staffs, patients, doctors and clinic admins storage
type Repository interface {
	CreateUser(postgres.UserCreate) (int, error)
	CreatePatient(postgres.PatientCreate) error
	CreateDoctor(postgres.DoctorCreate) error
	CreateClinicAdmin(postgres.ClinicAdminCreate) error
	CreateStaff(postgres.StaffCreate) error

	GetUser(string) (postgres.User, error)
}

// provide listing operations for authenticating/authorizing users
type Service interface {
	CreateUser(UserRegister) error
	CreatePatient(PatientRegister) error
	CreateDoctor(DoctorRegister) error
	CreateClinicAdmin(ClinicAdminRegister) error

	LoginUser(UserLogin) (User, string, error)
}

type service struct {
	repo Repository
}

// creates a listing service with the necessary dependencies
func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) LoginUser(loginCredentials UserLogin) (User, string, error) {
	// returns the user, their token and error
	var _user postgres.User
	var user User
	var token string

	_user, err := s.repo.GetUser(loginCredentials.Email)

	println(_user)

	if err != nil {
		return user, "", errors.New("ERROR: LoginUser - " + err.Error())
	}

	// decrypt received password and check if equal

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginCredentials.Password))

	// check if password matches
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return user, "", errors.New("ERROR: LoginUser - Incorrect password -" + err.Error())
	}

	// generate jwt token
	// expiresAt := time.Now().Add(time.Minute * 100000).Unix()

	// tk := user.Token{
	// 	UserID: user.ID,
	// 	Name:   user.Name,
	// 	Email:  user.Email,
	// 	StandardClaims: &jwt.StandardClaims{
	// 		ExpiresAt: expiresAt,
	// 	},
	// }

	// token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	// tokenString, error := token.SignedString([]byte("secret"))
	// if error != nil {
	// 	fmt.Println(error)
	// }

	return user, token, nil

}

// implement service methods
func (s *service) CreateUser(user UserRegister) error {

	// parse user to the role's corresponding model - patient, doctor or clinic admin
	// then call the corresponding functions to create those models and commit to db
	if user.UserDetails.Role == postgres.PatientRole {
		var _patient PatientRegister
		json.Unmarshal([]byte(user.RoleDetails), &_patient)

		_patient.UserDetails.FirstName = user.UserDetails.FirstName
		_patient.UserDetails.LastName = user.UserDetails.LastName
		_patient.UserDetails.DOB = user.UserDetails.DOB
		_patient.UserDetails.PhoneNumber = user.UserDetails.PhoneNumber
		_patient.UserDetails.Email = user.UserDetails.Email
		_patient.UserDetails.Password = user.UserDetails.Password
		_patient.UserDetails.Role = user.UserDetails.Role

		return s.CreatePatient(_patient)

	} else if user.UserDetails.Role == postgres.DoctorRole {
		var _doctor DoctorRegister
		json.Unmarshal([]byte(user.RoleDetails), &_doctor)

		_doctor.UserDetails.FirstName = user.UserDetails.FirstName
		_doctor.UserDetails.LastName = user.UserDetails.LastName
		_doctor.UserDetails.DOB = user.UserDetails.DOB
		_doctor.UserDetails.PhoneNumber = user.UserDetails.PhoneNumber
		_doctor.UserDetails.Email = user.UserDetails.Email
		_doctor.UserDetails.Password = user.UserDetails.Password
		_doctor.UserDetails.Role = user.UserDetails.Role

		return s.CreateDoctor(_doctor)

	} else if user.UserDetails.Role == postgres.ClinicAdminRole {
		var _clinicAdmin ClinicAdminRegister
		json.Unmarshal([]byte(user.RoleDetails), &_clinicAdmin)

		_clinicAdmin.UserDetails.LastName = user.UserDetails.LastName
		_clinicAdmin.UserDetails.DOB = user.UserDetails.DOB
		_clinicAdmin.UserDetails.PhoneNumber = user.UserDetails.PhoneNumber
		_clinicAdmin.UserDetails.Email = user.UserDetails.Email
		_clinicAdmin.UserDetails.Password = user.UserDetails.Password
		_clinicAdmin.UserDetails.Role = user.UserDetails.Role
		_clinicAdmin.UserDetails.FirstName = user.UserDetails.FirstName

		return s.CreateClinicAdmin(_clinicAdmin)

	} else {
		return errors.New("ERROR: CreateUser - role '" + user.UserDetails.Role + "' does not exist")
	}

}

func (s *service) CreatePatient(patient PatientRegister) error {
	var _user postgres.UserCreate
	var _patient postgres.PatientCreate

	// user details
	_user.FirstName = patient.UserDetails.FirstName
	_user.LastName = patient.UserDetails.LastName
	_user.DOB = patient.UserDetails.DOB
	_user.PhoneNumber = patient.UserDetails.PhoneNumber
	_user.Email = patient.UserDetails.Email
	_user.Password = patient.UserDetails.Password
	_user.Role = patient.UserDetails.Role

	// get the generated user id to use it for creating a patient and achieve inheritance
	user_id, err := s.repo.CreateUser(_user)

	if err != nil {
		return errors.New("ERROR: CreatePatient - " + err.Error())
	}

	// patient details
	_patient.ID = user_id
	_patient.MedicalHistory = patient.MedicalHistory

	return s.repo.CreatePatient(_patient)
}

func (s *service) CreateDoctor(doctor DoctorRegister) error {
	var _user postgres.UserCreate
	var _staff postgres.StaffCreate
	var _doctor postgres.DoctorCreate

	// user details
	_user.FirstName = doctor.UserDetails.FirstName
	_user.LastName = doctor.UserDetails.LastName
	_user.DOB = doctor.UserDetails.DOB
	_user.PhoneNumber = doctor.UserDetails.PhoneNumber
	_user.Email = doctor.UserDetails.Email
	_user.Password = doctor.UserDetails.Password
	_user.Role = doctor.UserDetails.Role

	// get the generated user id to use it for creating a doctor and achieve inheritance
	user_id, err := s.repo.CreateUser(_user)

	if err != nil {
		return errors.New("ERROR: CreateDoctor - " + err.Error())
	}

	// staff details
	_staff.ID = user_id
	_staff.WorkDays = doctor.WorkDays
	_staff.WorkTime = doctor.WorkTime
	_staff.BreakTime = doctor.BreakTime

	if err = s.repo.CreateStaff(_staff); err != nil {
		return err
	}

	// doctor details
	_doctor.ID = user_id
	_doctor.Specialization = doctor.Specialization

	return s.repo.CreateDoctor(_doctor)
}

func (s *service) CreateClinicAdmin(clinicAdmin ClinicAdminRegister) error {
	var _user postgres.UserCreate
	var _staff postgres.StaffCreate
	var _clinicAdmin postgres.ClinicAdminCreate

	// user details
	_user.FirstName = clinicAdmin.UserDetails.FirstName
	_user.LastName = clinicAdmin.UserDetails.LastName
	_user.DOB = clinicAdmin.UserDetails.DOB
	_user.PhoneNumber = clinicAdmin.UserDetails.PhoneNumber
	_user.Email = clinicAdmin.UserDetails.Email
	_user.Password = clinicAdmin.UserDetails.Password
	_user.Role = clinicAdmin.UserDetails.Role

	// get the generated user id to use it for creating a doctor and achieve inheritance
	user_id, err := s.repo.CreateUser(_user)

	if err != nil {
		return errors.New("ERROR: CreateDoctor - " + err.Error())
	}

	// staff details
	_staff.ID = user_id
	_staff.WorkDays = clinicAdmin.WorkDays
	_staff.WorkTime = clinicAdmin.WorkTime
	_staff.BreakTime = clinicAdmin.BreakTime

	if err = s.repo.CreateStaff(_staff); err != nil {
		return err
	}

	// clinic admin details
	_clinicAdmin.ID = user_id

	return s.repo.CreateClinicAdmin(_clinicAdmin)
}
