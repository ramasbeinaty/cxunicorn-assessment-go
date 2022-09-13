package auth

import "time"

type ClinicAdminRegister struct {
	UserDetails User `json:"user_details"`

	WorkDays  []int32       `json:"work_days" db:"work_days"`
	WorkTime  []time.Time `json:"work_time" db:"work_time"`
	BreakTime []time.Time `json:"break_time" db:"break_time"`
}
