package repository

import (
	"ecommerce/internal/domain/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

// PaymentRepositoryEnhanced defines the enhanced interface for payment data operations
type PaymentRepositoryEnhanced interface {
	// Payment CRUD
	CreatePayment(payment *model.Payment) error
	GetPaymentByID(id uint) (*model.Payment, error)
	GetPaymentByOrderID(orderID uint) (*model.Payment, error)
	GetPaymentByTransactionID(transactionID string) (*model.Payment, error)
	UpdatePayment(payment *model.Payment) error

	// Payment Queries
	GetPaymentsByUser(userID uint, limit, offset int) ([]model.Payment, int64, error)
	GetPaymentsByStatus(status model.PaymentStatus, limit, offset int) ([]model.Payment, int64, error)
	GetPendingPayments() ([]model.Payment, error)

	// Payment Status
	UpdatePaymentStatus(paymentID uint, status model.PaymentStatus) error
	MarkPaymentAsPaid(paymentID uint, paidAt time.Time) error
	MarkPaymentAsFailed(paymentID uint, reason string) error

	// Payment Transactions
	CreateTransaction(transaction *model.PaymentTransaction) error
	GetTransactionsByPaymentID(paymentID uint) ([]model.PaymentTransaction, error)

	// Refunds
	CreateRefund(refund *model.Refund) error
	GetRefundByID(id uint) (*model.Refund, error)
	GetRefundByNumber(refundNumber string) (*model.Refund, error)
	GetRefundsByPaymentID(paymentID uint) ([]model.Refund, error)
	GetRefundsByOrderID(orderID uint) ([]model.Refund, error)
	UpdateRefund(refund *model.Refund) error
	UpdateRefundStatus(refundID uint, status string) error

	// Payment Methods
	CreatePaymentMethod(method *model.PaymentMethodModel) error
	GetPaymentMethodsByUser(userID uint) ([]model.PaymentMethodModel, error)
	GetPaymentMethodByID(id uint) (*model.PaymentMethodModel, error)
	UpdatePaymentMethod(method *model.PaymentMethodModel) error
	DeletePaymentMethod(id uint) error
	SetDefaultPaymentMethod(userID, methodID uint) error

	// Analytics
	GetPaymentStats(userID uint) (*model.PaymentStats, error)
	GetRevenueByDateRange(startDate, endDate time.Time) (float64, error)
	GetPaymentCountByStatus(status model.PaymentStatus) (int64, error)

	// Cleanup
	DeleteExpiredPendingPayments(olderThan time.Duration) (int64, error)
}

type paymentRepositoryEnhanced struct {
	db *gorm.DB
}

// NewPaymentRepositoryEnhanced creates a new enhanced payment repository
func NewPaymentRepositoryEnhanced(db *gorm.DB) PaymentRepositoryEnhanced {
	return &paymentRepositoryEnhanced{db: db}
}

// ==================== PAYMENT CRUD ====================

func (r *paymentRepositoryEnhanced) CreatePayment(payment *model.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepositoryEnhanced) GetPaymentByID(id uint) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.Preload("Order").
		Preload("Transactions").
		Preload("Refunds").
		First(&payment, id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepositoryEnhanced) GetPaymentByOrderID(orderID uint) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.Where("order_id = ?", orderID).
		Preload("Transactions").
		Preload("Refunds").
		First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepositoryEnhanced) GetPaymentByTransactionID(transactionID string) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.Where("transaction_id = ?", transactionID).
		Preload("Order").
		First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepositoryEnhanced) UpdatePayment(payment *model.Payment) error {
	return r.db.Save(payment).Error
}

// ==================== PAYMENT QUERIES ====================

func (r *paymentRepositoryEnhanced) GetPaymentsByUser(userID uint, limit, offset int) ([]model.Payment, int64, error) {
	var payments []model.Payment
	var total int64

	if err := r.db.Model(&model.Payment{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where("user_id = ?", userID).
		Preload("Order").
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&payments).Error

	return payments, total, err
}

func (r *paymentRepositoryEnhanced) GetPaymentsByStatus(status model.PaymentStatus, limit, offset int) ([]model.Payment, int64, error) {
	var payments []model.Payment
	var total int64

	if err := r.db.Model(&model.Payment{}).Where("status = ?", status).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where("status = ?", status).
		Preload("Order").
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&payments).Error

	return payments, total, err
}

func (r *paymentRepositoryEnhanced) GetPendingPayments() ([]model.Payment, error) {
	var payments []model.Payment
	err := r.db.Where("status = ?", model.PaymentStatusPending).
		Where("created_at < ?", time.Now().Add(-30*time.Minute)).
		Find(&payments).Error
	return payments, err
}

// ==================== PAYMENT STATUS ====================

