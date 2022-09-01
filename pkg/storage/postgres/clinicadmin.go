package postgres

// defines the storage format of a clinic admin
type ClinicAdmin struct {
	Staff
	ID int32 `json:"id"`
}
