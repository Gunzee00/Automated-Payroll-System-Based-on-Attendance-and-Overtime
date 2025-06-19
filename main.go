package main

import (
	"log"
	"net/http"

	"dealls-test/config"
	"dealls-test/models"
	"dealls-test/routes"
)

func main() {
	// 1. Init database
	config.InitDB()
	config.InitTestDB()
	// SeedUsers()
	// SeedAdminUser()

	// 2. Auto migrate models
	if err := config.DB.AutoMigrate(&models.User{}, &models.AttendancePeriod{}); err != nil {
		log.Fatal("❌ Gagal migrate:", err)
	}

	r := routes.SetupRoutes()

	log.Println("🚀 Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("❌ Server error:", err)
	}
}
