package postgres

import (
	"clinicapp/pkg/storage/postgres/utils"
	"errors"
	"fmt"
	"log"
	"time"

	"database/sql"

	"github.com/blockloop/scan"
	_ "github.com/lib/pq"
)

type Storage struct {
	DB *sql.DB
}

// start and returns a new DB
func NewStorage() (*Storage, error) {
	var err error

	s := new(Storage)

	s.DB, err = sql.Open("postgres", "postgres://postgres:123456@localhost:5432/clinic_app2?sslmode=disable")

	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	// defer s.DB.Close()

	s.DB.SetConnMaxLifetime(0)
	s.DB.SetMaxOpenConns(3)

	return s, nil
}

func (s *Storage) GetDoctor(id int) (Doctor, error) {
	var doctor Doctor

	row := s.DB.QueryRow(`SELECT d.*, s.*, u.* FROM doctors AS d
				JOIN staffs AS s ON s.id = d.id
				JOIN users AS u ON u.id = d.id
				WHERE d.id = $1`, id)

	// err := scan.RowStrict(&doctor, row)

	if err := row.Scan(&doctor.ID, &doctor.Specialization, &doctor.ID,
		&doctor.WorkDays, &doctor.WorkTime, &doctor.BreakTime,
		&doctor.ID, &doctor.FirstName, &doctor.LastName, &doctor.DOB,
		&doctor.PhoneNumber, &doctor.Email, &doctor.Password, &doctor.Role,
		&doctor.CreatedAt, &doctor.IsActive, &doctor.IsVerified); err != nil {
		return doctor, errors.New(fmt.Sprintln("ERROR: GetDoctor - ", err))
	}

	return doctor, nil
}

func (s *Storage) GetAllDoctors() []Doctor {
	var doctors []Doctor = []Doctor{}

	// rows, _ := s.DB.Query(`SELECT * FROM doctors, staffs, users
	// 					   WHERE staffs.id = doctors.id AND users.id = doctors.id
	// 						`)

	rows, err := s.DB.Query(`SELECT * FROM doctors
						 	JOIN staffs ON doctors.id = staffs.id
							JOIN users ON doctors.id = users.id
							`)
	if err != nil {
		fmt.Println("GetAllDoctors - Was not able to execute query", err.Error())
		return doctors
	}

	defer rows.Close()

	for rows.Next() {
		var doctor Doctor

		if err = rows.Scan(&doctor.ID, &doctor.Specialization, &doctor.ID,
			&doctor.WorkDays, &doctor.WorkTime, &doctor.BreakTime,
			&doctor.ID, &doctor.FirstName, &doctor.LastName, &doctor.DOB,
			&doctor.PhoneNumber, &doctor.Email, &doctor.Password, &doctor.Role,
			&doctor.CreatedAt, &doctor.IsActive, &doctor.IsVerified); err != nil {
			fmt.Println("ERROR: GetAllDoctors - ", err)
			return doctors
		}

		doctors = append(doctors, doctor)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("GetAllDoctors - Was not able to scan rows", err.Error())
		return doctors
	}

	// if err := scan.RowsStrict(&doctors, rows); err != nil {
	// 	if err == sql.ErrNoRows {
	// 		fmt.Println("Warning: GetAllDoctors - No doctors exist")
	// 		return doctors
	// 	}

	return doctors

}

func (s *Storage) DoctorExists(id int) bool {

	var doctorID int

	row, _ := s.DB.Query(`SELECT id FROM doctors
							WHERE doctors.id = $1`, id)

	if err := scan.RowStrict(&doctorID, row); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("INFO: doctorExists - doctor id does not exist")
			return false
		}

		fmt.Println("ERROR: doctorExists -", err)
		return false
	}

	fmt.Println("INFO: doctorExists - doctor id exists")
	return true
}

func (s *Storage) CreateAppointment(a AppointmentCreate) error {

	rows, err := s.DB.Query(`
		INSERT INTO appointments (patient_id, doctor_id, created_by, start_datetime,
		end_datetime)
		VALUES ($1, $2, $3, $4, $5)`,
		a.PatientID, a.DoctorID, a.CreatedBy, a.StartDatetime, a.EndDatetime,
	)

	if err != nil {
		return errors.New("ERROR: Create Appointment - " + err.Error())
	}

	defer rows.Close()

	return nil
}

func (s *Storage) GetAppointments() []Appointment {
	var appointments []Appointment

	return appointments
}

func (s *Storage) GetAppointment(id int) Appointment {
	var appointment Appointment

	row := s.DB.QueryRow(`
		SELECT *
		FROM appointments 
		WHERE id = $1
		`, id)

	if err := row.Scan(&appointment); err != nil {
		fmt.Println("GetAppointment - ", err)
		return appointment
	}

	return appointment
}

func (s *Storage) EditAppointment(id int, appointment AppointmentEdit) error {
	rows, err := s.DB.Query(`
		UPDATE appointments 
		SET start_datetime = $1, end_datetime = $2, is_canceled = $3
		WHERE id = $4`, appointment.StartDatetime, appointment.EndDatetime, appointment.IsCanceled, id)

	if err != nil {
		return errors.New("EditAppointment - Was not able to execute query - " + err.Error())
	}

	defer rows.Close()

	return nil
}

// get the number of appointments a doctor has with unique patients in a given date
func (s *Storage) GetNumberOfAppointmentsWithDistinctPatient(doctorID int, date time.Time) int {
	var appointments_count int = 0

	row := s.DB.QueryRow(`
		SELECT COUNT(DISTINCT(patient_id, $1))
		FROM appointments a
		WHERE CAST(a.start_datetime as DATE) = CAST($2 as DATE)
		`, doctorID, date)

	if err := row.Scan(&appointments_count); err != nil {
		fmt.Println("GetNumberOfAppointmentsWithDistinctPatient - ", err)
		return appointments_count
	}

	return appointments_count

}

