package handlers_test

import (
	"bytes"
	"dealls-test/config"
	"dealls-test/handlers"
	"dealls-test/models"
	"dealls-test/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateAttendancePeriod_Success(t *testing.T) {
	// Inisialisasi DB test
	config.InitDB()

	username := fmt.Sprintf("testadmin_%d", time.Now().UnixNano())
	rawPassword := "secret123"
	hashedPassword, err := utils.HashPassword(rawPassword)
	assert.NoError(t, err)

	admin := models.User{
		Username: username,
		Password: hashedPassword,
		Role:     "admin",
	}
	if err := config.DB.Create(&admin).Error; err != nil {
		t.Fatalf("failed to create admin user: %v", err)
	}
	defer config.DB.Delete(&admin)

	token, err := utils.GenerateJWT(admin.ID, admin.Role)
	assert.NoError(t, err)

	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, 7)

	requestBody, _ := json.Marshal(map[string]string{
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
	})

	// ✅ Simpan X-Request-ID yang akan diuji
	testRequestID := uuid.New().String()

	req := httptest.NewRequest(http.MethodPost, "/attendance-periods", bytes.NewBuffer(requestBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Request-ID", testRequestID)
	req.RemoteAddr = "127.0.0.1:1234"

	rr := httptest.NewRecorder()
	handlers.CreateAttendancePeriod(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Attendance period created", response["message"])

	var createdID uint
	if data, ok := response["data"].(map[string]interface{}); ok {
		if id, ok := data["id"].(float64); ok {
			createdID = uint(id)
		}
	}

	// ✅ Cek bahwa audit log dengan request_id tersebut tersimpan
	var auditLog models.AuditLog
	err = config.DB.Where("request_id = ? AND table_name = ? AND record_id = ?", testRequestID, "attendance_periods", createdID).First(&auditLog).Error
	assert.NoError(t, err, "AuditLog with request_id should exist")
	assert.Equal(t, "CREATE", auditLog.Action)
	assert.Equal(t, testRequestID, auditLog.RequestID)

	// Cleanup data test
	config.DB.Delete(&models.AttendancePeriod{}, createdID)
	config.DB.Delete(&auditLog)
}

