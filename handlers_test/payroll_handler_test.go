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
func TestRunPayroll_Success(t *testing.T) {
	// Init DB
	config.InitDB()

	// Buat admin dummy
	admin := models.User{
		Username: "admin_payroll_test",
		Password: "password",
		Role:     "admin",
	}
	config.DB.Create(&admin)
	defer config.DB.Delete(&admin)

	// Buat employee dummy
	salary := 100000.0
	employee := models.User{
		Username: "employee_payroll_test",
		Password: "password",
		Role:     "employee",
		Salary:   &salary,
	}
	config.DB.Create(&employee)
	defer config.DB.Delete(&employee)

	// Buat attendance period
	start := time.Now().AddDate(0, 0, -7).Truncate(24 * time.Hour)
	end := time.Now().Truncate(24 * time.Hour)
	period := models.AttendancePeriod{
		StartDate: start,
		EndDate:   end,
		Status:    "open",
		CreatedAt: time.Now(),
		CreatedBy: admin.ID,
		UpdatedBy: admin.ID,
		IPAddress: "127.0.0.1",
	}
	config.DB.Create(&period)
	defer config.DB.Delete(&period)

	// Attendance dummy
	attendance := models.Attendance{
		UserID:         employee.ID,
		AttendanceDate: start.AddDate(0, 0, 1),
		CreatedAt:      time.Now(),
		CreatedBy:      employee.ID,
		UpdatedBy:      employee.ID,
		IPAddress:      "127.0.0.1",
	}
	config.DB.Create(&attendance)
	defer config.DB.Delete(&attendance)

	// Overtime dummy
	overtime := models.Overtime{
		UserID:       employee.ID,
		OvertimeDate: start.AddDate(0, 0, 2),
		Hours:        2,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		CreatedBy:    employee.ID,
		UpdatedBy:    employee.ID,
		IPAddress:    "127.0.0.1",
	}
	config.DB.Create(&overtime)
	defer config.DB.Delete(&overtime)

	// Reimbursement dummy
	reimbursement := models.Reimbursement{
		UserID:    employee.ID,
		Amount:    50000,
		CreatedAt: start.AddDate(0, 0, 3),
		CreatedBy: employee.ID,
		UpdatedBy: employee.ID,
		IPAddress: "127.0.0.1",
	}
	config.DB.Create(&reimbursement)
	defer config.DB.Delete(&reimbursement)

	// Generate token
	token, err := utils.GenerateJWT(admin.ID, admin.Role)
	assert.NoError(t, err)

	// Request body
	body, _ := json.Marshal(map[string]interface{}{
		"attendance_period_id": period.ID,
	})

	requestID := uuid.New().String()

	// Request
	req := httptest.NewRequest(http.MethodPost, "/run-payroll", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Request-ID", requestID)
	req.RemoteAddr = "127.0.0.1:1234"

	rr := httptest.NewRecorder()
	handlers.RunPayroll(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Payroll processed successfully for attendance period", resp["message"])

	// ✅ Validasi payroll tersimpan
	var payroll models.Payroll
	err = config.DB.Where("user_id = ? AND attendance_period_id = ?", employee.ID, period.ID).First(&payroll).Error
	assert.NoError(t, err)
	assert.Equal(t, employee.ID, payroll.UserID)

	// ✅ Validasi payslip tersimpan
	var payslip models.Payslip
	err = config.DB.Where("user_id = ? AND attendance_period_id = ?", employee.ID, period.ID).First(&payslip).Error
	assert.NoError(t, err)
	assert.Equal(t, payroll.TotalSalary, payslip.TotalSalary)
	assert.Equal(t, 1, payslip.AttendanceCount)

	// ✅ Optional: Validasi audit log (jika aktif)
	var audit models.AuditLog
	err = config.DB.Where("record_id = ? AND table_name = ? AND request_id = ?", payroll.ID, "payrolls", requestID).First(&audit).Error
	assert.NoError(t, err)

	// Cleanup
	config.DB.Delete(&payroll)
	config.DB.Delete(&payslip)
	config.DB.Where("record_id = ? AND table_name = ?", payroll.ID, "payrolls").Delete(&models.AuditLog{})
}
