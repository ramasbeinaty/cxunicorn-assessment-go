package postgres

import (
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

func (s *Storage) GetAppointments() []Appointment {
	var appointments []Appointment

	return appointments
}

func (s *Storage) CreateAppointment(a AppointmentCreate) error {

	row, err := s.DB.Query(`
		INSERT INTO appointments (patient_id, doctor_id, created_by, start_datetime,
		end_datetime)
		VALUES ($1, $2, $3, $4, $5)`,
		a.PatientID, a.DoctorID, a.CreatedBy, a.StartDatetime, a.EndDatetime,
	)

	if err != nil {
		return errors.New("ERROR: GetAllDoctors - Was not able to execute query" + err.Error())
	}

	defer row.Close()

	if err != nil {
		return errors.New("ERROR: Create Appointment - " + err.Error())
	}

	return nil
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

// get the number of appointments a doctor has with unique patients in a given date
func (s *Storage) GetNumberOfAppointmentsWithDistinctPatient(doctor_id int, date time.Time) int {
	var appointments_count int = 0

	row := s.DB.QueryRow(`
	SELECT COUNT(DISTINCT(patient_id, doctor_id))
	FROM appointments a
	WHERE CAST(a.start_datetime as DATE) = CAST($1 as DATE)
	`, date)

	if err := row.Scan(&appointments_count); err != nil {
		fmt.Println("GetNumberOfAppointmentsWithDistinctPatient - ", err)
		return appointments_count
	}

	return appointments_count

}

// // get doctor with given id
// func (s *Storage) GetDoctor(id string) (listing.Doctor, error) {
// 	// if row := s.DB.QueryRow(`
// 	// SELECT first_name, last_name, email, work_shift, specialization, doctors.id from users, staffs, doctors
// 	// WHERE users.id = $1 AND staffs.id = $1 AND doctors.id = $1`, id); row != nil {
// 	// 	if err := row.Scan(&doctor.FirstName, &doctor.LastName, &doctor.Email,
// 	// 		&doctor.WorkShift, &doctor.Specialization, &doctor.ID); err != nil {
// 	// 		return doctor, errors.New(fmt.Sprintln("ERROR: GetDoctor - ", err))
// 	// 	}
// 	// 	if row.Err() == sql.ErrNoRows {
// 	// 		return doctor, errors.New(fmt.Sprintln("ERROR: GetDoctor - ", listing.ErrIdNotFound))
// 	// 	}
// 	// 	return doctor, errors.New(fmt.Sprintln("ERROR: GetDoctor - ", row))
// 	// }

// 	// row := s.DB.QueryRow(`
// 	// SELECT * FROM users, staffs, doctors
// 	// WHERE users.id = $1 AND staffs.id = $1 AND doctors.id = $1`, id)

// 	// scanning := structscanner.Select(s.DB, &d, "",
// 	// 	`SELECT * FROM users, staffs, doctors
// 	// 	WHERE users.id = 1 AND staffs.id = 1 AND doctors.id = 1`)

// 	var doctor listing.Doctor

// 	row, _ := s.DB.Query(`SELECT * FROM users, staffs, doctors
// 							WHERE users.id = $1 AND
// 							staffs.id = $1 AND doctors.id = $1`, id)

// 	if err := scan.RowStrict(&doctor, row); err != nil {
// 		if err == sql.ErrNoRows {
// 			return doctor, errors.New(fmt.Sprintln("ERROR: GetDoctor - ", listing.ErrIdNotFound))
// 		}
// 		return doctor, errors.New(fmt.Sprintln("ERROR: GetDoctor - ", err))

// 	}

// 	return doctor, nil
// }
