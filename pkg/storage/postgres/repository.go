package postgres

import (
	"errors"
	"fmt"
	"log"

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

// func (s *Storage) GetAllDoctors() []listing.Doctor {
// 	var doctors []listing.Doctor = []listing.Doctor{}

// 	// rows, _ := s.DB.Query(`SELECT * FROM doctors, staffs, users
// 	// 					   WHERE staffs.id = doctors.id AND users.id = doctors.id
// 	// 						`)

// 	rows, _ := s.DB.Query(`SELECT * FROM doctors
// 						 	JOIN staffs ON doctors.id = staffs.id
// 							JOIN users ON doctors.id = users.id
// 							`)

// 	if err := scan.RowsStrict(&doctors, rows); err != nil {
// 		if err == sql.ErrNoRows {
// 			fmt.Println("Warning: GetAllDoctors", listing.ErrEmpty)
// 			return doctors
// 		}

// 		fmt.Println("ERROR: GetAllDoctors - ", err)
// 		return doctors
// 	}

// 	return doctors

// }

// func (s *Storage) CreateAppointment(a booking.Appointment) error {

// 	row, err := s.DB.Query(`
// 		INSERT INTO appointments (patient_id, doctor_id, created_by, start_datetime,
// 		end_datetime)
// 		VALUES ($1, $2, $3, $4, $5)`,
// 		a.PatientID, a.DoctorID, a.CreatedBy, a.StartDatetime, a.EndDatetime,
// 	)

// 	row.Close()

// 	if err != nil {
// 		fmt.Sprintln("ERROR: Create Appointment - ", err)
// 		return errors.New(booking.Err.Error())
// 	}

// 	return nil
// }

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

func (s *Storage) GetDoctor(id int) (Doctor, error) {
	var doctor Doctor

	// var str string

	// rows := s.DB.QueryRow(`SELECT first_name FROM users where users.id = $1`, id)
	// rows.Scan(&str)

	row := s.DB.QueryRow(`SELECT d.*, s.*, u.* FROM doctors AS d
				JOIN staffs AS s ON s.id = d.id
				JOIN users AS u ON u.id = d.id
				WHERE d.id = $1`, id)

	// err := scan.RowStrict(&doctor, row)

	// var workTime PgTime

	if err := row.Scan(&doctor.ID, &doctor.Specialization, &doctor.ID,
		&doctor.WorkDays, &doctor.WorkTime, &doctor.BreakTime,
		&doctor.ID, &doctor.FirstName, &doctor.LastName, &doctor.DOB,
		&doctor.PhoneNumber, &doctor.Email, &doctor.Password, &doctor.Role,
		&doctor.CreatedAt, &doctor.IsActive, &doctor.IsVerified); err != nil {
		return doctor, errors.New(fmt.Sprintln("ERROR: GetDoctor - ", err))
	}

	return doctor, nil
}
