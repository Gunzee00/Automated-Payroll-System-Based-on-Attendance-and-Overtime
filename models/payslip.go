package models

import "time"

type Payslip struct {
	ID                 uint      `json:"id" gorm:"primaryKey"`
	UserID             uint      `json:"user_id"`
	AttendancePeriodID uint      `json:"attendance_period_id"`

	AttendanceCount    int       `json:"attendance_count"`
	BaseSalary         float64   `json:"base_salary"`
	TotalOvertime      float64   `json:"total_overtime"`
	TotalReimbursement float64   `json:"total_reimbursement"`
	TotalSalary        float64   `json:"total_salary"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	CreatedBy  uint   `json:"created_by"`
	UpdatedBy  uint   `json:"updated_by"`
	IPAddress  string `json:"ip_address"`
	RequestID  string `json:"request_id" gorm:"type:varchar(100)"`
}
