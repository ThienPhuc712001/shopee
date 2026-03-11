// Package testutil provides testing utilities and helpers
package testutil

import (
	"ecommerce/internal/domain/model"
	"fmt"
	"time"
)

// ============================================================================
// TEST DATA FACTORIES
// ============================================================================

// CreateTestUser creates a test user with default values
func CreateTestUser(id uint, email string) *model.User {
	return &model.User{
		ID:        id,
		Email:     email,
		Password:  "$2a$10$hashedpassword",
		FirstName: "Test",
		LastName:  "User",
		Phone:     "123456789",
		Role:      model.RoleCustomer,
		Status:    model.StatusActive,
		CreatedAt: time.Now().AddDate(0, -1, 0),
		UpdatedAt: time.Now(),
	}
}

// CreateTestUserWithPassword creates a test user with hashed password
func CreateTestUserWithPassword(id uint, email, password string) *model.User {
	user := CreateTestUser(id, email)
	user.Password = password // Should be hashed in real usage
	return user
}

// CreateTestAdmin creates a test admin user
func CreateTestAdmin(id uint, email string) *model.User {
	user := CreateTestUser(id, email)
	user.Role = model.UserRole(model.AdminRoleType("admin"))
	return user
}

// CreateTestSeller creates a test seller user
func CreateTestSeller(id uint, email string) *model.User {
	user := CreateTestUser(id, email)
	user.Role = model.RoleSeller
	return user
}

// CreateTestShop creates a test shop
func CreateTestShop(id uint, userID uint, name string) *model.Shop {
	return &model.Shop{
		ID:          id,
		UserID:      userID,
		Name:        name,
		Slug:        fmt.Sprintf("shop-%d", id),
		Status:      "active",
		Description: "Test shop description",
		CreatedAt:   time.Now().AddDate(0, -1, 0),
		UpdatedAt:   time.Now(),
	}
}

// CreateTestCategory creates a test category
func CreateTestCategory(id uint, parentID *uint, name string) *model.Category {
	return &model.Category{
		ID:        id,
		ParentID:  parentID,
		Name:      name,
		Slug:      fmt.Sprintf("category-%d", id),
		Level:     1,
		IsActive:  true,
		CreatedAt: time.Now(),
	}
}

// CreateTestProduct creates a test product
func CreateTestProduct(id uint, shopID, categoryID uint, price float64) *model.Product {
	return &model.Product{
		ID:          id,
		ShopID:      shopID,
		CategoryID:  categoryID,
		Name:        fmt.Sprintf("Test Product %d", id),
		Slug:        fmt.Sprintf("test-product-%d", id),
		Description: "Test product description",
		Price:       price,
		Stock:       100,
		Status:      "active",
		CreatedAt:   time.Now().AddDate(0, -1, 0),
		UpdatedAt:   time.Now(),
	}
}

// CreateTestProductWithStock creates a test product with specific stock
func CreateTestProductWithStock(id uint, shopID, categoryID uint, price float64, stock int) *model.Product {
	product := CreateTestProduct(id, shopID, categoryID, price)
	product.Stock = stock
	return product
}

