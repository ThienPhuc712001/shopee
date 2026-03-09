package repository

import (
	"ecommerce/internal/domain/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

// OrderRepositoryEnhanced defines the enhanced interface for order data operations
type OrderRepositoryEnhanced interface {
	// Order CRUD
	CreateOrder(order *model.Order) error
	CreateOrderWithItems(order *model.Order, items []*model.OrderItem) error
	GetOrderByID(id uint) (*model.Order, error)
	GetOrderByOrderNumber(orderNumber string) (*model.Order, error)
	UpdateOrder(order *model.Order) error
	DeleteOrder(id uint) error

	// Order Queries
	GetOrdersByUser(userID uint, limit, offset int) ([]model.Order, int64, error)
	GetOrdersByShop(shopID uint, limit, offset int) ([]model.Order, int64, error)
	GetOrdersByStatus(status model.OrderStatus, limit, offset int) ([]model.Order, int64, error)
	GetAllOrders(limit, offset int) ([]model.Order, int64, error)
	FindAll(limit, offset int) ([]model.Order, int64, error)
	FindByID(id uint) (*model.Order, error)

	// Order Status
	UpdateOrderStatus(orderID uint, status model.OrderStatus) error
	UpdateStatus(orderID uint, status model.OrderStatus) error
	AddStatusHistory(orderID uint, status model.OrderStatus, fromStatus model.OrderStatus, message string, changedBy *uint, ipAddress string) error
	GetStatusHistory(orderID uint) ([]model.OrderStatusHistory, error)

	// Order Items
	CreateOrderItem(item *model.OrderItem) error
	BulkCreateOrderItems(items []*model.OrderItem) error
	GetOrderItemsByOrderID(orderID uint) ([]model.OrderItem, error)
	UpdateOrderItem(item *model.OrderItem) error

	// Order Shipping
	CreateOrderShipping(shipping *model.OrderShipping) error
	UpdateOrderShipping(shipping *model.OrderShipping) error
	GetOrderShippingByOrderID(orderID uint) (*model.OrderShipping, error)

	// Order Tracking
	AddTrackingEvent(event *model.OrderTracking) error
	GetTrackingByOrderID(orderID uint) ([]model.OrderTracking, error)

	// Order Calculations
	CalculateOrderTotal(orderID uint) (float64, error)
	GetOrderStatistics(userID uint) (*model.OrderStats, error)

	// Bulk Operations
	BulkUpdateOrderStatus(orderIDs []uint, status model.OrderStatus) error

	// Analytics
	GetOrderCountByStatus(status model.OrderStatus) (int64, error)
	GetRevenueByDateRange(startDate, endDate time.Time) (float64, error)
	GetRecentOrders(limit int) ([]model.Order, error)
}

type orderRepositoryEnhanced struct {
	db *gorm.DB
}

// NewOrderRepositoryEnhanced creates a new enhanced order repository
func NewOrderRepositoryEnhanced(db *gorm.DB) OrderRepositoryEnhanced {
	return &orderRepositoryEnhanced{db: db}
}

// ==================== ORDER CRUD ====================

func (r *orderRepositoryEnhanced) CreateOrder(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepositoryEnhanced) CreateOrderWithItems(order *model.Order, items []*model.OrderItem) error {
	tx := r.db.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// Create order
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Create order items
	for _, item := range items {
		item.OrderID = order.ID
		if err := tx.Create(item).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Add initial status history
	history := &model.OrderStatusHistory{
		OrderID:    order.ID,
		Status:     order.Status,
		FromStatus: "",
		Message:    "Order created",
	}
	if err := tx.Create(history).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *orderRepositoryEnhanced) GetOrderByID(id uint) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("Shop").
		Preload("Items", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Preload("StatusHistory", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Preload("Shipping").
		Preload("Tracking", func(db *gorm.DB) *gorm.DB {
			return db.Order("timestamp DESC")
		}).
		First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepositoryEnhanced) GetOrderByOrderNumber(orderNumber string) (*model.Order, error) {
	var order model.Order
	err := r.db.Where("order_number = ?", orderNumber).
		Preload("Shop").
		Preload("Items").
		Preload("StatusHistory").
		Preload("Shipping").
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepositoryEnhanced) UpdateOrder(order *model.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepositoryEnhanced) DeleteOrder(id uint) error {
	return r.db.Delete(&model.Order{}, id).Error
}

// ==================== ORDER QUERIES ====================

func (r *orderRepositoryEnhanced) GetOrdersByUser(userID uint, limit, offset int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	if err := r.db.Model(&model.Order{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where("user_id = ?", userID).
		Preload("Shop").
		Preload("Items", func(db *gorm.DB) *gorm.DB {
			return db.Limit(5) // Load limited items for list view
		}).
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&orders).Error

	return orders, total, err
}

func (r *orderRepositoryEnhanced) GetOrdersByShop(shopID uint, limit, offset int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	if err := r.db.Model(&model.Order{}).Where("shop_id = ?", shopID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where("shop_id = ?", shopID).
		Preload("User").
		Preload("Items").
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&orders).Error

	return orders, total, err
}

func (r *orderRepositoryEnhanced) FindAll(limit, offset int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	if err := r.db.Model(&model.Order{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("User").
		Preload("Shop").
		Preload("Items", func(db *gorm.DB) *gorm.DB {
			return db.Limit(5)
		}).
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&orders).Error

	return orders, total, err
}

func (r *orderRepositoryEnhanced) FindByID(id uint) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("Shop").
		Preload("Items", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Preload("StatusHistory", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Preload("Shipping").
		Preload("Tracking", func(db *gorm.DB) *gorm.DB {
			return db.Order("timestamp DESC")
		}).
		First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepositoryEnhanced) GetOrdersByStatus(status model.OrderStatus, limit, offset int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	if err := r.db.Model(&model.Order{}).Where("status = ?", status).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where("status = ?", status).
		Preload("User").
		Preload("Shop").
		Preload("Items").
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&orders).Error

	return orders, total, err
}

func (r *orderRepositoryEnhanced) GetAllOrders(limit, offset int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	if err := r.db.Model(&model.Order{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("User").
		Preload("Shop").
		Preload("Items", func(db *gorm.DB) *gorm.DB {
			return db.Limit(5)
		}).
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&orders).Error

	return orders, total, err
}

// ==================== ORDER STATUS ====================

func (r *orderRepositoryEnhanced) UpdateOrderStatus(orderID uint, status model.OrderStatus) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	// Set status-specific timestamps
	now := time.Now()
	switch status {
	case model.OrderStatusPaid:
		updates["paid_at"] = now
	case model.OrderStatusProcessing:
		updates["confirmed_at"] = now
	case model.OrderStatusShipped:
		updates["shipped_at"] = now
	case model.OrderStatusDelivered:
		updates["delivered_at"] = now
		updates["completed_at"] = now
	case model.OrderStatusCancelled:
		updates["cancelled_at"] = now
	}

	return r.db.Model(&model.Order{}).Where("id = ?", orderID).Updates(updates).Error
}

func (r *orderRepositoryEnhanced) UpdateStatus(orderID uint, status model.OrderStatus) error {
	return r.UpdateOrderStatus(orderID, status)
}

func (r *orderRepositoryEnhanced) AddStatusHistory(orderID uint, status model.OrderStatus, fromStatus model.OrderStatus, message string, changedBy *uint, ipAddress string) error {
	history := &model.OrderStatusHistory{
		OrderID:    orderID,
		Status:     status,
		FromStatus: fromStatus,
		Message:    message,
		ChangedBy:  changedBy,
		IPAddress:  ipAddress,
	}
	return r.db.Create(history).Error
}

func (r *orderRepositoryEnhanced) GetStatusHistory(orderID uint) ([]model.OrderStatusHistory, error) {
	var history []model.OrderStatusHistory
	err := r.db.Where("order_id = ?", orderID).
		Order("created_at DESC").
		Find(&history).Error
	return history, err
}

// ==================== ORDER ITEMS ====================

func (r *orderRepositoryEnhanced) CreateOrderItem(item *model.OrderItem) error {
	return r.db.Create(item).Error
}

func (r *orderRepositoryEnhanced) BulkCreateOrderItems(items []*model.OrderItem) error {
	return r.db.CreateInBatches(items, 50).Error
}

func (r *orderRepositoryEnhanced) GetOrderItemsByOrderID(orderID uint) ([]model.OrderItem, error) {
	var items []model.OrderItem
	err := r.db.Where("order_id = ?", orderID).
		Preload("Product").
		Preload("Variant").
		Order("created_at ASC").
		Find(&items).Error
	return items, err
}

func (r *orderRepositoryEnhanced) UpdateOrderItem(item *model.OrderItem) error {
	return r.db.Save(item).Error
}

// ==================== ORDER SHIPPING ====================

func (r *orderRepositoryEnhanced) CreateOrderShipping(shipping *model.OrderShipping) error {
	return r.db.Create(shipping).Error
}

func (r *orderRepositoryEnhanced) UpdateOrderShipping(shipping *model.OrderShipping) error {
	return r.db.Save(shipping).Error
}

func (r *orderRepositoryEnhanced) GetOrderShippingByOrderID(orderID uint) (*model.OrderShipping, error) {
	var shipping model.OrderShipping
	err := r.db.Where("order_id = ?", orderID).First(&shipping).Error
	if err != nil {
		return nil, err
	}
	return &shipping, nil
}

// ==================== ORDER TRACKING ====================

func (r *orderRepositoryEnhanced) AddTrackingEvent(event *model.OrderTracking) error {
	return r.db.Create(event).Error
}

func (r *orderRepositoryEnhanced) GetTrackingByOrderID(orderID uint) ([]model.OrderTracking, error) {
	var events []model.OrderTracking
	err := r.db.Where("order_id = ?", orderID).
		Order("timestamp DESC").
		Find(&events).Error
	return events, err
}

// ==================== ORDER CALCULATIONS ====================

func (r *orderRepositoryEnhanced) CalculateOrderTotal(orderID uint) (float64, error) {
	var total float64
	err := r.db.Model(&model.OrderItem{}).
		Where("order_id = ?", orderID).
		Select("COALESCE(SUM(final_amount), 0)").
		Scan(&total).Error
	return total, err
}

func (r *orderRepositoryEnhanced) GetOrderStatistics(userID uint) (*model.OrderStats, error) {
	stats := &model.OrderStats{}

	// Total orders
	r.db.Model(&model.Order{}).
		Where("user_id = ?", userID).
		Where("deleted_at IS NULL").
		Count(&stats.TotalOrders)

	// Total revenue
	r.db.Model(&model.Order{}).
		Where("user_id = ?", userID).
		Where("status NOT IN ?", []model.OrderStatus{model.OrderStatusCancelled, model.OrderStatusRefunded}).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&stats.TotalRevenue)

	// Pending orders
	r.db.Model(&model.Order{}).
		Where("user_id = ?", userID).
		Where("status = ?", model.OrderStatusPending).
		Count(&stats.PendingOrders)

	// Completed orders
	r.db.Model(&model.Order{}).
		Where("user_id = ?", userID).
		Where("status = ?", model.OrderStatusDelivered).
		Count(&stats.CompletedOrders)

	// Cancelled orders
	r.db.Model(&model.Order{}).
		Where("user_id = ?", userID).
		Where("status = ?", model.OrderStatusCancelled).
		Count(&stats.CancelledOrders)

	// Average order value
	if stats.TotalOrders > 0 {
		stats.AverageOrderValue = stats.TotalRevenue / float64(stats.TotalOrders)
	}

	return stats, nil
}

// ==================== BULK OPERATIONS ====================

func (r *orderRepositoryEnhanced) BulkUpdateOrderStatus(orderIDs []uint, status model.OrderStatus) error {
	return r.db.Model(&model.Order{}).
		Where("id IN ?", orderIDs).
		Update("status", status).Error
}

// ==================== ANALYTICS ====================

func (r *orderRepositoryEnhanced) GetOrderCountByStatus(status model.OrderStatus) (int64, error) {
	var count int64
	err := r.db.Model(&model.Order{}).
		Where("status = ?", status).
		Where("deleted_at IS NULL").
		Count(&count).Error
	return count, err
}

func (r *orderRepositoryEnhanced) GetRevenueByDateRange(startDate, endDate time.Time) (float64, error) {
	var revenue float64
	err := r.db.Model(&model.Order{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Where("status NOT IN ?", []model.OrderStatus{model.OrderStatusCancelled, model.OrderStatusRefunded}).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&revenue).Error
	return revenue, err
}

func (r *orderRepositoryEnhanced) GetRecentOrders(limit int) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Preload("User").
		Preload("Shop").
		Limit(limit).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

// Error definitions
var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrOrderItemNotFound  = errors.New("order item not found")
	ErrInvalidOrderStatus = errors.New("invalid order status transition")
	ErrOrderCannotCancel  = errors.New("order cannot be cancelled")
	ErrOrderCannotRefund  = errors.New("order cannot be refunded")
)
