package service

import (
	"ecommerce/internal/domain/model"
	"ecommerce/internal/repository"
	"errors"
	"fmt"
	"time"
)

// OrderServiceEnhanced defines the enhanced order service interface
type OrderServiceEnhanced interface {
	// Checkout
	CheckoutCart(userID uint, input *model.OrderInput) (*model.Order, error)

	// Order Management
	CreateOrder(userID uint, input *model.OrderInput) (*model.Order, error)
	GetOrderByID(id uint) (*model.Order, error)
	GetOrderByOrderNumber(orderNumber string) (*model.Order, error)
	GetUserOrders(userID uint, page, limit int) ([]model.Order, int64, error)
	GetShopOrders(shopID uint, page, limit int) ([]model.Order, int64, error)

	// Order Status
	UpdateOrderStatus(orderID uint, status model.OrderStatus, userID uint, ipAddress string) (*model.Order, error)
	CancelOrder(orderID uint, reason string, userID uint) (*model.Order, error)
	ConfirmOrder(orderID uint, userID uint) (*model.Order, error)
	ShipOrder(orderID uint, trackingNumber string, carrierName string, userID uint) (*model.Order, error)
	CompleteOrder(orderID uint) (*model.Order, error)

	// Inventory Management
	LockInventory(orderID uint) error
	ReleaseInventory(orderID uint) error
	DecreaseInventory(orderID uint) error

	// Order Calculations
	CalculateOrderTotal(orderID uint) (float64, error)
	CalculateShippingFee(shippingInfo *model.ShippingInfo, subtotal float64) float64

	// Order Tracking
	AddTrackingEvent(orderID uint, status string, location, description string) error
	GetOrderTracking(orderID uint) ([]model.OrderTracking, error)

	// Analytics
	GetOrderStatistics(userID uint) (*model.OrderStats, error)
	GetRecentOrders(limit int) ([]model.Order, error)
}

type orderServiceEnhanced struct {
	orderRepo   repository.OrderRepositoryEnhanced
	cartRepo    repository.CartRepositoryEnhanced
	productRepo repository.ProductRepositoryEnhanced
}

