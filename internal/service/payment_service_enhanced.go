package service

import (
	"ecommerce/internal/domain/model"
	"ecommerce/internal/repository"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// PaymentServiceEnhanced defines the enhanced payment service interface
type PaymentServiceEnhanced interface {
	// Payment Creation
	CreatePaymentIntent(userID, orderID uint, method model.PaymentMethod, provider string) (*model.PaymentResponse, error)
	CreatePayment(userID, orderID uint, input *model.PaymentInput) (*model.Payment, error)

	// Payment Processing
	ProcessPayment(paymentID uint, gatewayResponse map[string]interface{}) (*model.Payment, error)
	ConfirmPayment(transactionID string) (*model.Payment, error)
	CancelPayment(paymentID uint, reason string) (*model.Payment, error)

	// Webhook Handling
	HandleWebhook(input *model.PaymentWebhookInput) error
	VerifyWebhookSignature(payload, signature, secret string) bool

	// Refunds
	RequestRefund(userID uint, input *model.RefundInput) (*model.Refund, error)
	ProcessRefund(refundID uint) (*model.Refund, error)
	ApproveRefund(refundID uint, adminID uint) (*model.Refund, error)
	RejectRefund(refundID uint, adminID uint, reason string) (*model.Refund, error)

	// Payment Methods
	SavePaymentMethod(userID uint, input *model.PaymentMethodInput) (*model.PaymentMethodModel, error)
	GetPaymentMethods(userID uint) ([]model.PaymentMethodModel, error)
	DeletePaymentMethod(userID, methodID uint) error
	SetDefaultPaymentMethod(userID, methodID uint) error

	// Payment Status
	GetPaymentByOrderID(orderID uint) (*model.Payment, error)
	GetPaymentByTransactionID(transactionID string) (*model.Payment, error)
	GetUserPayments(userID uint, page, limit int) ([]model.Payment, int64, error)

	// Analytics
	GetPaymentStats(userID uint) (*model.PaymentStats, error)

	// Gateway Integration (mock - implement actual gateway calls)
	InitializeGateway(provider string, config *model.PaymentGatewayConfig) error
	ChargePayment(provider string, amount float64, token string) (map[string]interface{}, error)
	RefundPayment(provider string, transactionID string, amount float64) (map[string]interface{}, error)

	// Cleanup
	CleanupExpiredPayments() (int64, error)
}

type paymentServiceEnhanced struct {
	paymentRepo    repository.PaymentRepositoryEnhanced
	orderRepo      repository.OrderRepositoryEnhanced
	gatewayConfigs map[string]*model.PaymentGatewayConfig
}

// NewPaymentServiceEnhanced creates a new enhanced payment service
func NewPaymentServiceEnhanced(
	paymentRepo repository.PaymentRepositoryEnhanced,
	orderRepo repository.OrderRepositoryEnhanced,
) PaymentServiceEnhanced {
	return &paymentServiceEnhanced{
		paymentRepo:    paymentRepo,
		orderRepo:      orderRepo,
		gatewayConfigs: make(map[string]*model.PaymentGatewayConfig),
	}
}

// ==================== PAYMENT CREATION ====================

func (s *paymentServiceEnhanced) CreatePaymentIntent(userID, orderID uint, method model.PaymentMethod, provider string) (*model.PaymentResponse, error) {
	// Get order to get amount
	order, err := s.orderRepo.GetOrderByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	// Create payment record
	payment := &model.Payment{
		OrderID:         orderID,
		UserID:          userID,
		PaymentMethod:   method,
		PaymentProvider: provider,
		Amount:          order.TotalAmount,
		Currency:        "USD",
		Status:          model.PaymentStatusPending,
	}

	if err := s.paymentRepo.CreatePayment(payment); err != nil {
		return nil, err
	}

	response := &model.PaymentResponse{
		PaymentID:     payment.ID,
		TransactionID: payment.TransactionID,
		OrderID:       orderID,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		Status:        payment.Status,
		PaymentMethod: payment.PaymentMethod,
	}

	// Generate gateway-specific data
	switch provider {
	case "stripe":
		// In production: Create Stripe PaymentIntent
		response.ClientSecret = "pi_xxx_secret_xxx" // Mock
	case "paypal":
		// In production: Create PayPal order
		response.PaymentURL = "https://paypal.com/checkout/xxx" // Mock
	}

	return response, nil
}

func (s *paymentServiceEnhanced) CreatePayment(userID, orderID uint, input *model.PaymentInput) (*model.Payment, error) {
	// Get order
	order, err := s.orderRepo.GetOrderByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	// Validate payment method
	if err := s.validatePaymentMethod(input.PaymentMethod); err != nil {
		return nil, err
	}

	// Create payment
	payment := &model.Payment{
		OrderID:         orderID,
		UserID:          userID,
		PaymentMethod:   input.PaymentMethod,
		PaymentProvider: input.Provider,
		Amount:          order.TotalAmount,
		Currency:        "USD",
		Status:          model.PaymentStatusPending,
	}

	if err := s.paymentRepo.CreatePayment(payment); err != nil {
		return nil, err
	}

	// Create initial transaction record
	transaction := &model.PaymentTransaction{
		PaymentID: payment.ID,
		Type:      model.PaymentTypeCharge,
		Amount:    payment.Amount,
		Status:    "pending",
	}
	s.paymentRepo.CreateTransaction(transaction)

	// Save payment method if requested
	if input.SaveMethod {
		// In production, get token from gateway
	}

	return payment, nil
}

func (s *paymentServiceEnhanced) validatePaymentMethod(method model.PaymentMethod) error {
	validMethods := map[model.PaymentMethod]bool{
		model.PaymentMethodCreditCard:   true,
		model.PaymentMethodDebitCard:    true,
		model.PaymentMethodBankTransfer: true,
		model.PaymentMethodEWallet:      true,
		model.PaymentMethodCOD:          true,
		model.PaymentMethodPayPal:       true,
	}

	if !validMethods[method] {
		return ErrInvalidPaymentMethod
	}

	return nil
}

// ==================== PAYMENT PROCESSING ====================

func (s *paymentServiceEnhanced) ProcessPayment(paymentID uint, gatewayResponse map[string]interface{}) (*model.Payment, error) {
	payment, err := s.paymentRepo.GetPaymentByID(paymentID)
	if err != nil {
		return nil, ErrPaymentNotFound
	}

	if payment.IsPaid() {
		return nil, ErrPaymentAlreadyPaid
	}

	// Update payment with gateway response
	responseJSON, _ := json.Marshal(gatewayResponse)
	payment.GatewayResponse = string(responseJSON)
	payment.Status = model.PaymentStatusPaid

	now := time.Now()
	payment.PaidAt = &now

	if err := s.paymentRepo.UpdatePayment(payment); err != nil {
		return nil, err
	}

	// Create transaction record
	transaction := &model.PaymentTransaction{
		PaymentID:     payment.ID,
		Type:          model.PaymentTypeCharge,
		Amount:        payment.Amount,
		Status:        "completed",
		GatewayID:     fmt.Sprintf("%v", gatewayResponse["id"]),
		GatewayResponse: string(responseJSON),
		ProcessedAt:   &now,
	}
	s.paymentRepo.CreateTransaction(transaction)

	// Update order status
	s.orderRepo.UpdateOrderStatus(payment.OrderID, model.OrderStatusPaid)

	return payment, nil
}

func (s *paymentServiceEnhanced) ConfirmPayment(transactionID string) (*model.Payment, error) {
	payment, err := s.paymentRepo.GetPaymentByTransactionID(transactionID)
	if err != nil {
		return nil, ErrPaymentNotFound
	}

	if payment.IsPaid() {
		return payment, nil
	}

	now := time.Now()
	payment.Status = model.PaymentStatusPaid
	payment.PaidAt = &now

	if err := s.paymentRepo.UpdatePayment(payment); err != nil {
		return nil, err
	}

	// Update order
	s.orderRepo.UpdateOrderStatus(payment.OrderID, model.OrderStatusPaid)

	return payment, nil
}

func (s *paymentServiceEnhanced) CancelPayment(paymentID uint, reason string) (*model.Payment, error) {
	payment, err := s.paymentRepo.GetPaymentByID(paymentID)
	if err != nil {
		return nil, ErrPaymentNotFound
	}

	if payment.IsPaid() {
		return nil, errors.New("cannot cancel completed payment")
	}

	payment.Status = model.PaymentStatusCancelled
	payment.FailureReason = reason

	if err := s.paymentRepo.UpdatePayment(payment); err != nil {
		return nil, err
	}

	return payment, nil
}

// ==================== WEBHOOK HANDLING ====================

func (s *paymentServiceEnhanced) HandleWebhook(input *model.PaymentWebhookInput) error {
	// Verify signature
	if !s.VerifyWebhookSignature(input.TransactionID, input.Signature, "webhook_secret") {
		return ErrInvalidSignature
	}

	// Check for duplicate processing
	payment, err := s.paymentRepo.GetPaymentByTransactionID(input.TransactionID)
	if err == nil && payment != nil {
		// Payment already processed - idempotency
		return nil
	}

	// Process based on event type
	switch input.Event {
	case "payment.success", "payment.completed":
		if payment != nil {
			s.paymentRepo.MarkPaymentAsPaid(payment.ID, time.Now())
			s.orderRepo.UpdateOrderStatus(payment.OrderID, model.OrderStatusPaid)
		}
	case "payment.failed", "payment.declined":
		if payment != nil {
			s.paymentRepo.MarkPaymentAsFailed(payment.ID, "Gateway declined")
		}
	case "refund.completed":
		// Handle refund webhook
	}

	return nil
}

func (s *paymentServiceEnhanced) VerifyWebhookSignature(payload, signature, secret string) bool {
	// In production, verify HMAC signature from payment gateway
	// Example for Stripe:
	// expectedSignature := computeHMAC(payload, secret)
	// return hmac.Equal([]byte(signature), []byte(expectedSignature))

	// Mock verification
	return len(signature) > 0
}

// Helper function to compute HMAC
func computeHMAC(payload, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}

// ==================== REFUNDS ====================

func (s *paymentServiceEnhanced) RequestRefund(userID uint, input *model.RefundInput) (*model.Refund, error) {
	payment, err := s.paymentRepo.GetPaymentByID(input.PaymentID)
	if err != nil {
		return nil, ErrPaymentNotFound
	}

	// Validate refund amount
	refundAmount := input.Amount
	if refundAmount <= 0 {
		refundAmount = payment.GetRefundableAmount() // Full refund
	}

	if refundAmount > payment.GetRefundableAmount() {
		return nil, ErrRefundExceedsAmount
	}

	// Validate refund type
	if input.Type != "full" && input.Type != "partial" {
		return nil, ErrInvalidRefundType
	}

	refund := &model.Refund{
		PaymentID:   input.PaymentID,
		OrderID:     payment.OrderID,
		Amount:      refundAmount,
		Reason:      input.Reason,
		Status:      "pending",
		Type:        input.Type,
		RequestedBy: &userID,
	}

	if err := s.paymentRepo.CreateRefund(refund); err != nil {
		return nil, err
	}

	return refund, nil
}

func (s *paymentServiceEnhanced) ProcessRefund(refundID uint) (*model.Refund, error) {
	refund, err := s.paymentRepo.GetRefundByID(refundID)
	if err != nil {
		return nil, ErrRefundNotFound
	}

	// Call payment gateway to process refund
	gatewayResponse, err := s.RefundPayment(refund.Payment.PaymentProvider, refund.Payment.TransactionID, refund.Amount)
	if err != nil {
		refund.Status = "failed"
		s.paymentRepo.UpdateRefund(refund)
		return nil, err
	}

	// Update refund
	refund.Status = "processed"
	refund.GatewayRefundID = fmt.Sprintf("%v", gatewayResponse["id"])
	now := time.Now()
	refund.ProcessedAt = &now

	if err := s.paymentRepo.UpdateRefund(refund); err != nil {
		return nil, err
	}

	// Update payment refunded amount - simplified, in production use proper repository method
	return refund, nil
}

func (s *paymentServiceEnhanced) ApproveRefund(refundID uint, adminID uint) (*model.Refund, error) {
	refund, err := s.paymentRepo.GetRefundByID(refundID)
	if err != nil {
		return nil, ErrRefundNotFound
	}

	refund.Status = "approved"
	refund.ApprovedBy = &adminID
	now := time.Now()
	refund.ApprovedAt = &now

	if err := s.paymentRepo.UpdateRefund(refund); err != nil {
		return nil, err
	}

	// Process the refund
	return s.ProcessRefund(refundID)
}

func (s *paymentServiceEnhanced) RejectRefund(refundID uint, adminID uint, reason string) (*model.Refund, error) {
	refund, err := s.paymentRepo.GetRefundByID(refundID)
	if err != nil {
		return nil, ErrRefundNotFound
	}

	refund.Status = "rejected"
	refund.ApprovedBy = &adminID
	refund.Notes = reason
	now := time.Now()
	refund.ApprovedAt = &now

	if err := s.paymentRepo.UpdateRefund(refund); err != nil {
		return nil, err
	}

	return refund, nil
}

// ==================== PAYMENT METHODS ====================

func (s *paymentServiceEnhanced) SavePaymentMethod(userID uint, input *model.PaymentMethodInput) (*model.PaymentMethodModel, error) {
	method := &model.PaymentMethodModel{
		UserID:      userID,
		Type:        input.Type,
		Provider:    input.Provider,
		Name:        input.Name,
		LastFour:    input.LastFour,
		ExpiryMonth: input.ExpiryMonth,
		ExpiryYear:  input.ExpiryYear,
		Token:       input.Token,
		IsDefault:   input.IsDefault,
	}

	if err := s.paymentRepo.CreatePaymentMethod(method); err != nil {
		return nil, err
	}

	return method, nil
}

func (s *paymentServiceEnhanced) GetPaymentMethods(userID uint) ([]model.PaymentMethodModel, error) {
	return s.paymentRepo.GetPaymentMethodsByUser(userID)
}

func (s *paymentServiceEnhanced) DeletePaymentMethod(userID, methodID uint) error {
	method, err := s.paymentRepo.GetPaymentMethodByID(methodID)
	if err != nil {
		return ErrPaymentMethodNotFound
	}

	if method.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.paymentRepo.DeletePaymentMethod(methodID)
}

func (s *paymentServiceEnhanced) SetDefaultPaymentMethod(userID, methodID uint) error {
	return s.paymentRepo.SetDefaultPaymentMethod(userID, methodID)
}

// ==================== PAYMENT STATUS ====================

func (s *paymentServiceEnhanced) GetPaymentByOrderID(orderID uint) (*model.Payment, error) {
	return s.paymentRepo.GetPaymentByOrderID(orderID)
}

func (s *paymentServiceEnhanced) GetPaymentByTransactionID(transactionID string) (*model.Payment, error) {
	return s.paymentRepo.GetPaymentByTransactionID(transactionID)
}

func (s *paymentServiceEnhanced) GetUserPayments(userID uint, page, limit int) ([]model.Payment, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	return s.paymentRepo.GetPaymentsByUser(userID, limit, offset)
}

// ==================== ANALYTICS ====================

func (s *paymentServiceEnhanced) GetPaymentStats(userID uint) (*model.PaymentStats, error) {
	return s.paymentRepo.GetPaymentStats(userID)
}

// ==================== GATEWAY INTEGRATION (MOCK) ====================

func (s *paymentServiceEnhanced) InitializeGateway(provider string, config *model.PaymentGatewayConfig) error {
	s.gatewayConfigs[provider] = config
	return nil
}

func (s *paymentServiceEnhanced) ChargePayment(provider string, amount float64, token string) (map[string]interface{}, error) {
	// Mock gateway response - implement actual gateway integration
	return map[string]interface{}{
		"id":      "charge_xxx",
		"status":  "succeeded",
		"amount":  amount,
		"created": time.Now().Unix(),
	}, nil
}

func (s *paymentServiceEnhanced) RefundPayment(provider string, transactionID string, amount float64) (map[string]interface{}, error) {
	// Mock gateway response - implement actual gateway integration
	return map[string]interface{}{
		"id":              "refund_xxx",
		"status":          "succeeded",
		"amount":          amount,
		"transaction_id":  transactionID,
		"created":         time.Now().Unix(),
	}, nil
}

// ==================== CLEANUP ====================

func (s *paymentServiceEnhanced) CleanupExpiredPayments() (int64, error) {
	return s.paymentRepo.DeleteExpiredPendingPayments(30 * time.Minute)
}
