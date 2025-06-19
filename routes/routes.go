package routes

import (
	"dealls-test/handlers"
	"dealls-test/middleware"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()

	// Admin attendance period
	api.HandleFunc("/admin/attendance-periods", middleware.AdminOnly(handlers.CreateAttendancePeriod)).Methods("POST")

	// Employee submit attendance
	api.HandleFunc("/attendance/submit", middleware.JWTAuthMiddleware(handlers.SubmitAttendance)).Methods("POST")

	// Auth
	api.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	api.HandleFunc("/attendance/overtime", middleware.JWTAuthMiddleware(handlers.SubmitOvertime)).Methods("POST")

	api.HandleFunc("/employee/reimbursements", handlers.SubmitReimbursement).Methods("POST")

	api.HandleFunc("/admin/run-payroll", middleware.AdminOnly(handlers.RunPayroll)).Methods("POST")

	// Endpoint untuk melihat payslip karyawan
	api.HandleFunc("/employee/payslips/generate", middleware.JWTAuthMiddleware(handlers.GeneratePayslip)).Methods("POST")

	api.HandleFunc("/employee/payslips", middleware.JWTAuthMiddleware(handlers.GetMyPayslips)).Methods("GET")

	api.HandleFunc("/admin/payslips/summary", middleware.AdminOnly(handlers.GetPayslipSummary)).Methods("GET")

	// api.HandleFunc("/employee/payslips", middleware.JWTAuthMiddleware(handlers.GetMyPayslips)).Methods("GET")

	// api.HandleFunc("/employee/generate-payslip", middleware.JWTAuthMiddleware(handlers.GeneratePayslip)).Methods("POST")

	// api.HandleFunc("/admin/payslips/summary", middleware.AdminOnly(handlers.GetPayslipSummary)).Methods("GET")

	return r
}
