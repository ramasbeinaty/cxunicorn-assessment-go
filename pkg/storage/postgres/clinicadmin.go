package postgres

// defines the storage format of a clinic admin
type ClinicAdmin struct {
	Staff
	ID int `json:"id"`
}

type ClinicAdminCreate struct {
	ID int `json:"id"`
}
