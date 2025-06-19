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

func RunPayroll(w http.ResponseWriter, r *http.Request) {
	ip := utils.GetIP(r)

	// Ambil atau generate request_id
	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = uuid.New().String()
	}

	// Ambil token dan validasi
	tokenStr := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	claims, err := utils.ParseToken(tokenStr)
	if err != nil || claims.Role != "admin" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Decode input
	var input struct {
		AttendancePeriodID uint `json:"attendance_period_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Ambil periode kehadiran
	var period models.AttendancePeriod
	if err := config.DB.First(&period, input.AttendancePeriodID).Error; err != nil {
		http.Error(w, "Attendance period not found", http.StatusNotFound)
		return
	}

	// Ambil semua karyawan
	var users []models.User
	if err := config.DB.Where("role = ?", "employee").Find(&users).Error; err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	// Proses setiap karyawan
	for _, user := range users {
		// Skip jika payroll sudah dibuat
		var existing models.Payroll
		if err := config.DB.
			Where("user_id = ? AND attendance_period_id = ?", user.ID, input.AttendancePeriodID).
			First(&existing).Error; err == nil {
			continue
		}

		// Hitung jumlah kehadiran
		var attendanceCount int64
		config.DB.Table("attendances").
			Where("user_id = ? AND attendance_date BETWEEN ? AND ?", user.ID, period.StartDate, period.EndDate).
			Count(&attendanceCount)

		// Hitung lembur
		var overtimeTotal float64
		config.DB.Table("overtimes").
			Select("COALESCE(SUM(hours), 0) * 100000").
			Where("user_id = ? AND overtime_date BETWEEN ? AND ?", user.ID, period.StartDate, period.EndDate).
			Scan(&overtimeTotal)

		// Hitung reimbursement
		var reimbursementTotal float64
		config.DB.Table("reimbursements").
			Select("COALESCE(SUM(amount), 0)").
			Where("user_id = ? AND created_at BETWEEN ? AND ?", user.ID, period.StartDate, period.EndDate).
			Scan(&reimbursementTotal)

		// Gaji pokok
		var baseSalary float64 = 0
		if user.Salary != nil {
			baseSalary = *user.Salary
		}

		// Total take home pay
		total := float64(attendanceCount)*baseSalary + overtimeTotal + reimbursementTotal

		// Simpan ke tabel payroll
		payroll := models.Payroll{
			UserID:             user.ID,
			AttendancePeriodID: input.AttendancePeriodID,
			BaseSalary:         baseSalary,
			TotalOvertime:      overtimeTotal,
			TotalReimbursement: reimbursementTotal,
			TotalSalary:        total,
			CreatedAt:          time.Now(),
			CreatedBy:          claims.UserID,
			IPAddress:          ip,
			RequestID:          requestID,
		}
		config.DB.Create(&payroll)

		// Simpan ke tabel payslip ✅ Termasuk AttendanceCount
		payslip := models.Payslip{
			UserID:             user.ID,
			AttendancePeriodID: input.AttendancePeriodID,
			AttendanceCount:    int(attendanceCount), // ✅ Diperbaiki
			BaseSalary:         baseSalary,
			TotalOvertime:      overtimeTotal,
			TotalReimbursement: reimbursementTotal,
			TotalSalary:        total,
			CreatedAt:          time.Now(),
			CreatedBy:          claims.UserID,
			IPAddress:          ip,
			RequestID:          requestID,
		}
		config.DB.Create(&payslip)
	}

	// Kirim response
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Payroll processed successfully for attendance period",
		"request_id": requestID,
	})
}
