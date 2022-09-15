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
			created_at timestamp WITH TIME ZONE NOT NULL DEFAULT (now()),
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
			work_days INTEGER[7] NOT NULL,
			work_time TIME WITH TIME ZONE[] NOT NULL,
			break_time TIME WITH TIME ZONE[] NOT NULL,
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

	createAppointmentsTable := `
	CREATE TABLE appointments (
		id serial PRIMARY KEY,
		patient_id INTEGER NOT NULL,
		doctor_id INTEGER NOT NULL,
		created_by INTEGER NOT NULL,
		created_at timestamp WITH TIME ZONE NOT NULL DEFAULT (now()),
		start_datetime timestamp WITH TIME ZONE NOT NULL,
		end_datetime timestamp WITH TIME ZONE NOT NULL,
		is_canceled BOOLEAN DEFAULT FALSE
	);`

	_, err = db.Exec(createAppointmentsTable)
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
		'$2a$10$jnbQAO21L2/Cfvm3Sy0qRugExJJBkOi2hstrZtMw78njs/YxW5kwO',
		'rama',
		'doe',
		'2021-10-10',
		'055-123456',
		'doctor'
	),
	(
		'raneem@gmail.com',
		'$2a$10$jnbQAO21L2/Cfvm3Sy0qRugExJJBkOi2hstrZtMw78njs/YxW5kwO',
        'raneem',
        'doe',
        '2020-09-20',
        '055-1234567',
        'doctor'
	),
	(
		'ronaldo@gmail.com',
		'$2a$10$jnbQAO21L2/Cfvm3Sy0qRugExJJBkOi2hstrZtMw78njs/YxW5kwO',
        'ronaldo',
        'doe',
        '2020-09-20',
        '055-1234567',
        'doctor'
	),
	(
		'patrik@gmail.com',
		'$2a$10$jnbQAO21L2/Cfvm3Sy0qRugExJJBkOi2hstrZtMw78njs/YxW5kwO',
        'patrik',
        'doe',
        '2020-09-20',
        '055-1234567',
        'patient'
	),
	(
		'prashant@gmail.com',
		'$2a$10$jnbQAO21L2/Cfvm3Sy0qRugExJJBkOi2hstrZtMw78njs/YxW5kwO',
        'prashant',
        'doe',
        '2020-09-20',
        '055-1234567',
        'patient'
	),
	(
		'amanda@gmail.com',
		'$2a$10$jnbQAO21L2/Cfvm3Sy0qRugExJJBkOi2hstrZtMw78njs/YxW5kwO',
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
		work_days,
		work_time,
		break_time
	) VALUES (
		1,
		'{0, 1, 2}',
		'{8:00:00Z, 18:00:00Z}',
		'{13:00:00Z, 14:00:00Z}'
	),
	(
		2,
		'{3, 4, 5}',
		'{8:00:00Z, 18:00:00Z}',
		'{13:00:00Z, 14:00:00Z}'
	),
	(
		3,
		'{1, 2, 3}',
		'{8:00:00Z, 18:00:00Z}',
		'{13:00:00Z, 14:00:00Z}'
	),
	(
		6,
		'{2, 3, 4}',
		'{8:00:00Z, 18:00:00Z}',
		'{13:00:00Z, 14:00:00Z}'
	)

	`
	// insertStaffs := `
	// INSERT INTO staffs (
	// 	id,
	// 	work_days,
	// 	work_time,
	// 	break_time,
	// 	unavailable_datetimes
	// ) VALUES (
	// 	1,
	// 	'{Mon, Tues, Wed}',
	// 	'18:00:00+04:00',
	// 	'{13:00:00+04:00, 14:00:00+04:00}',
	// 	'{}'
	// ),
	// (
	// 	2,
	// 	'{Sat, Tues, Wed}',
	// 	'18:00:00+04:00',
	// 	'{13:00:00+04:00, 14:00:00+04:00}',
	// 	'{}'
	// ),
	// (
	// 	3,
	// 	'{Fri, Sat, Sun}',
	// 	'18:00:00+04:00',
	// 	'{13:00:00+04:00, 14:00:00+04:00}',
	// 	'{}'
	// ),
	// (
	// 	6,
	// 	'{Fri, Sat, Sun}',
	// 	'18:00+04:00',
	// 	'{13:00:00+04:00, 14:00:00+04:00}',
	// 	'{}'
	// )

	// `

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
