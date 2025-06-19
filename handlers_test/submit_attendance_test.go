package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"dealls-test/config"
	"dealls-test/handlers"
	"dealls-test/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSubmitAttendance_Success(t *testing.T) {
	// Inisialisasi DB test
	config.InitTestDB()
	config.DB = config.TestDB

	// Buat user dummy
	user := models.User{
		Username: "user_absen",
		Password: "secret",
		Role:     "employee",
	}
	config.TestDB.Create(&user)
	defer config.TestDB.Delete(&user)

	today := time.Now().Truncate(24 * time.Hour)

	// Pastikan tidak ada attendance dan audit sebelumnya
	config.TestDB.Where("user_id = ? AND attendance_date = ?", user.ID, today).Delete(&models.Attendance{})
	config.TestDB.Where("table_name = ? AND performed_by = ?", "attendances", user.ID).Delete(&models.AuditLog{})

	// Generate X-Request-ID
	requestID := uuid.New().String()

	// Siapkan request
	req := httptest.NewRequest("POST", "/attendances", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", user.ID))
	req.Header.Set("X-Request-ID", requestID)
	req.RemoteAddr = "127.0.0.1:1234"

	rr := httptest.NewRecorder()
	handlers.SubmitAttendance(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Attendance submitted successfully", response["message"])

	// Validasi attendance tersimpan
	var attendance models.Attendance
	err = config.TestDB.Where("user_id = ? AND attendance_date = ?", user.ID, today).First(&attendance).Error
	assert.NoError(t, err)
	assert.Equal(t, user.ID, attendance.UserID)

	// Validasi audit log tercatat dengan request_id
	var audit models.AuditLog
	err = config.TestDB.Where("table_name = ? AND record_id = ? AND request_id = ?", "attendances", attendance.ID, requestID).First(&audit).Error
	assert.NoError(t, err)
	assert.Equal(t, user.ID, audit.PerformedBy)

	// Clean-up
	config.TestDB.Delete(&attendance)
	config.TestDB.Delete(&audit)
}
