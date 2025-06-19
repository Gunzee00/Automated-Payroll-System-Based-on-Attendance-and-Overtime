package models
import "time"
type Payroll struct {
	ID                 uint      `json:"id" gorm:"primaryKey"`
	UserID             uint      `json:"user_id"`
	AttendancePeriodID uint      `json:"attendance_period_id"`
	BaseSalary         float64   `json:"base_salary"`
	TotalOvertime      float64   `json:"total_overtime"`
	TotalReimbursement float64   `json:"total_reimbursement"`
	TotalSalary        float64   `json:"total_salary"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedBy          uint      `json:"created_by"`
	IPAddress          string    `json:"ip_address"`
	RequestID          string    `gorm:"type:varchar(100)" json:"request_id"`
}
