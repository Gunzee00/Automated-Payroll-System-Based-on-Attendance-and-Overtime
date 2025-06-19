package models

import "time"

type AttendancePeriod struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy uint      `json:"created_by"`
	UpdatedBy uint      `json:"updated_by"`
	IPAddress string    `json:"ip_address"`
	RequestID string    `gorm:"type:varchar(100)" json:"request_id"`
}
