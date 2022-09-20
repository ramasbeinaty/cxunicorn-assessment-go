package listing

import (
	"clinicapp/pkg/storage/postgres"
	"errors"
	"fmt"
	"time"
)

var ErrIdNotFound = errors.New("doctor with given id not found")
var ErrEmpty = errors.New("no doctor was found")

// provide access to the doctor storage
type Repository interface {
	// returns a doctor with given id
	GetDoctor(int) (postgres.Doctor, error)

	// returns all doctors in storage
	GetAllDoctors() []postgres.Doctor

	// returns all the appointments of a specific doctor
	GetAllAppointmentsOfDoctor(int, time.Time) []postgres.Appointment
}

// provide listing operations for struct doctor
type Service interface {
	GetDoctor(int) (Doctor, error)
	GetAllDoctors() []Doctor
	GetAllAppointmentsOfDoctor(int, time.Time) []Appointment
	GetAvailableSlotsPerDay(int, time.Time) [][]time.Time
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

func (s *service) GetAllDoctors() []Doctor {
	var _doctors []postgres.Doctor
	var doctors []Doctor = []Doctor{}

	_doctors = s.repo.GetAllDoctors()

	for _, _doctor := range _doctors {
		var doctor Doctor

		doctor.ID = _doctor.ID
		doctor.Email = _doctor.Email
		doctor.FirstName = _doctor.FirstName
		doctor.LastName = _doctor.LastName
		doctor.Specialization = _doctor.Specialization

		doctors = append(doctors, doctor)
	}

	return doctors
}

func (s *service) GetAllAppointmentsOfDoctor(doctorID int, date time.Time) []Appointment {
	var _appointments []postgres.Appointment
	var appointments []Appointment = []Appointment{}

	_appointments = s.repo.GetAllAppointmentsOfDoctor(doctorID, date)

	for _, _appointment := range _appointments {
		var appointment Appointment

		appointment.ID = _appointment.ID
		appointment.DoctorID = _appointment.DoctorID
		appointment.CreatedBy = _appointment.CreatedBy
		appointment.CreatedAt = _appointment.CreatedAt
		appointment.StartDatetime = _appointment.StartDatetime
		appointment.EndDatetime = _appointment.EndDatetime

		appointments = append(appointments, appointment)
	}

	return appointments
}

func (s *service) GetAvailableSlotsPerDay(doctorID int, slotDate time.Time) [][]time.Time {
	var _appointments []postgres.Appointment
	var availableSlots [][]time.Time = [][]time.Time{}

	// var availableSlotsBeforeBreakTime [][]time.Time = [][]time.Time{}

	_doctor, err := s.repo.GetDoctor(doctorID)

	if err != nil {
		fmt.Println("ERROR: GetAvailableSlotsPerDay - Failed to find specified doctor")
		return availableSlots
	}

	// timezoneLocation := time.LoadLocation("UTC")

	// re-declare and initialize doctor work and break times so dates are not taken into consideration when comparing timestamps
	_doctorWorkTime := []time.Time{
		// doctor's start work time
		time.Date(slotDate.Year(), slotDate.Month(), slotDate.Day(),
			_doctor.WorkTime[0].Hour(), _doctor.WorkTime[0].Minute(), _doctor.WorkTime[0].Second(), _doctor.WorkTime[0].Nanosecond(),
			_doctor.WorkTime[0].Location()),

		// doctor's end work time
		time.Date(slotDate.Year(), slotDate.Month(), slotDate.Day(),
			_doctor.WorkTime[1].Hour(), _doctor.WorkTime[1].Minute(), _doctor.WorkTime[1].Second(), _doctor.WorkTime[1].Nanosecond(),
			_doctor.WorkTime[1].Location()),
	}

	_doctorBreakTime := []time.Time{
		// doctor's start work time
		time.Date(slotDate.Year(), slotDate.Month(), slotDate.Day(),
			_doctor.BreakTime[0].Hour(), _doctor.BreakTime[0].Minute(), _doctor.BreakTime[0].Second(), _doctor.BreakTime[0].Nanosecond(),
			_doctor.BreakTime[0].Location()),

		// doctor's end work time
		time.Date(slotDate.Year(), slotDate.Month(), slotDate.Day(),
			_doctor.BreakTime[1].Hour(), _doctor.BreakTime[1].Minute(), _doctor.BreakTime[1].Second(), _doctor.BreakTime[1].Nanosecond(),
			_doctor.BreakTime[1].Location()),
	}

	_appointments = s.repo.GetAllAppointmentsOfDoctor(doctorID, slotDate)

	// BASE CASE of available slots when no appointments are booked
	// if doctor does not any appointments schedules, they would have 2 available slot ranges
	// which is starting from work time till break time starts,
	// and from when the break time ends till the end of work time
	if len(_appointments) < 1 {
		availableSlots = append(availableSlots,
			[]time.Time{_doctorWorkTime[0], _doctorBreakTime[0]},
			[]time.Time{_doctorBreakTime[1], _doctorWorkTime[1]},
		)

		return availableSlots
	}

	// FIND AVAILABLE TIME SLOTS BEFORE BREAK TIME
	// when at least one appointment is booked
	// append the ascending ordered appointments until the start of appointment
	// is not after the break time start
	// if it is, append the break end time and exit the loop
	// (ending the first half of the day, before break time)

	// APPEND FIRST AVAILABLE SLOT
	// if first appointment does start at the same time as work time starts and is before break time,
	// append the slot between [start of work, start of appointment]
	var counter int = 0

	if !_appointments[counter].StartDatetime.UTC().Equal(_doctorWorkTime[0]) &&
		_appointments[counter].StartDatetime.UTC().Before(_doctorBreakTime[0]) {
		availableSlots = append(availableSlots,
			[]time.Time{_doctorWorkTime[0], _appointments[counter].StartDatetime.UTC()})
	}

	// APPEND THE REST OF APPOINTMENTS BEFORE BREAK TIME
	appendedBreakTimeStart := false

	for {

		if counter+1 >= len(_appointments) {
			break
		}

		// make sure the appointment lies in the first half (before break time), otherwise break
		if _appointments[counter+1].StartDatetime.UTC().After(_doctorBreakTime[0]) {
			break
		}

		availableSlots = append(availableSlots,
			[]time.Time{_appointments[counter].EndDatetime.UTC(),
				_appointments[counter+1].StartDatetime.UTC()})

		counter++

	}

	var appendedWorkTime = false

	// APPEND REST OF APPOINTMENTS AFTER BREAK TIME
	// make sure start of break time is appended
	// and append the end of break time to start populating the second half
	if !appendedBreakTimeStart {
		availableSlots = append(availableSlots,
			[]time.Time{_appointments[counter].EndDatetime.UTC(), _doctorBreakTime[0]})

		if counter+1 >= len(_appointments) {
			availableSlots = append(availableSlots,
				[]time.Time{_doctorBreakTime[1], _doctorWorkTime[1]})

			appendedWorkTime = true

		} else {
			counter++

			availableSlots = append(availableSlots,
				[]time.Time{_doctorBreakTime[1], _appointments[counter].StartDatetime.UTC()})
		}

	}

	// append appointments in the second half (after break time)
	for {

		// check if the last appointment booked is reached
		// if so, add the last possible slot of the date (end date time of appointment, end date time of work)
		// if the last appointment does not end at the same time of work
		// also make sure break time is taken into consideration
		if counter+1 >= len(_appointments) {
			// if appendedBreakTimeStart && !appendedBreakTimeEnd {
			// 	availableSlots = append(availableSlots,
			// 		[]time.Time{_doctorBreakTime[1], _doctorWorkTime[1]})

			// 	break
			// }

			if !appendedWorkTime && !_appointments[counter].EndDatetime.UTC().Equal(_doctorWorkTime[1]) {
				availableSlots = append(availableSlots,
					[]time.Time{_appointments[counter].EndDatetime.UTC(), _doctorWorkTime[1]})
			}

			break

		}

		availableSlots = append(availableSlots,
			[]time.Time{_appointments[counter].EndDatetime.UTC(),
				_appointments[counter+1].StartDatetime.UTC()})

		counter++
	}

	return availableSlots

}
