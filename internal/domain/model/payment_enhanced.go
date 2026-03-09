package model

import (
	"time"

	"gorm.io/gorm"
)

// PaymentStatus defines the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending       PaymentStatus = "pending"
	PaymentStatusProcessing    PaymentStatus = "processing"
	PaymentStatusPaid          PaymentStatus = "paid"
	PaymentStatusFailed        PaymentStatus = "failed"
	PaymentStatusCancelled     PaymentStatus = "cancelled"
	PaymentStatusRefunded      PaymentStatus = "refunded"
	PaymentStatusPartialRefund PaymentStatus = "partial_refund"
)

// PaymentType defines the type of payment transaction
type PaymentType string

const (
	PaymentTypeCharge   PaymentType = "charge"
	PaymentTypeRefund   PaymentType = "refund"
	PaymentTypeChargeback PaymentType = "chargeback"
)

// PaymentMethod defines the payment method
type PaymentMethod string

const (
	PaymentMethodCreditCard  PaymentMethod = "credit_card"
	PaymentMethodDebitCard   PaymentMethod = "debit_card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodEWallet     PaymentMethod = "e_wallet"
	PaymentMethodCOD         PaymentMethod = "cod"
	PaymentMethodPayPal      PaymentMethod = "paypal"
)

// Payment represents a payment transaction
type Payment struct {
	ID               uint          `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	OrderID          uint          `gorm:"type:int;not null;index" json:"order_id"`
	Order            *Order        `gorm:"foreignKey:OrderID" json:"-"`
	UserID           uint          `gorm:"type:int;not null;index" json:"user_id"`
	User             *User         `gorm:"foreignKey:UserID" json:"-"`
	TransactionID    string        `gorm:"type:varchar(255);uniqueIndex" json:"transaction_id"`
	PaymentMethod    PaymentMethod `gorm:"type:varchar(50);not null" json:"payment_method"`
	PaymentProvider  string        `gorm:"type:varchar(50)" json:"payment_provider"` // stripe, paypal, vnpay, etc.
	Amount           float64       `gorm:"type:decimal(18,2);not null" json:"amount"`
	Currency         string        `gorm:"type:varchar(3);default:'USD'" json:"currency"`
	Status           PaymentStatus `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`
	GatewayResponse  string        `gorm:"type:text" json:"gateway_response"` // JSON
	Metadata         string        `gorm:"type:text" json:"metadata"`         // JSON
	PaidAt           *time.Time    `gorm:"type:datetime" json:"paid_at"`
	FailedAt         *time.Time    `gorm:"type:datetime" json:"failed_at"`
	FailureReason    string        `gorm:"type:varchar(500)" json:"failure_reason"`
	RefundedAmount   float64       `gorm:"type:decimal(18,2);default:0" json:"refunded_amount"`

	// Relationships
	Transactions []PaymentTransaction `gorm:"foreignKey:PaymentID" json:"transactions,omitempty"`
	Refunds      []Refund             `gorm:"foreignKey:PaymentID" json:"refunds,omitempty"`

	CreatedAt time.Time      `gorm:"type:datetime;not null;index" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// PaymentMethodModel represents a saved payment method for a user
type PaymentMethodModel struct {
	ID          uint      `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	UserID      uint      `gorm:"type:int;index" json:"user_id"`
	User        *User     `gorm:"foreignKey:UserID" json:"-"`
	Type        string    `gorm:"type:varchar(50);not null" json:"type"` // credit_card, bank_account, e_wallet
	Provider    string    `gorm:"type:varchar(50);not null" json:"provider"` // stripe, paypal, etc.
	Name        string    `gorm:"type:varchar(100)" json:"name"` // e.g., "Visa ending 4242"
	LastFour    string    `gorm:"type:varchar(4)" json:"last_four"`
	ExpiryMonth int       `gorm:"type:int" json:"expiry_month"`
	ExpiryYear  int       `gorm:"type:int" json:"expiry_year"`
	IsDefault   bool      `gorm:"type:bit;default:false" json:"is_default"`
	Token       string    `gorm:"type:varchar(500)" json:"-"` // Payment provider token (never expose)
	Metadata    string    `gorm:"type:text" json:"metadata"` // JSON

	CreatedAt time.Time      `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// PaymentTransaction represents a payment transaction record
type PaymentTransaction struct {
	ID             uint        `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	PaymentID      uint        `gorm:"type:int;not null;index" json:"payment_id"`
	Payment        *Payment    `gorm:"foreignKey:PaymentID" json:"-"`
	Type           PaymentType `gorm:"type:varchar(50);not null" json:"type"` // charge, refund, chargeback
	Amount         float64     `gorm:"type:decimal(18,2);not null" json:"amount"`
	Status         string      `gorm:"type:varchar(20);not null" json:"status"`
	GatewayID      string      `gorm:"type:varchar(255)" json:"gateway_id"`
	GatewayResponse string     `gorm:"type:text" json:"gateway_response"` // JSON
	ProcessedAt    *time.Time  `gorm:"type:datetime" json:"processed_at"`

	CreatedAt time.Time `gorm:"type:datetime;not null;index" json:"created_at"`
}

// Refund represents a refund request
type Refund struct {
	ID              uint        `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	PaymentID       uint        `gorm:"type:int;not null;index" json:"payment_id"`
	Payment         *Payment    `gorm:"foreignKey:PaymentID" json:"-"`
	OrderID         uint        `gorm:"type:int;not null;index" json:"order_id"`
	Order           *Order      `gorm:"foreignKey:OrderID" json:"-"`
	RefundNumber    string      `gorm:"type:varchar(50);uniqueIndex" json:"refund_number"`
	Amount          float64     `gorm:"type:decimal(18,2);not null" json:"amount"`
	Reason          string      `gorm:"type:varchar(500);not null" json:"reason"`
	Status          string      `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`
	Type            string      `gorm:"type:varchar(20);not null" json:"type"` // full, partial
	RequestedBy     *uint       `gorm:"type:int" json:"requested_by"`
	RequestedByUser *User       `gorm:"foreignKey:RequestedBy" json:"-"`
	ApprovedBy      *uint       `gorm:"type:int" json:"approved_by"`
	ApprovedByUser  *User       `gorm:"foreignKey:ApprovedBy" json:"-"`
	ApprovedAt      *time.Time  `gorm:"type:datetime" json:"approved_at"`
	ProcessedAt     *time.Time  `gorm:"type:datetime" json:"processed_at"`
	GatewayRefundID string      `gorm:"type:varchar(255)" json:"gateway_refund_id"`
	Notes           string      `gorm:"type:varchar(500)" json:"notes"`

	CreatedAt time.Time `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime;not null" json:"updated_at"`
}

// PaymentInput represents input for creating a payment
type PaymentInput struct {
	OrderID       uint          `json:"order_id" binding:"required"`
	PaymentMethod PaymentMethod `json:"payment_method" binding:"required"`
	Provider      string        `json:"provider"`
	SaveMethod    bool          `json:"save_method"` // Save payment method for future use
}

// PaymentConfirmInput represents input for confirming a payment
type PaymentConfirmInput struct {
	TransactionID string `json:"transaction_id" binding:"required"`
	PaymentID     uint   `json:"payment_id"`
}

// PaymentWebhookInput represents webhook payload from payment gateway
type PaymentWebhookInput struct {
	Event         string            `json:"event"`
	TransactionID string            `json:"transaction_id"`
	PaymentID     string            `json:"payment_id"`
	Status        string            `json:"status"`
	Amount        float64           `json:"amount"`
	Currency      string            `json:"currency"`
	Signature     string            `json:"signature"`
	Metadata      map[string]string `json:"metadata"`
	Timestamp     int64             `json:"timestamp"`
}

// RefundInput represents input for requesting a refund
type RefundInput struct {
	PaymentID uint   `json:"payment_id" binding:"required"`
	Amount    float64 `json:"amount"` // Optional for partial refunds
	Reason    string  `json:"reason" binding:"required"`
	Type      string  `json:"type" binding:"required"` // full or partial
}

// PaymentResponse represents payment response
type PaymentResponse struct {
	PaymentID     uint        `json:"payment_id"`
	TransactionID string      `json:"transaction_id"`
	OrderID       uint        `json:"order_id"`
	Amount        float64     `json:"amount"`
	Currency      string      `json:"currency"`
	Status        PaymentStatus `json:"status"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	PaymentURL    string      `json:"payment_url,omitempty"` // For redirect payments
	ClientSecret  string      `json:"client_secret,omitempty"` // For Stripe
}

// PaymentGatewayConfig represents payment gateway configuration
type PaymentGatewayConfig struct {
	Provider     string `json:"provider"`
	APIKey       string `json:"api_key"`
	SecretKey    string `json:"secret_key"`
	WebhookSecret string `json:"webhook_secret"`
	Environment  string `json:"environment"` // test, production
	MerchantID   string `json:"merchant_id"`
}

// BeforeCreate hook for Payment - generates transaction ID
func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.TransactionID == "" {
		p.TransactionID = GenerateTransactionID()
	}
	return nil
}

// BeforeCreate hook for Refund - generates refund number
func (r *Refund) BeforeCreate(tx *gorm.DB) error {
	if r.RefundNumber == "" {
		r.RefundNumber = GenerateRefundNumber()
	}
	return nil
}

// IsPaid checks if payment is completed
func (p *Payment) IsPaid() bool {
	return p.Status == PaymentStatusPaid
}

// IsPending checks if payment is pending
func (p *Payment) IsPending() bool {
	return p.Status == PaymentStatusPending || p.Status == PaymentStatusProcessing
}

// IsFailed checks if payment failed
func (p *Payment) IsFailed() bool {
	return p.Status == PaymentStatusFailed
}

// IsRefunded checks if payment is fully refunded
func (p *Payment) IsRefunded() bool {
	return p.Status == PaymentStatusRefunded
}

// CanRefund checks if payment can be refunded
func (p *Payment) CanRefund() bool {
	return p.Status == PaymentStatusPaid && p.RefundedAmount < p.Amount
}

// GetRefundableAmount returns the amount that can be refunded
func (p *Payment) GetRefundableAmount() float64 {
	return p.Amount - p.RefundedAmount
}

// GenerateTransactionID generates a unique transaction ID
func GenerateTransactionID() string {
	return "TXN-" + time.Now().Format("20060102150405") + "-" + generateRandomString(8)
}

// GenerateRefundNumber generates a unique refund number
func GenerateRefundNumber() string {
	return "REF-" + time.Now().Format("20060102") + "-" + generateRandomString(6)
}

// Helper function for generating random strings
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(time.Nanosecond)
	}
	return string(b)
}

// PaymentStats represents payment statistics
type PaymentStats struct {
	TotalPayments     int64   `json:"total_payments"`
	TotalAmount       float64 `json:"total_amount"`
	SuccessfulPayments int64  `json:"successful_payments"`
	FailedPayments    int64   `json:"failed_payments"`
	PendingPayments   int64   `json:"pending_payments"`
	TotalRefunds      int64   `json:"total_refunds"`
	RefundAmount      float64 `json:"refund_amount"`
}

// PaymentMethodInput represents input for saving a payment method
type PaymentMethodInput struct {
	Type        string `json:"type" binding:"required"`
	Provider    string `json:"provider" binding:"required"`
	Name        string `json:"name"`
	LastFour    string `json:"last_four"`
	ExpiryMonth int    `json:"expiry_month"`
	ExpiryYear  int    `json:"expiry_year"`
	Token       string `json:"token" binding:"required"`
	IsDefault   bool   `json:"is_default"`
}
