package handlers_test

import (
	"bytes"
	"dealls-test/config"
	"dealls-test/handlers"
	"dealls-test/models"
	"dealls-test/utils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	// "time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Float64Ptr(f float64) *float64 {
	return &f
}

func TestSubmitReimbursement_Success(t *testing.T) {
	config.InitTestDB()
	config.DB = config.TestDB

	// Buat user dummy
	user := models.User{
		Username: "user_test",
		Password: "secret",
		Role:     "employee",
		Salary:   Float64Ptr(5000000),
	}
	config.TestDB.Create(&user)
	defer config.TestDB.Delete(&user)

	// Generate token dan request ID
	token, _ := utils.GenerateJWT(user.ID, user.Role)
	requestID := uuid.New().String()

	payload := map[string]interface{}{
		"amount":      150000,
		"description": "Transport meeting",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/reimbursements", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", requestID)
	req.RemoteAddr = "127.0.0.1:1001"

	rr := httptest.NewRecorder()
	handlers.SubmitReimbursement(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&resp)
	assert.Equal(t, "Reimbursement submitted successfully", resp["message"])
	assert.NotNil(t, resp["data"])

	// Validasi reimbursement tercatat
	var reimbursement models.Reimbursement
	err := config.TestDB.Where("user_id = ?", user.ID).Order("created_at DESC").First(&reimbursement).Error
	assert.NoError(t, err)
	assert.Equal(t, 150000.0, reimbursement.Amount)

	// Validasi audit log juga tercatat
	var audit models.AuditLog
	err = config.TestDB.Where("table_name = ? AND record_id = ? AND request_id = ?", "reimbursements", reimbursement.ID, requestID).First(&audit).Error
	assert.NoError(t, err)
	assert.Equal(t, user.ID, audit.PerformedBy)

	// Clean-up
	config.TestDB.Delete(&reimbursement)
	config.TestDB.Delete(&audit)
}

func TestSubmitReimbursement_InvalidInput(t *testing.T) {
	config.InitTestDB()

	token, _ := utils.GenerateJWT(1, "employee")
	req := httptest.NewRequest("POST", "/reimbursements", bytes.NewBuffer([]byte(`{"amount": -100}`)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handlers.SubmitReimbursement(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSubmitReimbursement_Unauthorized(t *testing.T) {
	req := httptest.NewRequest("POST", "/reimbursements", bytes.NewBuffer([]byte(`{"amount": 100}`)))
	rr := httptest.NewRecorder()

	handlers.SubmitReimbursement(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
