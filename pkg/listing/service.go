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

	_doctor, err := s.repo.GetDoctor(doctorID)

	if err != nil {
		fmt.Println("ERROR: GetAvailableSlotsPerDay - Failed to find specified doctor")
		// return [][]time.Time{}
		return availableSlots
	}

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

	// if doctor does not any appointments schedules, they would have 2 available slot ranges
	// which is starting from work time till break time starts,
	// and from when the break time ends till the end of work time
	if len(_appointments) < 1 {
		// availableSlots := [][]time.Time{
		// 	{_doctor.WorkTime[0], _doctor.BreakTime[0]},
		// 	{_doctor.BreakTime[1], _doctor.WorkTime[1]},
		// }

		availableSlots = append(availableSlots,
			[]time.Time{_doctorWorkTime[0], _doctorBreakTime[0]},
			[]time.Time{_doctorBreakTime[1], _doctorWorkTime[1]},
		)

		return availableSlots
	}

	// if work starts before first appointment, then set first available slot range
	// to [work start time - appointment start time]
	firstAppointment := _appointments[0]
	if _doctorWorkTime[0].Before(firstAppointment.StartDatetime) {
		availableSlots = append(availableSlots,
			[]time.Time{_doctorWorkTime[0], firstAppointment.StartDatetime})

		// disregard first appointment if it starts at the same time as when the doctor starts working
		// as it was already taken care of above
		// _appointments = _appointments[1:]

		// if len(_appointments) < 1 {

		// 	if firstAppointment.StartDatetime.Before(_doctorBreakTime[0]) {
		// 		availableSlots = append(availableSlots,
		// 			[]time.Time{firstAppointment.EndDatetime, _doctorBreakTime[0]},
		// 			[]time.Time{_doctorBreakTime[1], _doctorWorkTime[1]})
		// 	} else {
		// 		availableSlots = append(availableSlots,
		// 			[]time.Time{firstAppointment.EndDatetime, _doctorWorkTime[1]})
		// 	}

		// 	return availableSlots
		// }
	}

	reachedAppointmentsAfterBreak := false
	for i, _appointment := range _appointments {
		// find only the first next appointment that starts after break time
		// if found, append available slots with ranges [end time of appointment - start time of break]
		// and [end time of break - start of next appointment]
		// if _appointment.StartDatetime.After(_doctorBreakTime[0]) && !reachedAppointmentsAfterBreak {
		// 	availableSlots = append(availableSlots,
		// 		[]time.Time{firstAppointment.EndDatetime, _doctorBreakTime[0]},
		// 		[]time.Time{_doctorBreakTime[1], _appointment.StartDatetime},
		// 	)
		// 	reachedAppointmentsAfterBreak = true
		// 	continue
		// }

		// if i = len(_appointments)-1{

		// }

		if _appointment.StartDatetime.After(_doctorBreakTime[0]) && !reachedAppointmentsAfterBreak {
			availableSlots = append(availableSlots,
				[]time.Time{_appointment.EndDatetime, _doctorBreakTime[0]},
				[]time.Time{_doctorBreakTime[1], _appointments[i+1].StartDatetime},
			)
			reachedAppointmentsAfterBreak = true
			continue
		}

		// append an available slot from appointment's end time till the start of the next one
		availableSlots = append(availableSlots,
			[]time.Time{_appointment.EndDatetime, _appointments[i+1].StartDatetime})

	}

	// if last appointment ends before doctor's work time ends
	// append an available slot with range [end of last appointment - end of doctor's work time]
	lastAppointment := _appointments[len(_appointments)-1]

	// if lastAppointment.EndDatetime.Before(_doctorWorkTime[1]) {
	// 	availableSlots = append(availableSlots,
	// 		[]time.Time{lastAppointment.EndDatetime, _doctorWorkTime[1]},
	// 	)
	// }

	if lastAppointment.EndDatetime.Before(_doctorBreakTime[0]) {
		availableSlots = append(availableSlots,
			[]time.Time{lastAppointment.EndDatetime, _doctorBreakTime[0]},
			[]time.Time{_doctorBreakTime[1], _doctorWorkTime[1]})
	} else {
		availableSlots = append(availableSlots,
			[]time.Time{lastAppointment.EndDatetime, _doctorWorkTime[1]})
	}

	return availableSlots

	// // 2 is added to accomodate both the work and break time ranges
	// availableSlots := make([][]time.Time, len(_appointments) + 2)
	// for i := range availableSlots {
	// 	availableSlots[i] = make([]time.Time, 2)
	// }

	// firstAppointment = _appointments[0]
	// // if start of first appointment equals start of work time, then set the first available slot to start
	// // from when the first appointment ends
	// if firstAppointment[0] == doctor.WorkTime[0] {
	// 	availableSlots = append(availableSlots, {
	// 		{firstAppointment[1]}
	// 	})
	// }

	// availableSlots = append(availableSlots, _appointments[])

}
