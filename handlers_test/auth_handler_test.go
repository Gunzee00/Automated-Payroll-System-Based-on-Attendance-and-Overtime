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
)

func TestLoginHandler_Success(t *testing.T) {
	// Setup: buat user di DB
	password := "secret123"
	hashedPassword, _ := utils.HashPassword(password)

	user := models.User{
		Username: "testloginuser",
		Password: hashedPassword,
		Role:     "admin",
	}
	config.DB.Create(&user)
	defer config.DB.Delete(&user)

	// Prepare request JSON body
	requestBody, _ := json.Marshal(map[string]string{
		"username": user.Username,
		"password": password,
	})

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Recorder untuk response-nya
	rr := httptest.NewRecorder()

	// Jalankan handler
	handler := http.HandlerFunc(handlers.LoginHandler)
	handler.ServeHTTP(rr, req)

	// Cek status code
	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	// Cek isi body JSON
	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["token"] == nil {
		t.Error("expected token in response, got nil")
	}
	if response["username"] != user.Username {
		t.Errorf("expected username %s, got %v", user.Username, response["username"])
	}
	if response["role"] != user.Role {
		t.Errorf("expected role %s, got %v", user.Role, response["role"])
	}
}
