package cache

import (
	"time"
)

// defines the storage format of a staff
type Staff struct {
	User
	ID        int         `json:"id"`
	WorkDays  []int32     `json:"work_days"`
	WorkTime  []time.Time `json:"work_time"`
	BreakTime []time.Time `json:"break_time"`
}
