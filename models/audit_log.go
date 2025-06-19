package models

import (
    "time"
)

 
type AuditLog struct {
     ID          uint      `gorm:"primaryKey" json:"id"`
     TableName   string    `gorm:"not null" json:"table_name"`
     RecordID    uint      `gorm:"not null" json:"record_id"`
     Action      string    `gorm:"not null" json:"action"`
     ChangedData string    `gorm:"type:text" json:"changed_data"`
     PerformedBy uint      `gorm:"not null" json:"performed_by"`
     IPAddress   string    `gorm:"type:varchar(45)" json:"ip_address"` 
     RequestID   string    `gorm:"type:varchar(100)" json:"request_id"`
     CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}
