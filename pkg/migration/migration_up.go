package migration

import (
	"clinicapp/pkg/storage/postgres"
	"database/sql"
	"log"
)

func MigrateUp() {
	s, _ := postgres.NewStorage()
	createTables(s.DB)
	populateTables(s.DB)
}

func createTables(db *sql.DB) {
	var err error

	createRoleType := `
	CREATE TYPE role_type AS ENUM ('doctor', 'patient', 'clinic_admin');
	`

	_, err = db.Exec(createRoleType)
	if err != nil {
		log.Fatal("Failed to create type role - ", err)
	}

	createUsersTable := `
		CREATE TABLE users (
			id serial PRIMARY KEY,
			first_name varchar NOT NULL,
			last_name varchar NOT NULL,
			dob timestamp NOT NULL,
			phone_number varchar NOT NULL,
			email varchar NOT NULL UNIQUE,
			password varchar NOT NULL,
			role role_type NOT NULL,
			created_at timestamp NOT NULL DEFAULT (now()),
			is_active boolean DEFAULT true,
			is_verified boolean DEFAULT false
		);

		CREATE INDEX idx_users_email ON users (email);`

	_, err = db.Exec(createUsersTable)
	if err != nil {
		log.Fatal("Failed to create users table - ", err)
	}

	createPatientsTable := `
		CREATE TABLE patients (
			id serial NOT NULL UNIQUE,
			medical_history varchar NOT NULL,
			CONSTRAINT FK_id FOREIGN KEY (id) REFERENCES users(id)
		);
	`

	_, err = db.Exec(createPatientsTable)
	if err != nil {
		log.Fatal("Failed to create patients table - ", err)
	}

	createStaffsTable := `
		CREATE TABLE staffs (
			id serial NOT NULL UNIQUE,
			work_shift varchar NOT NULL,
			CONSTRAINT FK_id FOREIGN KEY (id) REFERENCES users(id)
		);
	`

	_, err = db.Exec(createStaffsTable)
	if err != nil {
		log.Fatal("Failed to create staffs table - ", err)
	}

	createDoctorsTable := `
		CREATE TABLE doctors (
			id serial NOT NULL UNIQUE,
			specialization varchar NOT NULL,
			CONSTRAINT FK_id FOREIGN KEY (id) REFERENCES staffs(id)
		);
	`

	_, err = db.Exec(createDoctorsTable)
	if err != nil {
		log.Fatal("Failed to create doctors table - ", err)
	}

	createClinicAdminsTable := `
		CREATE TABLE clinic_admins (
			id serial NOT NULL UNIQUE,
			CONSTRAINT FK_id FOREIGN KEY (id) REFERENCES staffs(id)
		);
	`

	_, err = db.Exec(createClinicAdminsTable)
	if err != nil {
		log.Fatal("Failed to create clinic_admins table - ", err)
	}

}

func populateTables(db *sql.DB) {
	var err error

	insertUsers := `Insert into users(
		email,
		password,
		first_name,
    	last_name,
		dob,
		phone_number,
		role
	) values (
		'rama@gmail.com',
		'123456',
		'rama',
		'doe',
		'2021-10-10',
		'055-123456',
		'doctor'
	),
	(
		'raneem@gmail.com',
		'123456',
        'raneem',
        'doe',
        '2020-09-20',
        '055-1234567',
        'doctor'
	),
	(
		'ronaldo@gmail.com',
		'123456',
        'ronaldo',
        'doe',
        '2020-09-20',
        '055-1234567',
        'doctor'
	),
	(
		'patrik@gmail.com',
		'123456',
        'patrik',
        'doe',
        '2020-09-20',
        '055-1234567',
        'patient'
	),
	(
		'prashant@gmail.com',
		'123456',
        'prashant',
        'doe',
        '2020-09-20',
        '055-1234567',
        'patient'
	),
	(
		'amanda@gmail.com',
		'123456',
        'amanda',
        'doe',
        '2020-09-20',
        '055-1234567',
        'clinic_admin'
	);
	`

	_, err = db.Exec(insertUsers)
	if err != nil {
		log.Fatal("Failed to populate users table - ", err)
	}

	// remember that the ids of patients are hard-coded
	insertPatients := `
	INSERT INTO patients(
		id,
		medical_history
	) VALUES (
		4,
		'asthma'
	),
	(
		5,
		'had a hernia surgery'
	)`

	// insertPatients := `
	// 	INSERT INTO patients(
	// 		id,
	// 		medical_history
	// 	) SELECT (SELECT id from users WHERE role='patient' AND id='4'), 'had brain surgery'
	// `

	_, err = db.Exec(insertPatients)
	if err != nil {
		log.Fatal("Failed to populate patients table - ", err)
	}

	insertStaffs := `
	INSERT INTO staffs (
		id,
		work_shift
	) VALUES (
		'1',
		'morning_shift'
	),
	(
		'2',
		'night_shift'
	),
	(
		'3',
		'morning_shift'
	),
	(
		'6',
		'morning_shift'
	)

	`

	_, err = db.Exec(insertStaffs)
	if err != nil {
		log.Fatal("Failed to populate staffs table - ", err)
	}

	insertDoctors := `
		INSERT INTO doctors (
			id,
			specialization
		) 	VALUES (
			1,
			'Cardiology'
		),
		(
			2,
			'Pediatrician'
		),
		(
			3,
			'Neurology'
		)
	`

	_, err = db.Exec(insertDoctors)
	if err != nil {
		log.Fatal("Failed to populate doctors table - ", err)
	}

	insertIntoClinicAdmins := `
		INSERT INTO clinic_admins (
			id
		) VALUES (
			6
		)
	`

	_, err = db.Exec(insertIntoClinicAdmins)
	if err != nil {
		log.Fatal("Failed to populate clinic_admins table - ", err)
	}

}
