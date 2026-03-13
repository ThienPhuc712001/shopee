package main

import (
	"ecommerce/internal/domain/model"
	"ecommerce/pkg/config"
	"ecommerce/pkg/database"
	"ecommerce/pkg/password"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("🌱 Starting database seeding...")

	// Seed users
	seedUsers(db)

	// Seed categories
	seedCategories(db)

	// Seed shops
	seedShops(db)

	// Seed products
	seedProducts(db)

	// Seed coupons
	seedCoupons(db)

	fmt.Println("✅ Database seeding completed successfully!")
}

func seedUsers(db *gorm.DB) {
	fmt.Println("📧 Seeding users...")

	users := []model.User{
		{
			Email:         "admin@example.com",
			Password:      mustHashPassword("admin123"),
			FirstName:     "Super",
			LastName:      "Admin",
			Phone:         "0900000001",
			Role:          model.RoleUserAdmin,
			Status:        model.StatusActive,
			EmailVerified: true,
		},
		{
			Email:         "customer@example.com",
			Password:      mustHashPassword("customer123"),
			FirstName:     "John",
			LastName:      "Doe",
			Phone:         "0900000002",
			Role:          model.RoleCustomer,
			Status:        model.StatusActive,
			EmailVerified: true,
		},
		{
			Email:         "seller@example.com",
			Password:      mustHashPassword("seller123"),
			FirstName:     "Jane",
			LastName:      "Smith",
			Phone:         "0900000003",
			Role:          model.RoleSeller,
			Status:        model.StatusActive,
			EmailVerified: true,
		},
	}

	for i := range users {
		if err := db.Create(&users[i]).Error; err != nil {
			log.Printf("Warning: Could not create user %s: %v", users[i].Email, err)
		} else {
			fmt.Printf("  ✓ Created user: %s (%s)\n", users[i].Email, users[i].Role)
		}
	}
}

func seedCategories(db *gorm.DB) {
	fmt.Println("📂 Seeding categories...")

	categories := []model.Category{
		{
			Name:        "Electronics",
			Slug:        "electronics",
			Description: "Electronic devices and accessories",
			Level:       0,
			SortOrder:   1,
			IsActive:    true,
		},
		{
			Name:        "Fashion",
			Slug:        "fashion",
			Description: "Clothing and fashion accessories",
			Level:       0,
			SortOrder:   2,
			IsActive:    true,
		},
		{
			Name:        "Home & Living",
			Slug:        "home-living",
			Description: "Home decor and living essentials",
			Level:       0,
			SortOrder:   3,
			IsActive:    true,
		},
		{
			Name:        "Books",
			Slug:        "books",
			Description: "Books and publications",
			Level:       0,
			SortOrder:   4,
			IsActive:    true,
		},
		{
			Name:        "Sports",
			Slug:        "sports",
			Description: "Sports equipment and accessories",
			Level:       0,
			SortOrder:   5,
			IsActive:    true,
		},
		{
			Name:        "Beauty",
			Slug:        "beauty",
			Description: "Beauty and personal care products",
			Level:       0,
			SortOrder:   6,
			IsActive:    true,
		},
		{
			Name:        "Toys & Games",
			Slug:        "toys-games",
			Description: "Toys and games for all ages",
			Level:       0,
			SortOrder:   7,
			IsActive:    true,
		},
		{
			Name:        "Watches",
			Slug:        "watches",
			Description: "Watches and accessories",
			Level:       0,
			SortOrder:   8,
			IsActive:    true,
		},
	}

	for i := range categories {
		if err := db.Create(&categories[i]).Error; err != nil {
			log.Printf("Warning: Could not create category %s: %v", categories[i].Name, err)
		} else {
			fmt.Printf("  ✓ Created category: %s\n", categories[i].Name)
		}
	}
}

func seedShops(db *gorm.DB) {
	fmt.Println("🏪 Seeding shops...")

	// Find seller user
	var seller model.User
	if err := db.Where("email = ?", "seller@example.com").First(&seller).Error; err != nil {
		log.Printf("Warning: Seller user not found: %v", err)
		return
	}

	shops := []model.Shop{
		{
			UserID:      seller.ID,
			Name:        "Tech Store",
			Slug:        "tech-store",
			Description: "Your one-stop shop for electronics and gadgets",
			Phone:       "0900000003",
			Email:       "techstore@example.com",
			Status:      model.ShopStatusActive,
		},
		{
			UserID:      seller.ID,
			Name:        "Fashion Hub",
			Slug:        "fashion-hub",
			Description: "Trendy fashion and accessories",
			Phone:       "0900000004",
			Email:       "fashionhub@example.com",
			Status:      model.ShopStatusActive,
		},
	}

	for i := range shops {
		if err := db.Create(&shops[i]).Error; err != nil {
			log.Printf("Warning: Could not create shop %s: %v", shops[i].Name, err)
		} else {
			fmt.Printf("  ✓ Created shop: %s\n", shops[i].Name)
		}
	}
}

