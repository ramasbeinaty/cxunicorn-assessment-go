package canceling

// defines the appointment fields that can be edited
type Appointment struct {
	IsCanceled    bool      `json:"is_canceled"`
}