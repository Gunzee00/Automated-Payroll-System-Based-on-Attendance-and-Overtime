package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"dealls-test/config"
	"dealls-test/models"
	"dealls-test/utils"

	"github.com/google/uuid"
)

func GeneratePayslip(w http.ResponseWriter, r *http.Request) {
	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = uuid.New().String()
	}

	ip := utils.GetIP(r)
	tokenStr := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	claims, err := utils.ParseToken(tokenStr)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input struct {
		AttendancePeriodID uint `json:"attendance_period_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.AttendancePeriodID == 0 {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Cek jika sudah ada payslip
	var existing models.Payslip
	if err := config.DB.Where("user_id = ? AND attendance_period_id = ?", claims.UserID, input.AttendancePeriodID).First(&existing).Error; err == nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":    "Payslip already exists",
			"data":       existing,
			"request_id": requestID,
		})
		return
	}

	// Ambil periode absensi
	var period models.AttendancePeriod
	if err := config.DB.First(&period, input.AttendancePeriodID).Error; err != nil {
		http.Error(w, "Attendance period not found", http.StatusNotFound)
		return
	}

	// Hitung kehadiran
	var attendanceCount int64
	config.DB.Table("attendances").
		Where("user_id = ? AND attendance_date BETWEEN ? AND ?", claims.UserID, period.StartDate, period.EndDate).
		Count(&attendanceCount)

	// Hitung total lembur
	var overtimeTotal float64
	config.DB.Table("overtimes").
		Select("COALESCE(SUM(hours), 0) * 100000").
		Where("user_id = ? AND overtime_date BETWEEN ? AND ?", claims.UserID, period.StartDate, period.EndDate).
		Scan(&overtimeTotal)

	// Hitung total reimbursement
	var reimbursementTotal float64
	config.DB.Table("reimbursements").
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ? AND created_at BETWEEN ? AND ?", claims.UserID, period.StartDate, period.EndDate).
		Scan(&reimbursementTotal)

	// Ambil data user
	var user models.User
	if err := config.DB.First(&user, claims.UserID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	baseSalary := 0.0
	if user.Salary != nil {
		baseSalary = *user.Salary
	}

	total := float64(attendanceCount)*baseSalary + overtimeTotal + reimbursementTotal

	// Simpan payslip
	payslip := models.Payslip{
		UserID:             user.ID,
		AttendancePeriodID: input.AttendancePeriodID,
		AttendanceCount:    int(attendanceCount),
		BaseSalary:         baseSalary,
		TotalOvertime:      overtimeTotal,
		TotalReimbursement: reimbursementTotal,
		TotalSalary:        total,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		CreatedBy:          user.ID,
		UpdatedBy:          user.ID,
		IPAddress:          ip,
		RequestID:          requestID,
	}

	if err := config.DB.Create(&payslip).Error; err != nil {
		log.Printf("[request_id: %s] Failed to insert payslip: %v", requestID, err)
		http.Error(w, "Failed to generate payslip", http.StatusInternalServerError)
		return
	}

	log.Printf("[request_id: %s] Payslip generated for user_id %d for period %d", requestID, user.ID, input.AttendancePeriodID)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Payslip generated successfully",
		"data":       payslip,
		"request_id": requestID,
	})
}
func GetMyPayslips(w http.ResponseWriter, r *http.Request) {
	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = uuid.New().String()
	}

	tokenStr := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	claims, err := utils.ParseToken(tokenStr)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var payslips []models.Payslip
	if err := config.DB.Where("user_id = ?", claims.UserID).Order("created_at desc").Find(&payslips).Error; err != nil {
		http.Error(w, "Failed to retrieve payslips", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Payslips retrieved successfully",
		"data":       payslips,
		"request_id": requestID,
	})
}
func GetPayslipSummary(w http.ResponseWriter, r *http.Request) {
	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = uuid.New().String()
	}

	tokenStr := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	claims, err := utils.ParseToken(tokenStr)
	if err != nil || claims.Role != "admin" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	type SummaryItem struct {
		UserID   uint    `json:"user_id"`
		Username string  `json:"username"`
		TotalPay float64 `json:"total_take_home_pay"`
	}

	var results []SummaryItem
	rows, err := config.DB.Raw(`
		SELECT u.id AS user_id, u.username, COALESCE(SUM(p.total_salary), 0) AS total_pay
		FROM users u
		LEFT JOIN payslips p ON p.user_id = u.id
		WHERE u.role = 'employee'
		GROUP BY u.id, u.username
	`).Rows()
	if err != nil {
		http.Error(w, "Failed to fetch payslip summary", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var totalAll float64
	for rows.Next() {
		var item SummaryItem
		if err := rows.Scan(&item.UserID, &item.Username, &item.TotalPay); err == nil {
			totalAll += item.TotalPay
			results = append(results, item)
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"summary":              results,
		"total_take_home_all":  totalAll,
		"request_id":           requestID,
	})
}