// CreateTestCart creates a test cart
func CreateTestCart(id uint, userID uint) *model.Cart {
	return &model.Cart{
		ID:        id,
		UserID:    userID,
		TotalItems: 0,
		Subtotal:  0,
		Total:     0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// CreateTestCartItem creates a test cart item
func CreateTestCartItem(id uint, cartID, productID uint, quantity int, price float64) *model.CartItem {
	return &model.CartItem{
		ID:         id,
		CartID:     cartID,
		ProductID:  productID,
		Quantity:   quantity,
		Price:      price,
		Subtotal:   price * float64(quantity),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// CreateTestOrder creates a test order
func CreateTestOrder(id uint, userID, shopID uint, total float64) *model.Order {
	return &model.Order{
		ID:           id,
		OrderNumber:  GenerateOrderNumber(),
		UserID:       userID,
		ShopID:       shopID,
		Status:       model.OrderStatus("pending"),
		Subtotal:     total,
		TotalAmount:  total,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// CreateTestOrderWithStatus creates a test order with specific status
func CreateTestOrderWithStatus(id uint, userID, shopID uint, total float64, status model.OrderStatus) *model.Order {
	order := CreateTestOrder(id, userID, shopID, total)
	order.Status = status
	return order
}

// CreateTestOrderItem creates a test order item
func CreateTestOrderItem(id uint, orderID, productID uint, quantity int, price float64) *model.OrderItem {
	return &model.OrderItem{
		ID:         id,
		OrderID:    orderID,
		ProductID:  productID,
		Quantity:   quantity,
		Price:      price,
		Subtotal:   price * float64(quantity),
		CreatedAt:  time.Now(),
	}
}

// CreateTestCoupon creates a test coupon
func CreateTestCoupon(id uint, code string, discountType model.DiscountType, discountValue float64) *model.Coupon {
	return &model.Coupon{
		ID:            id,
		Code:          code,
		Name:          fmt.Sprintf("Coupon %s", code),
		DiscountType:  discountType,
		DiscountValue: discountValue,
		MinOrderValue: 0,
		UsageLimit:    100,
		UsedCount:     0,
		IsActive:      true,
		StartDate:     time.Now().AddDate(0, -1, 0),
		EndDate:       time.Now().AddDate(1, 0, 0),
		CreatedAt:     time.Now(),
	}
}

// CreateTestPayment creates a test payment
func CreateTestPayment(id uint, orderID uint, amount float64, status model.PaymentStatus) *model.Payment {
	return &model.Payment{
		ID:         id,
		OrderID:    orderID,
		Amount:     amount,
		Status:     status,
		PaymentMethod: "credit_card",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// CreateTestAddress creates a test address
func CreateTestAddress(id uint, userID uint) *model.Address {
	return &model.Address{
		ID:       id,
		UserID:   userID,
		Name:     "Test User",
		Phone:    "123456789",
		Street:   "123 Test Street",
		Ward:     "Ward 1",
		District: "District 1",
		City:     "Ho Chi Minh City",
		Country:  "Vietnam",
		IsDefault: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// CreateTestReview creates a test review
func CreateTestReview(id uint, userID, productID uint, rating int) *model.Review {
	return &model.Review{
		ID:        id,
		UserID:    userID,
		ProductID: productID,
		Rating:    rating,
		Comment:   "Test review comment",
		IsApproved: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// CreateTestRefreshToken creates a test refresh token
func CreateTestRefreshToken(id uint, userID uint, token string, expiresAt time.Time) *model.RefreshToken {
	return &model.RefreshToken{
		ID:        int64(id),
		UserID:    int64(userID),
		Token:     token,
		ExpiresAt: expiresAt,
		Revoked:   false,
		CreatedAt: time.Now(),
	}
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// GenerateOrderNumber generates a unique order number
func GenerateOrderNumber() string {
	return fmt.Sprintf("ORD-%s", time.Now().Format("20060102150405"))
}

// GenerateEmail generates a unique test email
func GenerateEmail(prefix string, id int) string {
	return fmt.Sprintf("%s%d@test.com", prefix, id)
}

// Ptr returns a pointer to a value
func Ptr[T any](v T) *T {
	return &v
}

// TimePtr returns a pointer to a time.Time value
func TimePtr(t time.Time) *time.Time {
	return &t
}

// ============================================================================
// COMPARISON HELPERS
// ============================================================================

// FloatEquals checks if two floats are equal within a tolerance
func FloatEquals(a, b, tolerance float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff < tolerance
}

// TimeEquals checks if two times are equal within a tolerance
func TimeEquals(a, b time.Time, tolerance time.Duration) bool {
	diff := a.Sub(b)
	if diff < 0 {
		diff = -diff
	}
	return diff < tolerance
}
