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

	// appointment must be within the work days of the doctor
	isWithinDoctorsWorkDays := s.repo.IsAppointmentWithinDoctorWorkDays(a.DoctorID, a.StartDatetime.Weekday())

	if !isWithinDoctorsWorkDays {
		return errors.New("ERROR: CreateAppointment - appointment is not within the work days of the doctor")
	}

	// appointment must be within the work time of the doctor

	//
	// appointmentStartTime := time.Date(a.StartDatetime.Year(), a.StartDatetime.Month(), a.StartDatetime.Day(),
	// a.StartDatetime.Hour(), a.StartDatetime.Minute(), a.StartDatetime.Second(), a.StartDatetime.Nanosecond(),
	// a.StartDatetime.Location())

	// appointmentEndTime := time.Date(a.EndDatetime.Year(), a.EndDatetime.Month(), a.EndDatetime.Day(),
	// a.EndDatetime.Hour(), a.EndDatetime.Minute(), a.EndDatetime.Second(), a.EndDatetime.Nanosecond(),
	// a.EndDatetime.Location())

	doctorWorkTime := s.repo.GetDoctorWorkTime(a.DoctorID)

	doctorStartWorkTime := doctorWorkTime[0]
	doctorEndWorkTime := doctorWorkTime[1]

	appointmentIsWithinDoctorWorkTime := eventIsWithinTimeBounds(a.StartDatetime, a.EndDatetime, doctorStartWorkTime, doctorEndWorkTime)

	if !appointmentIsWithinDoctorWorkTime {
		return errors.New("ERROR: CreateAppointment - appointment is not within the work timings of the doctor")
	}

	// validate that appointment is not within doctor's break time
	doctorBreakTime := s.repo.GetDoctorBreakTime(a.DoctorID)

	doctorStartBreakTime := doctorBreakTime[0]
	doctorEndBreakTime := doctorBreakTime[1]

	appointmentIsWithinDoctorBreakTime := eventIsWithinTimeBounds(a.StartDatetime, a.EndDatetime, doctorStartBreakTime, doctorEndBreakTime)

	if appointmentIsWithinDoctorBreakTime {
		return errors.New("ERROR: CreateAppointment - appointment cannot be within the break timing of the doctor")
	}

	// appointment should not conflict with other previously set appointments

	// if validations come through, add the appointment to storage
	var appointments postgres.AppointmentCreate
	appointments.PatientID = a.PatientID
	appointments.DoctorID = a.DoctorID
	appointments.CreatedBy = a.CreatedBy
	appointments.StartDatetime = a.StartDatetime
	appointments.EndDatetime = a.EndDatetime

	return s.repo.CreateAppointment(appointments)
}

func eventIsWithinTimeBounds(event_start_time time.Time, event_end_time time.Time, start_time_bound time.Time, end_time_bound time.Time) bool {
	// specified event cannot start before the time bound
	if event_start_time.Before(start_time_bound) {
		fmt.Println("ERROR: eventIsWithinTimeBounds - event cannot start before time bound")
		return false
	}

	// specified event cannot start after the time bound
	if event_start_time.After(end_time_bound) {
		fmt.Println("ERROR: eventIsWithinTimeBounds - event cannot overlap or occur after time bound")
		return false
	}

	// specified event cannot end after the time bound
	if event_end_time.After(end_time_bound) {
		fmt.Println("ERROR: eventIsWithinTimeBounds - event cannot end after time bound")
		return false
	}

	return true
}
