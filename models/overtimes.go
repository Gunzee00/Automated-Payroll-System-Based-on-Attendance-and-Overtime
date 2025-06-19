package models

import "time"

type Overtime struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `json:"user_id"`
	OvertimeDate time.Time `json:"overtime_date"`
	Hours        int       `json:"hours"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedBy    uint      `json:"created_by"`
	UpdatedBy    uint      `json:"updated_by"`
	IPAddress    string    `json:"ip_address"`
	RequestID    string    `gorm:"type:varchar(100)" json:"request_id"`
}
