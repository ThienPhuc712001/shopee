package service

import "errors"

// Common service errors
var (
	// Cart errors
	ErrCartNotFound        = errors.New("cart not found")
	ErrCartItemNotFound    = errors.New("cart item not found")
	ErrCartEmpty           = errors.New("cart is empty")
	
	// Product errors
	ErrProductNotFound     = errors.New("product not found")
	ErrVariantNotFound     = errors.New("variant not found")
	ErrProductNotAvailable = errors.New("product is not available")
	ErrProductUnavailable  = errors.New("product is unavailable")
	ErrInvalidPrice        = errors.New("invalid price")
	ErrInvalidSKU          = errors.New("invalid SKU")
	ErrDuplicateSKU        = errors.New("SKU already exists")
	ErrUnauthorizedSeller  = errors.New("seller can only manage own products")
	
	// Order errors
	ErrOrderNotFound       = errors.New("order not found")
	ErrOrderCannotCancel   = errors.New("order cannot be cancelled at this stage")
	ErrOrderCannotRefund   = errors.New("order cannot be refunded")
	ErrInvalidOrderStatus  = errors.New("invalid order status transition")
	
	// Payment errors
	ErrPaymentNotFound       = errors.New("payment not found")
	ErrRefundNotFound        = errors.New("refund not found")
	ErrPaymentMethodNotFound = errors.New("payment method not found")
	ErrInvalidPaymentMethod  = errors.New("invalid payment method")
	ErrPaymentAlreadyPaid    = errors.New("payment already completed")
	ErrInvalidSignature      = errors.New("invalid webhook signature")
	ErrDuplicateTransaction  = errors.New("duplicate transaction")
	ErrRefundExceedsAmount   = errors.New("refund amount exceeds payment")
	ErrInvalidRefundType     = errors.New("invalid refund type")
	ErrGatewayError          = errors.New("payment gateway error")
	ErrPaymentTimeout        = errors.New("payment timeout expired")
	
	// Product errors
	ErrTooManyImages         = errors.New("too many images (max 9)")
	ErrImageTooLarge         = errors.New("image file too large")
	ErrInvalidImageFormat    = errors.New("invalid image format")
	ErrCategoryNotFound      = errors.New("category not found")
	
	// Inventory errors
	ErrInsufficientStock     = errors.New("insufficient stock available")
	ErrInventoryLockFailed   = errors.New("failed to lock inventory")
	
	// Cart specific errors
	ErrInvalidQuantity       = errors.New("quantity must be between 1 and 999")
	ErrPriceChanged          = errors.New("product price has changed")
	
	// Admin errors
	ErrAdminNotFound         = errors.New("admin user not found")
	ErrUnauthorized          = errors.New("unauthorized action")
	ErrInsufficientPermission = errors.New("insufficient permissions")
	ErrAdminAlreadyExists    = errors.New("admin with this email already exists")
	ErrInvalidRole           = errors.New("invalid admin role")
	
	// Validation errors
	ErrInvalidShippingInfo   = errors.New("invalid shipping information")
	ErrRefundExceedsPayment  = errors.New("refund amount exceeds payment")
)