func (s *Storage) GetAllAppointmentsOfDoctor(doctorID int, date time.Time) []Appointment {
	var appointments []Appointment = []Appointment{}

	rows, err := s.DB.Query(`
		SELECT * 
		FROM appointments a
		WHERE a.doctors_id = $1	AND CAST(a.start_datetime as DATE) = CAST($2 as DATE)`, doctorID, date)

	if err != nil {
		fmt.Println("GetAllAppointment - Was not able to execute query", err.Error())
		return appointments
	}

	defer rows.Close()

	for rows.Next() {
		var appointment Appointment

		if err = rows.Scan(&appointment.ID, &appointment.PatientID, &appointment.DoctorID,
			&appointment.CreatedAt, &appointment.CreatedBy, &appointment.StartDatetime,
			&appointment.EndDatetime, &appointment.IsCanceled); err != nil {
			fmt.Println("ERROR: GetAllAppointment - ", err)
			return appointments
		}

		appointments = append(appointments, appointment)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("GetAllAppointment - Was not able to scan rows", err.Error())
		return appointments
	}

	return appointments

}

func (s *Storage) GetAppointmentHoursPerDay(doctorID int, date time.Time) int {
	var hours int = 0

	row := s.DB.QueryRow(`
		SELECT extract(hour FROM total_appointments_duration) FROM 
		(SELECT SUM(duration) AS total_appointments_duration FROM 
			(SELECT end_datetime - start_datetime AS duration
			FROM appointments 
			WHERE doctor_id = $1 AND CAST(start_datetime AS DATE) = $2
			) appointment_duration
		) total_duration_hours;
		`, doctorID, date)

	if err := row.Scan(&hours); err != nil {
		fmt.Println("GetNumberOfAppointmentsWithDistinctPatient - ", err)
		return hours
	}

	return hours

	// 	--with appointment(start_datetime, end_datetime)
	// --as (select start_datetime, end_datetime
	// --	from appointments
	// --	where doctor_id = 2 and cast(start_datetime as date) = '2022-08-09'),
	// --	hours_sum as (select sum(a.end_datetime-a.start_datetime) from appointment a)
	// --select hours_sum from hours_sum;

	// --with appointment(start_datetime, end_datetime)
	// --as (select start_datetime, end_datetime
	// --	from appointments
	// --	where doctor_id = 2 and cast(start_datetime as date) = '2022-08-09'),
	// --	hours as (select extract('epoch' from a.end_datetime-a.start_datetime)/3600.00 from appointment a)
	// --select hours from hours;

}

func (s *Storage) IsAppointmentWithinDoctorWorkDays(doctorID int, day time.Weekday) bool {
	var _doctorID int

	row := s.DB.QueryRow(`
		SELECT id
		FROM staffs
		WHERE id=$1 and $2 = any (work_days);
		`, doctorID, day)

	if err := row.Scan(&_doctorID); err != nil {
		fmt.Println("IsAppointmentWithinDoctorWorkDays - ", err)
		return false
	}

	return true
}

func (s *Storage) IsAppointmentOverlapping(doctorID int, patientID int, startDateTime time.Time, endDatetime time.Time) bool {
	var isOverlapping bool = false

	row := s.DB.QueryRow(`
		SELECT CAST(CASE WHEN COUNT(*) > 0 THEN 1 ELSE 0 END AS BIT)
		FROM (SELECT start_datetime, end_datetime
		FROM appointments
		WHERE doctor_id=$1 or patient_id=$2) app 
		WHERE (app.start_datetime, app.end_datetime) OVERLAPS ($3, $4);
	`, doctorID, patientID, startDateTime, endDatetime)

	if err := row.Scan(&isOverlapping); err != nil {
		fmt.Println("IsAppointmentOverlapping - ", err)
		return false
	}

	return isOverlapping
}

func (s *Storage) GetDoctorWorkTime(doctorID int) []time.Time {
	var _workTime utils.TimeArray = []time.Time{}

	row := s.DB.QueryRow(`
		SELECT work_time
		FROM staffs
		WHERE id=$1;
		`, doctorID)

	if err := row.Scan(&_workTime); err != nil {
		fmt.Println("IsAppointmentWithinDoctorWorkDays - ", err)
		return _workTime
	}

	return _workTime
}

func (s *Storage) GetDoctorBreakTime(doctorID int) []time.Time {
	var _breakTime utils.TimeArray = []time.Time{}

	row := s.DB.QueryRow(`
		SELECT break_time
		FROM staffs
		WHERE id=$1;
		`, doctorID)

	if err := row.Scan(&_breakTime); err != nil {
		fmt.Println("IsAppointmentWithinDoctorWorkDays - ", err)
		return _breakTime
	}

	return _breakTime
}

func (s *Storage) CancelAppointment(id int) error {
	var canceled bool = true

	rows, err := s.DB.Query(`
		UPDATE appointments 
		SET is_canceled = $1
		WHERE id = $2`, canceled, id)

	if err != nil {
		return errors.New("CancelAppointment - Was not able to execute query - " + err.Error())
	}

	defer rows.Close()

	return nil
}

func (s *Storage) DeleteAppointment(id int) error {
	rows, err := s.DB.Query(`
		DELETE FROM appointments
		WHERE id = $1`, id)

	if err != nil {
		return errors.New("DeleteAppointment - Was not able to execute query - " + err.Error())
	}

	defer rows.Close()

	return nil
}
