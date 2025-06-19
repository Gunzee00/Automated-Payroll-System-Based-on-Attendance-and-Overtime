package config

import (
	"fmt"
	"log"
	"dealls-test/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var TestDB *gorm.DB

func InitTestDB() {
	dsn := "host=localhost user=user password=root dbname=db_test port=5432 sslmode=disable"
	var err error

	TestDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Gagal konek ke TEST database:", err)
	}

	fmt.Println("✅ Koneksi ke TEST DB berhasil!")

	// Auto migrate semua model
	err = TestDB.AutoMigrate(
		&models.AttendancePeriod{},  
		&models.Attendance{},  
		&models.AuditLog{},  
		&models.Overtime{},  
		&models.Payroll{},  
		&models.Payslip{},  
		&models.Reimbursement{},  
		&models.User{},  
		 
		// Tambahkan semua model: User, Payroll, Attendance, dll.
	)
	if err != nil {
		log.Fatal("❌ Gagal AutoMigrate di TEST DB:", err)
	}
}
