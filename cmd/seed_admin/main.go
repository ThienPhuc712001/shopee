package main

import (
	"ecommerce/internal/domain/model"
	"ecommerce/pkg/config"
	"ecommerce/pkg/database"
	"ecommerce/pkg/password"
	"fmt"
	"log"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	fmt.Println("==============================================")
	fmt.Println("       E-Commerce Admin Account Seeder       ")
	fmt.Println("==============================================")

	// Admin account details
	adminEmail := "admin@example.com"
	adminPassword := "Admin@123"
	adminFirstName := "Super"
	adminLastName := "Admin"
	adminPhone := "0900000001"

	// Check if admin already exists
	var existingAdmin model.User
	err = db.Where("email = ?", adminEmail).First(&existingAdmin).Error
	if err == nil {
		fmt.Printf("\n⚠️  Admin account already exists!\n")
		fmt.Printf("   Email: %s\n", existingAdmin.Email)
		fmt.Printf("   Role: %s\n", existingAdmin.Role)
		fmt.Printf("\n💡 To reset password, run: UPDATE users SET password = '<new_hash>' WHERE email = '%s';\n", adminEmail)
		return
	}

	// Hash password
	hashedPassword, err := password.Hash(adminPassword)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Create admin user
	admin := &model.User{
		Email:     adminEmail,
		Password:  hashedPassword,
		FirstName: adminFirstName,
		LastName:  adminLastName,
		Phone:     adminPhone,
		Role:      model.RoleUserAdmin,
		Status:    model.StatusActive,
		EmailVerified: true,
	}

	// Insert into database
	result := db.Create(admin)
	if result.Error != nil {
		log.Fatalf("Failed to create admin: %v", result.Error)
	}

	fmt.Println("\n✅ Admin account created successfully!")
	fmt.Println("\n📋 Login Credentials:")
	fmt.Println("   ┌─────────────────────────────────────┐")
	fmt.Printf("   │ Email:    %-28s│\n", adminEmail)
	fmt.Printf("   │ Password: %-28s│\n", adminPassword)
	fmt.Println("   └─────────────────────────────────────┘")
	fmt.Println("\n🔐 API Login:")
	fmt.Println("   POST http://localhost:8080/api/auth/login")
	fmt.Println("   Body: {\"email\": \"admin@example.com\", \"password\": \"Admin@123\"}")
	fmt.Println("\n💡 Remember to save these credentials securely!")
	fmt.Println("==============================================")
}