func seedProducts(db *gorm.DB) {
	fmt.Println("📦 Seeding products...")

	// Find categories
	var categories []model.Category
	db.Find(&categories)

	categoryMap := make(map[string]model.Category)
	for _, cat := range categories {
		categoryMap[cat.Slug] = cat
	}

	// Find a shop
	var shop model.Shop
	if err := db.Where("name = ?", "Tech Store").First(&shop).Error; err != nil {
		log.Printf("Warning: Shop not found: %v", err)
		return
	}

	products := []model.Product{
		{
			ShopID:      shop.ID,
			CategoryID:  categoryMap["electronics"].ID,
			Name:        "Premium Wireless Headphones",
			Slug:        "premium-wireless-headphones",
			Description: "High-quality wireless headphones with noise cancellation and 30-hour battery life.",
			ShortDescription: "Premium wireless headphones with noise cancellation",
			Price:       199.99,
			OriginalPrice: 299.99,
			Stock:       50,
			Brand:       "AudioTech",
			Status:      "active",
			IsFeatured:  true,
		},
		{
			ShopID:      shop.ID,
			CategoryID:  categoryMap["electronics"].ID,
			Name:        "Smart Watch Pro",
			Slug:        "smart-watch-pro",
			Description: "Advanced smartwatch with health monitoring, GPS, and 7-day battery life.",
			ShortDescription: "Advanced smartwatch with health monitoring",
			Price:       349.99,
			OriginalPrice: 449.99,
			Stock:       30,
			Brand:       "TechWear",
			Status:      "active",
			IsFeatured:  true,
		},
		{
			ShopID:      shop.ID,
			CategoryID:  categoryMap["electronics"].ID,
			Name:        "Professional Camera",
			Slug:        "professional-camera",
			Description: "Professional DSLR camera with 24MP sensor and 4K video recording.",
			ShortDescription: "Professional DSLR camera with 24MP sensor",
			Price:       899.99,
			OriginalPrice: 1199.99,
			Stock:       15,
			Brand:       "PhotoPro",
			Status:      "active",
			IsFeatured:  true,
		},
		{
			ShopID:      shop.ID,
			CategoryID:  categoryMap["electronics"].ID,
			Name:        "Gaming Laptop",
			Slug:        "gaming-laptop",
			Description: "High-performance gaming laptop with RTX graphics and 144Hz display.",
			ShortDescription: "High-performance gaming laptop",
			Price:       1299.99,
			OriginalPrice: 1599.99,
			Stock:       8,
			Brand:       "GameMaster",
			Status:      "active",
			IsFeatured:  true,
		},
		{
			ShopID:      shop.ID,
			CategoryID:  categoryMap["electronics"].ID,
			Name:        "Bluetooth Speaker",
			Slug:        "bluetooth-speaker",
			Description: "Portable Bluetooth speaker with 360° sound and 20-hour battery.",
			ShortDescription: "Portable Bluetooth speaker with 360° sound",
			Price:       79.99,
			OriginalPrice: 129.99,
			Stock:       100,
			Brand:       "SoundWave",
			Status:      "active",
			IsFeatured:  false,
		},
		{
			ShopID:      shop.ID,
			CategoryID:  categoryMap["electronics"].ID,
			Name:        "Mechanical Keyboard",
			Slug:        "mechananical-keyboard",
			Description: "RGB mechanical keyboard with Cherry MX switches.",
			ShortDescription: "RGB mechanical keyboard",
			Price:       129.99,
			OriginalPrice: 179.99,
			Stock:       60,
			Brand:       "KeyMaster",
			Status:      "active",
			IsFeatured:  false,
		},
		{
			ShopID:      shop.ID,
			CategoryID:  categoryMap["electronics"].ID,
			Name:        "USB-C Hub",
			Slug:        "usb-c-hub",
			Description: "7-in-1 USB-C hub with HDMI, USB 3.0, and SD card reader.",
			ShortDescription: "7-in-1 USB-C hub",
			Price:       49.99,
			OriginalPrice: 79.99,
			Stock:       150,
			Brand:       "ConnectPro",
			Status:      "active",
			IsFeatured:  false,
		},
		{
			ShopID:      shop.ID,
			CategoryID:  categoryMap["electronics"].ID,
			Name:        "Wireless Charger",
			Slug:        "wireless-charger",
			Description: "Fast wireless charger compatible with all Qi-enabled devices.",
			ShortDescription: "Fast wireless charger",
			Price:       39.99,
			OriginalPrice: 59.99,
			Stock:       200,
			Brand:       "ChargeFast",
			Status:      "active",
			IsFeatured:  false,
		},
	}

	for i := range products {
		if err := db.Create(&products[i]).Error; err != nil {
			log.Printf("Warning: Could not create product %s: %v", products[i].Name, err)
		} else {
			fmt.Printf("  ✓ Created product: %s\n", products[i].Name)
		}
	}
}

func seedCoupons(db *gorm.DB) {
	fmt.Println("🎫 Seeding coupons...")

	now := time.Now()

	coupons := []model.Coupon{
		{
			Code:           "WELCOME10",
			DiscountType:   "percentage",
			DiscountValue:  10,
			MinOrderValue:  50,
			MaxDiscount:    50,
			StartDate:      now,
			EndDate:        now.AddDate(0, 3, 0),
			UsageLimit:     1000,
			UsedCount:      0,
			IsActive:       true,
		},
		{
			Code:           "SUMMER20",
			DiscountType:   "percentage",
			DiscountValue:  20,
			MinOrderValue:  100,
			MaxDiscount:    100,
			StartDate:      now,
			EndDate:        now.AddDate(0, 2, 0),
			UsageLimit:     500,
			UsedCount:      0,
			IsActive:       true,
		},
		{
			Code:           "FREESHIP",
			DiscountType:   "fixed",
			DiscountValue:  5.99,
			MinOrderValue:  30,
			StartDate:      now,
			EndDate:        now.AddDate(0, 1, 0),
			UsageLimit:     2000,
			UsedCount:      0,
			IsActive:       true,
		},
	}

	for i := range coupons {
		if err := db.Create(&coupons[i]).Error; err != nil {
			log.Printf("Warning: Could not create coupon %s: %v", coupons[i].Code, err)
		} else {
			fmt.Printf("  ✓ Created coupon: %s\n", coupons[i].Code)
		}
	}
}

func mustHashPassword(pwd string) string {
	hashed, err := password.Hash(pwd)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	return hashed
}
