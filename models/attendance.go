package models

import "time"

type Attendance struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `json:"user_id"`
	AttendanceDate time.Time `gorm:"column:attendance_date" json:"attendance_date"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	CreatedBy      uint      `json:"created_by"`
	UpdatedBy      uint      `json:"updated_by"`
	IPAddress      string    `gorm:"type:inet" json:"ip_address"`
	    RequestID string `gorm:"type:varchar(100)" json:"request_id"`

}
