package handlers

import (
	"dealls-test/config"
	"dealls-test/models"
	"dealls-test/utils"
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateAttendancePeriod(w http.ResponseWriter, r *http.Request) {
	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = uuid.New().String()
	}

	tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	claims, err := utils.ParseToken(tokenString)
	if err != nil || claims.Role != "admin" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ipAddress := getIP(r)

	var input struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		http.Error(w, "Invalid start date format", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		http.Error(w, "Invalid end date format", http.StatusBadRequest)
		return
	}

	if endDate.Before(startDate) {
		http.Error(w, "End date cannot be before start date", http.StatusBadRequest)
		return
	}

	now := time.Now()

	attendancePeriod := models.AttendancePeriod{
		StartDate: startDate,
		EndDate:   endDate,
		Status:    "open",
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: claims.UserID,
		UpdatedBy: claims.UserID,
		IPAddress: ipAddress,
		RequestID: requestID, // ✅ Ditambahkan
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&attendancePeriod).Error; err != nil {
			return err
		}

		audit := models.AuditLog{
			TableName:   "attendance_periods",
			RecordID:    attendancePeriod.ID,
			Action:      "CREATE",
			ChangedData: toJSON(attendancePeriod),
			PerformedBy: claims.UserID,
			IPAddress:   ipAddress,
			RequestID:   requestID,
			CreatedAt:   now,
		}

		return tx.Create(&audit).Error
	})

	if err != nil {
		http.Error(w, "Failed to create attendance period", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Attendance period created",
		"data":    attendancePeriod,
	})
}


func getIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return strings.Split(ip, ",")[0]
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // fallback
	}

	return host
}

func toJSON(data interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		return "{}"
	}
	return string(b)
}
