package auth

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

const (
	DoctorRole      string = "doctor"
	PatientRole          = "patient"
	ClinicAdminRole      = "clinic_admin"
)

// type RoleTypes struct {
// 	RoleName Role `json:"Role"`
// }

// type EventName string

// const (
// 	NEW_USER       EventName = "NEW_USER"
// 	DIRECT_MESSAGE EventName = "DIRECT_MESSAGE"
// 	DISCONNECT     EventName = "DISCONNECT"
// )

// type ConnectionPayload struct {
// 	EventName    EventName   `json:"eventName" validate:"oneof=NEW_USER DIRECT_MESSAGE DISCONNECT"`
// 	EventPayload interface{} `json:"eventPayload"`
// }

// func (s *ConnectionPayload) Validate() error {
// 	validate := validator.New()
// 	return validate.Struct(s)
// }
