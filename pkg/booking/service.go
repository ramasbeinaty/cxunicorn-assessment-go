package booking

import (
	"clinicapp/pkg/storage/postgres"
	"errors"
	"fmt"
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

	// to validate appointment is within doctor's work days
	IsAppointmentWithinDoctorWorkDays(int, time.Weekday) bool

	// to validate appointment is within doctor's work hours
	GetDoctorWorkTime(int) []time.Time

	// to validate appointment is not within doctor's break time
	GetDoctorBreakTime(int) []time.Time

	// to validate appointment does not overlap previously booked appointments of neither the doctor's nor the patient's
	IsAppointmentOverlapping(int, int, time.Time, time.Time) bool
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
	datetimeDifference := a.EndDatetime.Sub(a.StartDatetime)

	// validate the appointment is at least 15 minutes
	if datetimeDifference.Minutes() < 15 {
		return errors.New("ERROR: CreateAppointment - appointment duration cannot less than 15 minutes")
	}

	// validate the appointment is at most 2 hours
	if datetimeDifference.Hours() > 2 {
		return errors.New("ERROR: CreateAppointment - appointment duration cannot greater than 2 hours")
	}

	// validate doctor does not have more than 12 different patients per day
	numberOfAppointments := s.repo.GetNumberOfAppointmentsWithDistinctPatient(a.DoctorID, a.StartDatetime)
	if numberOfAppointments > 12 {
		return errors.New("ERROR: CreateAppointment - doctor cannot have more than 12 different patients in a day")
	}

	// validate doctor does not have more than 8 hours per day
	appointmentHours := s.repo.GetAppointmentHoursPerDay(a.DoctorID, a.StartDatetime)
	if appointmentHours > 8 {
		return errors.New("ERROR: CreateAppointment - doctor cannot have more than 8 hours in a day")
	}

	// validate appointment is within the work days of the doctor
	isWithinDoctorsWorkDays := s.repo.IsAppointmentWithinDoctorWorkDays(a.DoctorID, a.StartDatetime.Weekday())

	if !isWithinDoctorsWorkDays {
		return errors.New("ERROR: CreateAppointment - appointment is not within the work days of the doctor")
	}

	// validate appointment is within the work time of the doctor

	doctorWorkTime := s.repo.GetDoctorWorkTime(a.DoctorID)

	doctorStartWorkTime := doctorWorkTime[0]
	doctorEndWorkTime := doctorWorkTime[1]

	// declare and initialize appointment time only so dates are not taken into consideration when comparing timestamps
	appointmentStartTime := time.Date(doctorEndWorkTime.Year(), doctorEndWorkTime.Month(), doctorEndWorkTime.Day(),
		a.StartDatetime.Hour(), a.StartDatetime.Minute(), a.StartDatetime.Second(), a.StartDatetime.Nanosecond(),
		a.StartDatetime.Location())

	appointmentEndTime := time.Date(doctorEndWorkTime.Year(), doctorEndWorkTime.Month(), doctorEndWorkTime.Day(),
		a.EndDatetime.Hour(), a.EndDatetime.Minute(), a.EndDatetime.Second(), a.EndDatetime.Nanosecond(),
		a.EndDatetime.Location())

	appointmentIsWithinDoctorWorkTime := appointmentIsWithinWorkTime(appointmentStartTime, appointmentEndTime, doctorStartWorkTime, doctorEndWorkTime)

	appointmentIsWithinBreakTime(appointmentStartTime, appointmentEndTime, doctorStartWorkTime, doctorEndWorkTime)

	if !appointmentIsWithinDoctorWorkTime {
		return errors.New("ERROR: CreateAppointment - appointment is not within the work timings of the doctor")
	}

	// validate that appointment is not within doctor's break time
	doctorBreakTime := s.repo.GetDoctorBreakTime(a.DoctorID)

	doctorStartBreakTime := doctorBreakTime[0]
	doctorEndBreakTime := doctorBreakTime[1]

	appointmentIsWithinBreakTime := appointmentIsWithinBreakTime(appointmentStartTime, appointmentEndTime, doctorStartBreakTime, doctorEndBreakTime)

	if appointmentIsWithinBreakTime {
		return errors.New("ERROR: CreateAppointment - appointment cannot occur within the break time of doctor")
	}

	// validate appointment doesn't overlap with previously booked appointments of both the doctor's and the patient's
	isOverlapping := s.repo.IsAppointmentOverlapping(a.DoctorID, a.PatientID, a.StartDatetime, a.EndDatetime)

	if isOverlapping {
		return errors.New("ERROR: CreateAppointment - appointment overlaps a previously booked appointment")
	}

	// if validations come through, add the appointment to storage
	var appointments postgres.AppointmentCreate
	appointments.PatientID = a.PatientID
	appointments.DoctorID = a.DoctorID
	appointments.CreatedBy = a.CreatedBy
	appointments.StartDatetime = a.StartDatetime
	appointments.EndDatetime = a.EndDatetime

	return s.repo.CreateAppointment(appointments)
}

func appointmentIsWithinWorkTime(appointmentStartTime time.Time, appointmentEndTime time.Time, workStartTime time.Time, workEndTime time.Time) bool {
	if appointmentStartTime.Before(workStartTime) {
		fmt.Println("INFO: appointmentIsWithinWorkTime - appointment starts before doctor's work time")
		return false
	}

	if appointmentStartTime.After(workEndTime) {
		fmt.Println("INFO: appointmentIsWithinWorkTime - appointment overlaps or occurs after doctor's work time")
		return false
	}

	if appointmentEndTime.After(workEndTime) {
		fmt.Println("INFO: appointmentIsWithinWorkTime - appointment ends after doctor's work time")
		return false
	}

	return true
}

func appointmentIsWithinBreakTime(appointmentStartTime time.Time, appointmentEndTime time.Time, breakStartTime time.Time, breakEndTime time.Time) bool {
	if appointmentStartTime.After(breakEndTime) {
		fmt.Println("INFO: appointmentIsWithinBreakTime - appointment starts after the break ends")
		return false
	}

	if appointmentStartTime.Before(breakStartTime) && appointmentEndTime.Before(breakStartTime) {
		fmt.Println("INFO: appointmentIsWithinBreakTime - appointment starts and ends before the break starts")
		return false
	}

	return true
}
