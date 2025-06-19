package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"dealls-test/config"
	"dealls-test/models"
	"dealls-test/utils"

	"github.com/google/uuid"
)

func SubmitOvertime(w http.ResponseWriter, r *http.Request) {
	// Ambil IP address dari fungsi utils
	ip := utils.GetIP(r)

	// Ambil atau generate request_id
	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = uuid.New().String()
	}

	// Parse token dari header
	tokenStr := r.Header.Get("Authorization")
	claims, err := utils.ParseToken(strings.TrimPrefix(tokenStr, "Bearer "))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input struct {
		Hours int `json:"hours"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Hours <= 0 || input.Hours > 3 {
		http.Error(w, "Invalid hours (1-3 allowed)", http.StatusBadRequest)
		return
	}

	today := time.Now().Truncate(24 * time.Hour)

	// Cek apakah sudah ada submission lembur hari ini
	var existing models.Overtime
	if err := config.DB.Where("user_id = ? AND overtime_date = ?", claims.UserID, today).First(&existing).Error; err == nil {
		http.Error(w, "Overtime already submitted today", http.StatusConflict)
		return
	}

	overtime := models.Overtime{
		UserID:       claims.UserID,
		OvertimeDate: today,
		Hours:        input.Hours,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		CreatedBy:    claims.UserID,
		UpdatedBy:    claims.UserID,
		IPAddress:    ip,
		RequestID:    requestID, // ✅ Tambahkan request_id di sini
	}

	if err := config.DB.Create(&overtime).Error; err != nil {
		http.Error(w, "Failed to submit overtime", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Overtime submitted successfully",
		"data":    overtime,
	})
}
