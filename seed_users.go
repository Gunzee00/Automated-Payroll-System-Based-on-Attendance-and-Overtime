package main

import (
	"fmt"
	"log"
	"math/rand"

	"dealls-test/config"
	"dealls-test/models"

	"golang.org/x/crypto/bcrypt"
)

func SeedUsers() {
	config.InitDB()

	for i := 1; i <= 100; i++ {
		username := fmt.Sprintf("employee%d", i)
		password := username // password sama dengan username

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("❌ Gagal hash password untuk %s: %v", username, err)
			continue
		}
		salary := float64(4000000 + rand.Intn(2000000))

		user := models.User{
			Username: username,
			Password: string(hashedPassword),
			Role:     "employee",
			Salary:   &salary,
		}

		result := config.DB.Create(&user)
		if result.Error != nil {
			log.Printf("❌ Gagal insert user %s: %v", username, result.Error)
		}
	}

	fmt.Println("✅ Sukses insert 100 data dummy employee.")
}
func SeedAdminUser() {
	config.InitDB()

	username := "admin"
	password := "admin"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("❌ Gagal hash password untuk admin: %v", err)
	}

	admin := models.User{
		Username: username,
		Password: string(hashedPassword),
		Role:     "admin",
		Salary:   nil, // biarkan NULL
	}

	result := config.DB.Create(&admin)
	if result.Error != nil {
		log.Fatalf("❌ Gagal insert akun admin: %v", result.Error)
	}

	fmt.Println("✅ Akun admin berhasil dibuat.")
}
