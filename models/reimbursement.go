package models

import "time"

type Reimbursement struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   uint      `json:"created_by"`
	UpdatedBy   uint      `json:"updated_by"`
	IPAddress   string    `json:"ip_address"`
	RequestID   string    `gorm:"type:varchar(100)" json:"request_id"`
}