func (r *paymentRepositoryEnhanced) UpdatePaymentStatus(paymentID uint, status model.PaymentStatus) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	switch status {
	case model.PaymentStatusPaid:
		updates["paid_at"] = time.Now()
	case model.PaymentStatusFailed:
		updates["failed_at"] = time.Now()
	}

	return r.db.Model(&model.Payment{}).Where("id = ?", paymentID).Updates(updates).Error
}

func (r *paymentRepositoryEnhanced) MarkPaymentAsPaid(paymentID uint, paidAt time.Time) error {
	return r.db.Model(&model.Payment{}).Where("id = ?", paymentID).Updates(map[string]interface{}{
		"status":     model.PaymentStatusPaid,
		"paid_at":    paidAt,
		"updated_at": time.Now(),
	}).Error
}

func (r *paymentRepositoryEnhanced) MarkPaymentAsFailed(paymentID uint, reason string) error {
	return r.db.Model(&model.Payment{}).Where("id = ?", paymentID).Updates(map[string]interface{}{
		"status":          model.PaymentStatusFailed,
		"failure_reason":  reason,
		"failed_at":       time.Now(),
		"updated_at":      time.Now(),
	}).Error
}

// ==================== PAYMENT TRANSACTIONS ====================

func (r *paymentRepositoryEnhanced) CreateTransaction(transaction *model.PaymentTransaction) error {
	return r.db.Create(transaction).Error
}

func (r *paymentRepositoryEnhanced) GetTransactionsByPaymentID(paymentID uint) ([]model.PaymentTransaction, error) {
	var transactions []model.PaymentTransaction
	err := r.db.Where("payment_id = ?", paymentID).
		Order("created_at DESC").
		Find(&transactions).Error
	return transactions, err
}

// ==================== REFUNDS ====================

func (r *paymentRepositoryEnhanced) CreateRefund(refund *model.Refund) error {
	return r.db.Create(refund).Error
}

func (r *paymentRepositoryEnhanced) GetRefundByID(id uint) (*model.Refund, error) {
	var refund model.Refund
	err := r.db.Preload("Payment").Preload("Order").First(&refund, id).Error
	if err != nil {
		return nil, err
	}
	return &refund, nil
}

func (r *paymentRepositoryEnhanced) GetRefundByNumber(refundNumber string) (*model.Refund, error) {
	var refund model.Refund
	err := r.db.Where("refund_number = ?", refundNumber).
		Preload("Payment").Preload("Order").
		First(&refund).Error
	if err != nil {
		return nil, err
	}
	return &refund, nil
}

func (r *paymentRepositoryEnhanced) GetRefundsByPaymentID(paymentID uint) ([]model.Refund, error) {
	var refunds []model.Refund
	err := r.db.Where("payment_id = ?", paymentID).
		Order("created_at DESC").
		Find(&refunds).Error
	return refunds, err
}

func (r *paymentRepositoryEnhanced) GetRefundsByOrderID(orderID uint) ([]model.Refund, error) {
	var refunds []model.Refund
	err := r.db.Where("order_id = ?", orderID).
		Preload("Payment").
		Order("created_at DESC").
		Find(&refunds).Error
	return refunds, err
}

func (r *paymentRepositoryEnhanced) UpdateRefund(refund *model.Refund) error {
	return r.db.Save(refund).Error
}

func (r *paymentRepositoryEnhanced) UpdateRefundStatus(refundID uint, status string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	if status == "processed" {
		updates["processed_at"] = time.Now()
	}

	return r.db.Model(&model.Refund{}).Where("id = ?", refundID).Updates(updates).Error
}

// ==================== PAYMENT METHODS ====================

func (r *paymentRepositoryEnhanced) CreatePaymentMethod(method *model.PaymentMethodModel) error {
	// If this is set as default, unset other defaults
	if method.IsDefault {
		r.db.Model(&model.PaymentMethodModel{}).
			Where("user_id = ?", method.UserID).
			Update("is_default", false)
	}

	return r.db.Create(method).Error
}

func (r *paymentRepositoryEnhanced) GetPaymentMethodsByUser(userID uint) ([]model.PaymentMethodModel, error) {
	var methods []model.PaymentMethodModel
	err := r.db.Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("is_default DESC, created_at DESC").
		Find(&methods).Error
	return methods, err
}

func (r *paymentRepositoryEnhanced) GetPaymentMethodByID(id uint) (*model.PaymentMethodModel, error) {
	var method model.PaymentMethodModel
	err := r.db.First(&method, id).Error
	if err != nil {
		return nil, err
	}
	return &method, nil
}

func (r *paymentRepositoryEnhanced) UpdatePaymentMethod(method *model.PaymentMethodModel) error {
	return r.db.Save(method).Error
}

func (r *paymentRepositoryEnhanced) DeletePaymentMethod(id uint) error {
	return r.db.Delete(&model.PaymentMethodModel{}, id).Error
}