// NewOrderServiceEnhanced creates a new enhanced order service
func NewOrderServiceEnhanced(
	orderRepo repository.OrderRepositoryEnhanced,
	cartRepo repository.CartRepositoryEnhanced,
	productRepo repository.ProductRepositoryEnhanced,
) OrderServiceEnhanced {
	return &orderServiceEnhanced{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

// ==================== CHECKOUT ====================

func (s *orderServiceEnhanced) CheckoutCart(userID uint, input *model.OrderInput) (*model.Order, error) {
	// Validate input
	if err := s.validateCheckoutInput(input); err != nil {
		return nil, err
	}

	// Get cart
	cart, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return nil, ErrCartEmpty
	}

	// Validate cart items
	if len(cart.Items) == 0 {
		return nil, ErrCartEmpty
	}

	// Validate inventory for all items
	if err := s.validateCartInventory(cart); err != nil {
		return nil, err
	}

	// Group items by shop (for multi-shop orders)
	shopItems := make(map[uint][]model.CartItem)
	for _, item := range cart.Items {
		shopItems[item.ShopID] = append(shopItems[item.ShopID], item)
	}

	// Create order for each shop
	var mainOrder *model.Order
	for shopID, items := range shopItems {
		order, err := s.createSingleOrder(userID, shopID, items, input)
		if err != nil {
			// Rollback any created orders
			if mainOrder != nil {
				s.orderRepo.UpdateOrderStatus(mainOrder.ID, model.OrderStatusCancelled)
			}
			return nil, err
		}

		if mainOrder == nil {
			mainOrder = order
		}
	}

	// Clear cart after successful order
	_ = s.cartRepo.ClearCart(cart.ID)
	_ = s.cartRepo.UpdateCartTotals(cart.ID)

	// Lock inventory for all orders
	if err := s.LockInventory(mainOrder.ID); err != nil {
		return nil, ErrInventoryLockFailed
	}

	return mainOrder, nil
}

func (s *orderServiceEnhanced) validateCheckoutInput(input *model.OrderInput) error {
	if input == nil {
		return errors.New("checkout input is required")
	}

	if input.ShippingInfo.Name == "" || input.ShippingInfo.Phone == "" || input.ShippingInfo.Address == "" {
		return ErrInvalidShippingInfo
	}

	if input.PaymentMethod == "" {
		return ErrInvalidPaymentMethod
	}

	return nil
}

func (s *orderServiceEnhanced) validateCartInventory(cart *model.Cart) error {
	for _, item := range cart.Items {
		// Check product exists and is active
		product, err := s.productRepo.FindByID(item.ProductID)
		if err != nil {
			return fmt.Errorf("product %d not found: %w", item.ProductID, ErrProductUnavailable)
		}

		if product.Status != model.ProductStatusActive {
			return fmt.Errorf("product %s is no longer available: %w", product.Name, ErrProductUnavailable)
		}

		// Check stock
		availableStock := product.Stock
		if item.VariantID != nil {
			variants, _ := s.productRepo.FindVariantsByProductID(item.ProductID)
			for _, v := range variants {
				if uint(v.ID) == *item.VariantID {
					availableStock = v.Stock
					break
				}
			}
		}

		if availableStock < item.Quantity {
			return fmt.Errorf("insufficient stock for %s: %w", product.Name, ErrInsufficientStock)
		}
	}

	return nil
}

func (s *orderServiceEnhanced) createSingleOrder(userID, shopID uint, cartItems []model.CartItem, input *model.OrderInput) (*model.Order, error) {
	// Calculate totals
	var subtotal float64
	for _, item := range cartItems {
		subtotal += item.Subtotal
	}

	shippingFee := s.CalculateShippingFee(&input.ShippingInfo, subtotal)
	voucherDiscount := 0.0 // In production, validate and apply voucher
	taxAmount := 0.0       // In production, calculate tax based on location
	totalAmount := subtotal + shippingFee - voucherDiscount + taxAmount

	// Create order
	order := &model.Order{
		UserID:              userID,
		ShopID:              shopID,
		Status:              model.OrderStatusPending,
		PaymentStatus:       model.PaymentStatusPending,
		FulfillmentStatus:   model.FulfillmentUnfulfilled,
		Subtotal:            subtotal,
		ShippingFee:         shippingFee,
		ShippingDiscount:    0,
		ProductDiscount:     0,
		VoucherDiscount:     voucherDiscount,
		TaxAmount:           taxAmount,
		TotalAmount:         totalAmount,
		PaidAmount:          0,
		ShippingName:        input.ShippingInfo.Name,
		ShippingPhone:       input.ShippingInfo.Phone,
		ShippingAddress:     input.ShippingInfo.Address,
		ShippingWard:        input.ShippingInfo.Ward,
		ShippingDistrict:    input.ShippingInfo.District,
		ShippingCity:        input.ShippingInfo.City,
		ShippingState:       input.ShippingInfo.State,
		ShippingCountry:     input.ShippingInfo.Country,
		ShippingPostalCode:  input.ShippingInfo.PostalCode,
		ShippingMethod:      input.ShippingMethod,
		EstimatedDelivery:   input.EstimatedDelivery,
		BuyerNote:           input.BuyerNote,
	}

	// Create order items
	var orderItems []*model.OrderItem
	for _, cartItem := range cartItems {
		product, _ := s.productRepo.FindByID(cartItem.ProductID)

		productImage := ""
		if len(product.Images) > 0 {
			productImage = product.Images[0].URL
		}

		orderItem := &model.OrderItem{
			ProductID:    cartItem.ProductID,
			VariantID:    cartItem.VariantID,
			ProductName:  cartItem.ProductName,
			ProductImage: productImage,
			ProductSKU:   product.SKU,
			Quantity:     cartItem.Quantity,
			Price:        cartItem.Price,
			OriginalPrice: cartItem.OriginalPrice,
			Discount:     cartItem.Discount,
			ShopID:       cartItem.ShopID,
		}
		orderItems = append(orderItems, orderItem)
	}

	// Create order with items in transaction
	if err := s.orderRepo.CreateOrderWithItems(order, orderItems); err != nil {
		return nil, err
	}

	return order, nil
}

// ==================== ORDER MANAGEMENT ====================

func (s *orderServiceEnhanced) CreateOrder(userID uint, input *model.OrderInput) (*model.Order, error) {
	return s.CheckoutCart(userID, input)
}

func (s *orderServiceEnhanced) GetOrderByID(id uint) (*model.Order, error) {
	order, err := s.orderRepo.GetOrderByID(id)
	if err != nil {
		return nil, ErrOrderNotFound
	}
	return order, nil
}

func (s *orderServiceEnhanced) GetOrderByOrderNumber(orderNumber string) (*model.Order, error) {
	order, err := s.orderRepo.GetOrderByOrderNumber(orderNumber)
	if err != nil {
		return nil, ErrOrderNotFound
	}
	return order, nil
}

func (s *orderServiceEnhanced) GetUserOrders(userID uint, page, limit int) ([]model.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	return s.orderRepo.GetOrdersByUser(userID, limit, offset)
}

func (s *orderServiceEnhanced) GetShopOrders(shopID uint, page, limit int) ([]model.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	return s.orderRepo.GetOrdersByShop(shopID, limit, offset)
}

// ==================== ORDER STATUS ====================

func (s *orderServiceEnhanced) UpdateOrderStatus(orderID uint, status model.OrderStatus, userID uint, ipAddress string) (*model.Order, error) {
	order, err := s.orderRepo.GetOrderByID(orderID)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	// Validate status transition
	if !s.isValidStatusTransition(order.Status, status) {
		return nil, ErrInvalidOrderStatus
	}

	fromStatus := order.Status

	// Update status
	if err := s.orderRepo.UpdateOrderStatus(orderID, status); err != nil {
		return nil, err
	}

	// Add status history
	_ = s.orderRepo.AddStatusHistory(orderID, status, fromStatus, fmt.Sprintf("Status changed to %s", status), &userID, ipAddress)

	// Reload order
	return s.orderRepo.GetOrderByID(orderID)
}

func (s *orderServiceEnhanced) CancelOrder(orderID uint, reason string, userID uint) (*model.Order, error) {
	order, err := s.orderRepo.GetOrderByID(orderID)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	// Check if order can be cancelled
	if !order.CanCancel() {
		return nil, ErrOrderCannotCancel
	}

	// Update order
	order.Status = model.OrderStatusCancelled
	order.CancelReason = reason
	if err := s.orderRepo.UpdateOrder(order); err != nil {
		return nil, err
	}

	// Add status history
	_ = s.orderRepo.AddStatusHistory(orderID, model.OrderStatusCancelled, order.Status, fmt.Sprintf("Cancelled: %s", reason), &userID, "")

	// Release inventory
	_ = s.ReleaseInventory(orderID)

	return order, nil
}

func (s *orderServiceEnhanced) ConfirmOrder(orderID uint, userID uint) (*model.Order, error) {
	return s.UpdateOrderStatus(orderID, model.OrderStatusProcessing, userID, "")
}

func (s *orderServiceEnhanced) ShipOrder(orderID uint, trackingNumber string, carrierName string, userID uint) (*model.Order, error) {
	order, err := s.UpdateOrderStatus(orderID, model.OrderStatusShipped, userID, "")
	if err != nil {
		return nil, err
	}

	// Update tracking number
	order.TrackingNumber = trackingNumber
	order.ShippingCarrier = carrierName
	s.orderRepo.UpdateOrder(order)

	// Create shipping record
	shipping := &model.OrderShipping{
		OrderID:        orderID,
		TrackingNumber: trackingNumber,
		CarrierName:    carrierName,
		ShippedAt:      order.ShippedAt,
		Status:         "in_transit",
	}
	s.orderRepo.CreateOrderShipping(shipping)

	// Add tracking event
	s.AddTrackingEvent(orderID, "shipped", "", "Order has been shipped")

	return order, nil
}

func (s *orderServiceEnhanced) CompleteOrder(orderID uint) (*model.Order, error) {
	order, err := s.orderRepo.GetOrderByID(orderID)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	if order.Status != model.OrderStatusShipped {
		return nil, ErrInvalidOrderStatus
	}

	// Update status
	if err := s.orderRepo.UpdateOrderStatus(orderID, model.OrderStatusDelivered); err != nil {
		return nil, err
	}

	// Decrease inventory (permanent)
	if err := s.DecreaseInventory(orderID); err != nil {
		return nil, err
	}

	return s.orderRepo.GetOrderByID(orderID)
}

func (s *orderServiceEnhanced) isValidStatusTransition(from, to model.OrderStatus) bool {
	validTransitions := map[model.OrderStatus][]model.OrderStatus{
		model.OrderStatusPending:    {model.OrderStatusPaid, model.OrderStatusCancelled},
		model.OrderStatusPaid:       {model.OrderStatusProcessing, model.OrderStatusCancelled, model.OrderStatusRefunded},
		model.OrderStatusProcessing: {model.OrderStatusShipped, model.OrderStatusCancelled},
		model.OrderStatusShipped:    {model.OrderStatusDelivered, model.OrderStatusRefunded},
		model.OrderStatusDelivered:  {model.OrderStatusRefunded},
	}

	allowed, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, status := range allowed {
		if status == to {
			return true
		}
	}

	return false
}

// ==================== INVENTORY MANAGEMENT ====================

func (s *orderServiceEnhanced) LockInventory(orderID uint) error {
	order, err := s.orderRepo.GetOrderByID(orderID)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		if item.VariantID != nil {
			// Lock variant stock
			_ = s.productRepo.ReserveStock(item.ProductID, item.VariantID, item.Quantity)
		} else {
			// Lock product stock
			_ = s.productRepo.ReserveStock(item.ProductID, nil, item.Quantity)
		}
	}

	return nil
}

