package editing

import "time"

// defines the appointment fields that can be edited
type Appointment struct {
	StartDatetime time.Time `json:"start_datetime"`
	EndDatetime   time.Time `json:"end_datetime"`
	IsCanceled    bool      `json:"is_canceled"`
}