func (r *paymentRepositoryEnhanced) SetDefaultPaymentMethod(userID, methodID uint) error {
	tx := r.db.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// Unset all defaults for user
	if err := tx.Model(&model.PaymentMethodModel{}).
		Where("user_id = ?", userID).
		Update("is_default", false).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Set new default
	if err := tx.Model(&model.PaymentMethodModel{}).
		Where("id = ? AND user_id = ?", methodID, userID).
		Update("is_default", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// ==================== ANALYTICS ====================

func (r *paymentRepositoryEnhanced) GetPaymentStats(userID uint) (*model.PaymentStats, error) {
	stats := &model.PaymentStats{}

	// Total payments
	r.db.Model(&model.Payment{}).
		Where("user_id = ?", userID).
		Count(&stats.TotalPayments)

	// Total amount
	r.db.Model(&model.Payment{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&stats.TotalAmount)

	// Successful payments
	r.db.Model(&model.Payment{}).
		Where("user_id = ? AND status = ?", userID, model.PaymentStatusPaid).
		Count(&stats.SuccessfulPayments)

	// Failed payments
	r.db.Model(&model.Payment{}).
		Where("user_id = ? AND status = ?", userID, model.PaymentStatusFailed).
		Count(&stats.FailedPayments)

	// Pending payments
	r.db.Model(&model.Payment{}).
		Where("user_id = ? AND status IN ?", userID, []model.PaymentStatus{model.PaymentStatusPending, model.PaymentStatusProcessing}).
		Count(&stats.PendingPayments)

	// Total refunds
	r.db.Model(&model.Refund{}).
		Joins("JOIN payments ON refunds.payment_id = payments.id").
		Where("payments.user_id = ?", userID).
		Count(&stats.TotalRefunds)

	// Refund amount
	r.db.Model(&model.Refund{}).
		Joins("JOIN payments ON refunds.payment_id = payments.id").
		Where("payments.user_id = ?", userID).
		Select("COALESCE(SUM(refunds.amount), 0)").
		Scan(&stats.RefundAmount)

	return stats, nil
}

func (r *paymentRepositoryEnhanced) GetRevenueByDateRange(startDate, endDate time.Time) (float64, error) {
	var revenue float64
	err := r.db.Model(&model.Payment{}).
		Where("status = ?", model.PaymentStatusPaid).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&revenue).Error
	return revenue, err
}

func (r *paymentRepositoryEnhanced) GetPaymentCountByStatus(status model.PaymentStatus) (int64, error) {
	var count int64
	err := r.db.Model(&model.Payment{}).
		Where("status = ?", status).
		Count(&count).Error
	return count, err
}

// ==================== CLEANUP ====================

func (r *paymentRepositoryEnhanced) DeleteExpiredPendingPayments(olderThan time.Duration) (int64, error) {
	cutoffTime := time.Now().Add(-olderThan)

	result := r.db.Where("status = ?", model.PaymentStatusPending).
		Where("created_at < ?", cutoffTime).
		Delete(&model.Payment{})

	return result.RowsAffected, result.Error
}

// ==================== TRANSACTIONAL OPERATIONS ====================

// CreatePaymentWithTransaction creates payment and transaction record in a transaction
func (r *paymentRepositoryEnhanced) CreatePaymentWithTransaction(payment *model.Payment, transaction *model.PaymentTransaction) error {
	tx := r.db.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// Create payment
	if err := tx.Create(payment).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Create transaction
	transaction.PaymentID = payment.ID
	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// ProcessRefund processes a refund in a transaction
func (r *paymentRepositoryEnhanced) ProcessRefund(refund *model.Refund, updatePaymentAmount float64) error {
	tx := r.db.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// Create refund
	if err := tx.Create(refund).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update payment refunded amount
	if err := tx.Model(&model.Payment{}).Where("id = ?", refund.PaymentID).
		Update("refunded_amount", gorm.Expr("refunded_amount + ?", updatePaymentAmount)).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update payment status if fully refunded
	var payment model.Payment
	if err := tx.First(&payment, refund.PaymentID).Error; err != nil {
		tx.Rollback()
		return err
	}

	if payment.RefundedAmount+updatePaymentAmount >= payment.Amount {
		if err := tx.Model(&model.Payment{}).Where("id = ?", refund.PaymentID).
			Update("status", model.PaymentStatusRefunded).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// Error definitions
var (
	ErrPaymentNotFound      = errors.New("payment not found")
	ErrRefundNotFound       = errors.New("refund not found")
	ErrPaymentMethodNotFound = errors.New("payment method not found")
	ErrDuplicateTransaction = errors.New("duplicate transaction ID")
	ErrInvalidRefundAmount  = errors.New("invalid refund amount")
	ErrRefundExceedsPayment = errors.New("refund amount exceeds payment")
)
