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

func SubmitReimbursement(w http.ResponseWriter, r *http.Request) {
    ip := utils.GetIP(r)

    // Ambil atau generate request_id
    requestID := r.Header.Get("X-Request-ID")
    if requestID == "" {
        requestID = uuid.New().String()
    }

    // Parse token
    tokenStr := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
    claims, err := utils.ParseToken(tokenStr)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var input struct {
        Amount      float64 `json:"amount"`
        Description string  `json:"description"`
    }

    if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Amount <= 0 {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    reimbursement := models.Reimbursement{
        UserID:      claims.UserID,
        Amount:      input.Amount,
        Description: input.Description,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
        CreatedBy:   claims.UserID,
        UpdatedBy:   claims.UserID,
        IPAddress:   ip,
        RequestID:   requestID, // ✅ Tambahkan request_id ke struct
    }

    if err := config.DB.Create(&reimbursement).Error; err != nil {
        http.Error(w, "Failed to submit reimbursement", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]interface{}{
        "message":     "Reimbursement submitted successfully",
        "request_id":  requestID, // ✅ Sertakan juga di respons jika perlu ditelusuri
        "data":        reimbursement,
    })
}
