package postgres

import (
	"github.com/lib/pq"
)

// defines the storage format of a staff
type Staff struct {
	User
	ID        int            `json:"id" db:"id"`
	WorkDays  pq.StringArray `json:"work_days" db:"work_days"`
	WorkTime  TimeArray      `json:"work_time" db:"work_time"`
	BreakTime TimeArray      `json:"break_time" db:"break_time"`
}

// type Staff struct {
// 	User
// 	ID                   int           `json:"id" db:"id"`
// 	WorkDays             []string      `json:"work_days" db:"work_days"`
// 	WorkTime             []time.Time   `json:"work_time" db:"work_time"`
// 	BreakTime            []time.Time   `json:"break_time" db:"break_time"`
// 	UnavailableDatetimes [][]time.Time `json:"unavailable_datetimes" db:"unavailable_datetimes"`
// }
