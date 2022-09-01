package garbage

import (
	"log"
	"os"

	"database/sql"

	_ "github.com/lib/pq"
)

var DB *sql.DB

var (
	dbScheme = os.Getenv("DB_SCHEME")
	dbHost   = os.Getenv("DB_HOST")
	// dbPort, err2 = strconv.Atoi(os.Getenv("DB_PORT"))
	dbPort     = os.Getenv("DB_PORT")
	dbUsername = os.Getenv("DB_USERNAME")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbName     = os.Getenv("DB_NAME")
	dbSSLMode  = os.Getenv("DB_SSLMODE")
)

func NewDB() *sql.DB {
	/*
	start and return a db
	*/
	var err error

	// dsn := url.URL{
	// 	Scheme: os.Getenv("DB_SCHEME"),
	// 	Host:   os.Getenv("DB_HOST"),
	// 	User:   url.UserPassword(os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD")),
	// 	Path:   os.Getenv("DB_Name"),
	// }

	// q := dsn.Query()
	// q.Add("sslmode", os.Getenv("DB_SSLMODE"))

	// dsn.RawQuery = q.Encode()

	// DB, err = sql.Open("pgx", dsn.String())

	// DB, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
	// 	dbHost, dbPort, dbUsername, dbPassword, dbName, dbSSLMode))

	// if err != nil {
	// 	// This will not be a connection error, but a DSN parse error or
	// 	// another initialization error.
	// 	log.Fatal("unable to use data source name (dsn)", err)
	// }

	DB, err = sql.Open("postgres", "postgres://postgres:123456@localhost:5432/clinic_app2?sslmode=disable")

	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	// defer DB.Close()

	DB.SetConnMaxLifetime(0)
	DB.SetMaxOpenConns(3)

	return DB

	// row := DB.QueryRowContext(context.Background(), "SELECT role FROM users WHERE first_name = 'rama'")
	// if err := row.Err(); err != nil {
	// 	log.Fatal("Failed to get user role from db")
	// 	return
	// }

	// var role string

	// if err := row.Scan(&role); err != nil {
	// 	log.Fatal("Row.Scan error", err)
	// 	return
	// }

	// fmt.Println("role: ", role)

}
