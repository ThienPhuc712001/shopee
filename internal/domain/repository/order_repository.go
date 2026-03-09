package repository

import (
	"ecommerce/internal/domain/model"
)

// OrderRepository defines the interface for order data operations
type OrderRepository interface {
	// Create creates a new order
	Create(order *model.Order) error

	// Update updates an existing order
	Update(order *model.Order) error

	// Delete deletes an order by ID
	Delete(id uint) error

	// FindByID finds an order by ID
	FindByID(id uint) (*model.Order, error)

	// FindByOrderNumber finds an order by order number
	FindByOrderNumber(orderNumber string) (*model.Order, error)

	// FindByUserID finds all orders by user ID
	FindByUserID(userID uint, limit, offset int) ([]model.Order, int64, error)

	// FindByShopID finds all orders by shop ID
	FindByShopID(shopID uint, limit, offset int) ([]model.Order, int64, error)

	// FindByStatus finds orders by status
	FindByStatus(status model.OrderStatus, limit, offset int) ([]model.Order, int64, error)

	// FindAll retrieves all orders with pagination
	FindAll(limit, offset int) ([]model.Order, int64, error)

	// UpdateStatus updates order status
	UpdateStatus(id uint, status model.OrderStatus) error

	// UpdatePaymentStatus updates payment status
	UpdatePaymentStatus(id uint, status model.PaymentStatus) error

	// AddOrderItem adds an item to an order
	AddOrderItem(item *model.OrderItem) error

	// GetOrderItemsByOrderID gets all items in an order
	GetOrderItemsByOrderID(orderID uint) ([]model.OrderItem, error)

	// GetOrderStatistics gets order statistics for a user
	GetOrderStatistics(userID uint) (totalOrders int64, totalSpent float64, err error)
}
