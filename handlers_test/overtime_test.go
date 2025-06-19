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
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSubmitOvertime_Success(t *testing.T) {
	// Inisialisasi DB
	config.InitDB()

	// Siapkan user
	password := "secret"
	hashedPassword, _ := utils.HashPassword(password)
	user := models.User{
		Username: "testovertime",
		Password: hashedPassword,
		Role:     "employee",
	}
	config.DB.Create(&user)
	defer func() {
		config.DB.Where("user_id = ?", user.ID).Delete(&models.Overtime{})
		config.DB.Delete(&user)
	}()

	// Generate token
	token, err := utils.GenerateJWT(user.ID, user.Role)
	assert.NoError(t, err)

	// Buat body request lembur (valid)
	body, _ := json.Marshal(map[string]interface{}{
		"hours": 2,
	})

	// ✅ Generate request_id untuk diuji
	testRequestID := uuid.New().String()

	// Buat request
	req := httptest.NewRequest(http.MethodPost, "/submit-overtime", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Request-ID", testRequestID)
	req.RemoteAddr = "127.0.0.1:3000"

	rr := httptest.NewRecorder()
	handlers.SubmitOvertime(rr, req)

	// Validasi status
	assert.Equal(t, http.StatusOK, rr.Code)

	// Validasi body
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Overtime submitted successfully", response["message"])

	// Validasi bahwa overtime tersimpan dengan benar
	var overtime models.Overtime
	err = config.DB.Where("user_id = ? AND overtime_date = ?", user.ID, time.Now().Truncate(24*time.Hour)).First(&overtime).Error
	assert.NoError(t, err)
	assert.Equal(t, 2, overtime.Hours)
	assert.Equal(t, user.ID, overtime.UserID)

	// ✅ Validasi bahwa request_id tersimpan (jika kolom ini ada di model Overtime)
	if overtime.RequestID != "" {
		assert.Equal(t, testRequestID, overtime.RequestID)
	}
}
