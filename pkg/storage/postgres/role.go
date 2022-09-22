package postgres

const (
	DoctorRole      string = "doctor"
	PatientRole            = "patient"
	ClinicAdminRole        = "clinic_admin"
)

var Roles = newRoleRegistry()

func newRoleRegistry() *roleRegistry {
	return &roleRegistry{
		Patient:     "patient",
		Doctor:      "doctor",
		ClinicAdmin: "ClinicAdmin",
	}
}

type roleRegistry struct {
	Patient     string
	Doctor      string
	ClinicAdmin string
}