func (s *orderServiceEnhanced) ReleaseInventory(orderID uint) error {
	order, err := s.orderRepo.GetOrderByID(orderID)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		if item.VariantID != nil {
			// Release variant stock
			_ = s.productRepo.ReleaseStock(item.ProductID, item.VariantID, item.Quantity)
		} else {
			// Release product stock
			_ = s.productRepo.ReleaseStock(item.ProductID, nil, item.Quantity)
		}
	}

	return nil
}

func (s *orderServiceEnhanced) DecreaseInventory(orderID uint) error {
	order, err := s.orderRepo.GetOrderByID(orderID)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		if item.VariantID != nil {
			// Decrease variant stock
			_ = s.productRepo.DecreaseStock(item.ProductID, item.VariantID, item.Quantity)
		} else {
			// Decrease product stock
			_ = s.productRepo.DecreaseStock(item.ProductID, nil, item.Quantity)
		}
	}

	return nil
}

// ==================== ORDER CALCULATIONS ====================

func (s *orderServiceEnhanced) CalculateOrderTotal(orderID uint) (float64, error) {
	return s.orderRepo.CalculateOrderTotal(orderID)
}

func (s *orderServiceEnhanced) CalculateShippingFee(shippingInfo *model.ShippingInfo, subtotal float64) float64 {
	// Simplified shipping calculation
	// In production, calculate based on:
	// - Distance/zone
	// - Weight
	// - Shipping method
	// - Carrier rates

	if subtotal >= 500000 {
		return 0 // Free shipping for orders over 500,000
	}

	return 30000 // Default shipping fee
}

// ==================== ORDER TRACKING ====================

func (s *orderServiceEnhanced) AddTrackingEvent(orderID uint, status string, location, description string) error {
	order, err := s.orderRepo.GetOrderByID(orderID)
	if err != nil {
		return err
	}

	event := &model.OrderTracking{
		OrderID:        orderID,
		TrackingNumber: order.TrackingNumber,
		Status:         status,
		Location:       location,
		Description:    description,
		Timestamp:      time.Now(),
	}

	return s.orderRepo.AddTrackingEvent(event)
}

func (s *orderServiceEnhanced) GetOrderTracking(orderID uint) ([]model.OrderTracking, error) {
	return s.orderRepo.GetTrackingByOrderID(orderID)
}

// ==================== ANALYTICS ====================

func (s *orderServiceEnhanced) GetOrderStatistics(userID uint) (*model.OrderStats, error) {
	return s.orderRepo.GetOrderStatistics(userID)
}

func (s *orderServiceEnhanced) GetRecentOrders(limit int) ([]model.Order, error) {
	return s.orderRepo.GetRecentOrders(limit)
}
