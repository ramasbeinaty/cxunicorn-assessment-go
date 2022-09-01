package postgres

import "time"

// defines the storage format of an event
type Event struct {
	ID                 int32     `json:"id"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedByUser      int32     `json:"created_by_user"`
	EventStartDatetime time.Time `json:"event_start_datetime"`
	EventEndDatetime   time.Time `json:"event_end_datetime"`
	IsCanceled         bool      `json:"is_canceled"`
}
