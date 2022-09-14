package migration

import (
	"clinicapp/pkg/storage/postgres"
	"database/sql"
	"log"
)

func MigrateDown() {
	s, _ := postgres.NewStorage()
	deleteTables(s.DB)
}

func deleteTables(db *sql.DB) {
	// dropTables := `
	// 	DROP TABLE patients;
	// 	DROP TABLE users;
	// 	DROP TYPE role_type;
	// 	`
	dropTables := `
		DROP TABLE appointments;
		DROP TABLE doctors;
		DROP TABLE clinic_admins;
		DROP TABLE staffs;
		DROP TABLE patients;
		DROP TABLE users;
		DROP TYPE role_type;
		`
	// DROP INDEX idx_users_email;

	_, err := db.Exec(dropTables)
	if err != nil {
		log.Fatal("Failed to drop tables - ", err)
	}

}
