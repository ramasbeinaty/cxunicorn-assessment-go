package auth

var Roles = newRoleRegistry()

func newRoleRegistry() *roleRegistry {
	return &roleRegistry{
		Patient:     "patient",
		Doctor:      "doctor",
		ClinicAdmin: "clinic_admin",
	}
}

type roleRegistry struct {
	Patient     string
	Doctor      string
	ClinicAdmin string
}
