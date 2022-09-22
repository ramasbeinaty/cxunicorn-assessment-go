package auth

import (
	"clinicapp/pkg/storage/postgres"
	"encoding/json"
	"errors"

	"github.com/golang-jwt/jwt"
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
	CreateUser(UserRegister) (string, error)
	CreatePatient(PatientRegister) error
	CreateDoctor(DoctorRegister) error
	CreateClinicAdmin(ClinicAdminRegister) error

	LoginUser(UserLogin) (string, error)
	GenerateJWT(*Claims) (string, error)
	AuthenticateUser(UserLogin, User) error
	AuthorizeUser(string, string) (bool, error)
	GetTokenFromString(string, *Claims) (*jwt.Token, error)
	VerifyJWT(string) (bool, *Claims)
}

type service struct {
	repo Repository
}

// creates a listing service with the necessary dependencies
func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) LoginUser(loginCredentials UserLogin) (string, error) {
	// returns the user, their token and error
	var _user postgres.User
	var user User
	var tokenStr = ""

	// find user
	_user, err := s.repo.GetUser(loginCredentials.Email)
	if err != nil {
		return tokenStr, errors.New("LoginUser - Failed to find user - " + err.Error())
	}

	// map user
	user.Email = _user.Email
	user.FirstName = _user.FirstName
	user.LastName = _user.LastName
	user.Password = _user.Password
	user.PhoneNumber = _user.PhoneNumber
	user.Role = _user.Role
	user.DOB = _user.DOB
	user.ID = _user.ID

	// authenticate user
	err = s.AuthenticateUser(loginCredentials, user)
	if err != nil {
		return tokenStr, errors.New("LoginUser - Failed to authenticate user - " + err.Error())
	}

	// define the custom token claims
	var claims = &Claims{}
	claims.UserID = user.ID
	claims.Email = user.Email
	claims.Role = user.Role
	claims.Name = user.FirstName + " " + user.LastName

	// generate jwt token
	tokenStr, err = s.GenerateJWT(claims)

	if err != nil {
		return tokenStr, errors.New("LoginUser - Failed to generate jwt token - " + err.Error())
	}

	return tokenStr, nil
}

func (s *service) AuthenticateUser(loginCredentials UserLogin, user User) error {

	// decrypt received password and check if equal
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginCredentials.Password))

	// check if password matches
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return errors.New("ERROR: AuthenticateUser - Incorrect password -" + err.Error())
	}

	return nil
}

func (s *service) AuthorizeUser(authorizedRole string, userRole string) (bool, error) {
	if authorizedRole == "" {
		return false, errors.New("AuthorizeUser - authorized role field cannot be empty")
	}

	if userRole == "" {
		return false, errors.New("AuthorizeUser - user role field cannot be empty")
	}

	if authorizedRole != Roles.ClinicAdmin && authorizedRole != Roles.Patient && authorizedRole != Roles.Doctor {
		return false, errors.New("AuthorizeUser - authorized role field is not a valid role")
	}

	if userRole != Roles.ClinicAdmin && userRole != Roles.Patient && userRole != Roles.Doctor {
		return false, errors.New("AuthorizeUser - user role field is not a valid role")
	}

	if authorizedRole != userRole {
		return false, errors.New("AuthorizeUser - user role does not match authorized role")
	}

	return true, nil
}

// implement service methods
func (s *service) CreateUser(user UserRegister) (string, error) {
	var tokenStr string = ""

	// check if user exists
	_user, _ := s.repo.GetUser(user.UserDetails.Email)
	if _user != (postgres.User{}) {
		return tokenStr, errors.New("ERROR: CreateUser - user with email '" + user.UserDetails.Email + "' already exists. Try to login instead.")
	}

	// define the custom token claims
	var claims = &Claims{}
	claims.UserID = user.UserDetails.ID
	claims.Email = user.UserDetails.Email
	claims.Role = user.UserDetails.Role
	claims.Name = user.UserDetails.FirstName + " " + user.UserDetails.LastName

	// generate jwt token
	tokenStr, err := s.GenerateJWT(claims)

	if err != nil {
		return tokenStr, errors.New("ERROR: LoginUser - Failed to generate jwt token - " + err.Error())
	}

	// parse user to the role's corresponding model - patient, doctor or clinic admin
	// then call the corresponding functions to create those models and commit to db
	if user.UserDetails.Role == Roles.Patient {
		var _patient PatientRegister
		json.Unmarshal([]byte(user.RoleDetails), &_patient)

		_patient.UserDetails.FirstName = user.UserDetails.FirstName
		_patient.UserDetails.LastName = user.UserDetails.LastName
		_patient.UserDetails.DOB = user.UserDetails.DOB
		_patient.UserDetails.PhoneNumber = user.UserDetails.PhoneNumber
		_patient.UserDetails.Email = user.UserDetails.Email
		_patient.UserDetails.Password = user.UserDetails.Password
		_patient.UserDetails.Role = user.UserDetails.Role

		return tokenStr, s.CreatePatient(_patient)

	} else if user.UserDetails.Role == Roles.Doctor {
		var _doctor DoctorRegister
		json.Unmarshal([]byte(user.RoleDetails), &_doctor)

		_doctor.UserDetails.FirstName = user.UserDetails.FirstName
		_doctor.UserDetails.LastName = user.UserDetails.LastName
		_doctor.UserDetails.DOB = user.UserDetails.DOB
		_doctor.UserDetails.PhoneNumber = user.UserDetails.PhoneNumber
		_doctor.UserDetails.Email = user.UserDetails.Email
		_doctor.UserDetails.Password = user.UserDetails.Password
		_doctor.UserDetails.Role = user.UserDetails.Role

		return tokenStr, s.CreateDoctor(_doctor)

	} else if user.UserDetails.Role == Roles.ClinicAdmin {
		var _clinicAdmin ClinicAdminRegister
		json.Unmarshal([]byte(user.RoleDetails), &_clinicAdmin)

		_clinicAdmin.UserDetails.LastName = user.UserDetails.LastName
		_clinicAdmin.UserDetails.DOB = user.UserDetails.DOB
		_clinicAdmin.UserDetails.PhoneNumber = user.UserDetails.PhoneNumber
		_clinicAdmin.UserDetails.Email = user.UserDetails.Email
		_clinicAdmin.UserDetails.Password = user.UserDetails.Password
		_clinicAdmin.UserDetails.Role = user.UserDetails.Role
		_clinicAdmin.UserDetails.FirstName = user.UserDetails.FirstName

		return tokenStr, s.CreateClinicAdmin(_clinicAdmin)

	} else {
		return tokenStr, errors.New("ERROR: CreateUser - role '" + user.UserDetails.Role + "' does not exist")
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
