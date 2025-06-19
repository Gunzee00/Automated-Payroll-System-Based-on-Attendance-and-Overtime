package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"dealls-test/config"
	"dealls-test/handlers"
	"dealls-test/models"
	"dealls-test/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetMyPayslips_Success(t *testing.T) {
	// Init DB
	config.InitDB()

	// Buat user dummy
	user := models.User{
		Username: "testemployee",
		Password: "hashed",
		Role:     "employee",
	}
	config.DB.Create(&user)
	defer config.DB.Delete(&user)

	// Buat payslip dummy
	payslip := models.Payslip{
		UserID:             user.ID,
		AttendancePeriodID: 1,
		BaseSalary:         5000000,
		TotalOvertime:      200000,
		TotalReimbursement: 150000,
		TotalSalary:        5350000,
		AttendanceCount:    20,
		CreatedAt:          time.Now(),
		CreatedBy:          user.ID,
		IPAddress:          "127.0.0.1",
		RequestID:          uuid.New().String(),
	}
	config.DB.Create(&payslip)
	defer config.DB.Delete(&payslip)

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Role)
	assert.NoError(t, err)

	// Buat request dengan X-Request-ID
	testRequestID := uuid.New().String()
	req := httptest.NewRequest(http.MethodGet, "/api/employee/payslips", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Request-ID", testRequestID)

	// Recorder
	rr := httptest.NewRecorder()

	// Jalankan handler
	handlers.GetMyPayslips(rr, req)

	// Validasi status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Validasi isi response
	assert.Equal(t, "Payslips retrieved successfully", response["message"])
	assert.NotNil(t, response["data"])

	// Validasi bahwa request_id dikembalikan
	if val, ok := response["request_id"]; ok {
		assert.Equal(t, testRequestID, val)
	} else {
		t.Errorf("Expected request_id in response but not found")
	}
}
