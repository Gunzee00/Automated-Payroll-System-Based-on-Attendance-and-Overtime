package handlers

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"time"

	"dealls-test/config"
	"dealls-test/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SubmitAttendance(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Ambil atau generate request_id
	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = uuid.New().String()
	}

	today := time.Now().Truncate(24 * time.Hour)
	weekday := today.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		http.Error(w, "Attendance is not allowed on weekends", http.StatusBadRequest)
		return
	}

	var existing models.Attendance
	err := config.DB.Where("user_id = ? AND attendance_date = ?", userID, today).First(&existing).Error
	if err == nil {
		http.Error(w, "You have already submitted attendance today", http.StatusConflict)
		return
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		http.Error(w, "Failed to check attendance", http.StatusInternalServerError)
		return
	}

	ipAddress := getIPAddress(r)

	attendance := models.Attendance{
		UserID:         userID.(uint),
		AttendanceDate: today,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		CreatedBy:      userID.(uint),
		UpdatedBy:      userID.(uint),
		IPAddress:      ipAddress,
		RequestID:      requestID, // ✅ Tambahkan request_id ke model
	}

	if err := config.DB.Create(&attendance).Error; err != nil {
		http.Error(w, "Failed to submit attendance", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Attendance submitted successfully",
		"request_id": requestID, // ✅ Tambahkan ke response
		"data":       attendance,
	})
}

func getIPAddress(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	} else {
		ip = strings.Split(ip, ",")[0]
	}
	return ip
}
