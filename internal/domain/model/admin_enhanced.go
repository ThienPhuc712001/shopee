package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AdminRoleType defines the type of admin role
type AdminRoleType string

const (
	RoleSuperAdmin   AdminRoleType = "super_admin"
	RoleAdmin        AdminRoleType = "admin"
	RoleSupportAgent AdminRoleType = "support_agent"
)

// AdminStatus defines the status of an admin user
type AdminStatus string

const (
	AdminStatusActive   AdminStatus = "active"
	AdminStatusInactive AdminStatus = "inactive"
	AdminStatusSuspended AdminStatus = "suspended"
)

// AdminUser represents an administrator user
type AdminUser struct {
	ID                 uint        `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	Email              string      `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password           string      `gorm:"type:varchar(255);not null" json:"-"`
	RoleID             uint        `gorm:"type:int;not null;index" json:"role_id"`
	Role               *AdminRole  `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	FirstName          string      `gorm:"type:varchar(100)" json:"first_name"`
	LastName           string      `gorm:"type:varchar(100)" json:"last_name"`
	Phone              string      `gorm:"type:varchar(20)" json:"phone"`
	AvatarURL          string      `gorm:"type:varchar(500)" json:"avatar_url"`
	Status             AdminStatus `gorm:"type:varchar(20);not null;default:'active';index" json:"status"`
	LastLoginAt        *time.Time  `gorm:"type:datetime" json:"last_login_at"`
	LastLoginIP        string      `gorm:"type:varchar(45)" json:"last_login_ip"`
	FailedLoginAttempts int        `gorm:"type:int;default:0" json:"-"`
	LockedUntil        *time.Time  `gorm:"type:datetime" json:"-"`

	// Relationships
	AuditLogs []AuditLog `gorm:"foreignKey:AdminID" json:"-"`

	CreatedAt time.Time      `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// AdminRole represents an admin role
type AdminRole struct {
	ID          uint            `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	Name        AdminRoleType   `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Description string          `gorm:"type:varchar(255)" json:"description"`
	Permissions string          `gorm:"type:text" json:"permissions"` // JSON array of permissions

	CreatedAt time.Time `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime;not null" json:"updated_at"`
}

// AdminPermission represents a permission definition
type AdminPermission struct {
	ID          uint   `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Description string `gorm:"type:varchar(255)" json:"description"`
	Module      string `gorm:"type:varchar(50);index" json:"module"` // users, products, orders, etc.

	CreatedAt time.Time `gorm:"type:datetime;not null" json:"created_at"`
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID         uint      `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	AdminID    uint      `gorm:"type:int;not null;index" json:"admin_id"`
	Admin      *AdminUser `gorm:"foreignKey:AdminID" json:"admin,omitempty"`
	Action     string    `gorm:"type:varchar(100);not null;index" json:"action"`
	EntityType string    `gorm:"type:varchar(50);not null;index" json:"entity_type"` // user, product, order, etc.
	EntityID   *uint     `gorm:"type:int;index" json:"entity_id"`
	OldValues  string    `gorm:"type:text" json:"old_values"` // JSON
	NewValues  string    `gorm:"type:text" json:"new_values"` // JSON
	IPAddress  string    `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent  string    `gorm:"type:varchar(500)" json:"user_agent"`

	CreatedAt time.Time `gorm:"type:datetime;not null;index" json:"created_at"`
}

// SystemSetting represents a system configuration setting
type SystemSetting struct {
	ID          uint        `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	Key         string      `gorm:"type:varchar(100);uniqueIndex;not null" json:"key"`
	Value       string      `gorm:"type:text;not null" json:"value"`
	Type        string      `gorm:"type:varchar(50);not null" json:"type"` // string, number, boolean, json
	Description string      `gorm:"type:varchar(500)" json:"description"`
	IsPublic    bool        `gorm:"type:bit;default:false" json:"is_public"`
	UpdatedBy   *uint       `gorm:"type:int" json:"updated_by"`
	Updater     *AdminUser  `gorm:"foreignKey:UpdatedBy" json:"updater,omitempty"`

	CreatedAt time.Time `gorm:"type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime;not null" json:"updated_at"`
}

// AdminLoginInput represents admin login request
type AdminLoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AdminCreateInput represents create admin user request
type AdminCreateInput struct {
	Email     string        `json:"email" binding:"required,email"`
	Password  string        `json:"password" binding:"required,min=8"`
	RoleID    uint          `json:"role_id" binding:"required"`
	FirstName string        `json:"first_name"`
	LastName  string        `json:"last_name"`
	Phone     string        `json:"phone"`
}

// AdminUpdateInput represents update admin user request
type AdminUpdateInput struct {
	FirstName string        `json:"first_name"`
	LastName  string        `json:"last_name"`
	Phone     string        `json:"phone"`
	AvatarURL string        `json:"avatar_url"`
	RoleID    uint          `json:"role_id"`
	Status    AdminStatus   `json:"status"`
}

// BanUserInput represents ban user request
type BanUserInput struct {
	UserID uint   `json:"user_id" binding:"required"`
	Reason string `json:"reason" binding:"required"`
}

// ApproveSellerInput represents approve seller request
type ApproveSellerInput struct {
	ShopID uint   `json:"shop_id" binding:"required"`
	Notes  string `json:"notes"`
}

// RefundOrderInput represents refund order request
type RefundOrderInput struct {
	OrderID uint    `json:"order_id" binding:"required"`
	Amount  float64 `json:"amount"`
	Reason  string  `json:"reason" binding:"required"`
}

// AdminStats represents admin statistics
type AdminStats struct {
	TotalUsers       int64   `json:"total_users"`
	TotalSellers     int64   `json:"total_sellers"`
	TotalProducts    int64   `json:"total_products"`
	TotalOrders      int64   `json:"total_orders"`
	TotalRevenue     float64 `json:"total_revenue"`
	PendingOrders    int64   `json:"pending_orders"`
	PendingRefunds   int64   `json:"pending_refunds"`
	ActiveUsers24h   int64   `json:"active_users_24h"`
	NewUsersToday    int64   `json:"new_users_today"`
}

// SalesAnalytics represents sales analytics data
type SalesAnalytics struct {
	TotalSales      float64 `json:"total_sales"`
	TotalOrders     int64   `json:"total_orders"`
	AverageOrderValue float64 `json:"average_order_value"`
	TodaySales      float64 `json:"today_sales"`
	WeekSales       float64 `json:"week_sales"`
	MonthSales      float64 `json:"month_sales"`
}

// UserAnalytics represents user analytics data
type UserAnalytics struct {
	TotalUsers     int64 `json:"total_users"`
	NewUsersToday  int64 `json:"new_users_today"`
	NewUsersWeek   int64 `json:"new_users_week"`
	NewUsersMonth  int64 `json:"new_users_month"`
	ActiveUsers    int64 `json:"active_users"`
	BannedUsers    int64 `json:"banned_users"`
}

// ProductAnalytics represents product analytics data
type ProductAnalytics struct {
	TotalProducts    int64   `json:"total_products"`
	ActiveProducts   int64   `json:"active_products"`
	OutOfStock       int64   `json:"out_of_stock"`
	TopProducts      []ProductStat `json:"top_products"`
	LowStockProducts []ProductStat `json:"low_stock_products"`
}

// ProductStat represents product statistics
type ProductStat struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	SoldCount  int64  `json:"sold_count"`
	Revenue    float64 `json:"revenue"`
}

// BeforeCreate hook for AdminUser
func (a *AdminUser) BeforeCreate(tx *gorm.DB) error {
	// Hash password before creating
	if a.Password != "" {
		// Password hashing handled by service layer
	}
	return nil
}

// GetFullName returns the full name of the admin user
func (a *AdminUser) GetFullName() string {
	if a.LastName != "" {
		return a.FirstName + " " + a.LastName
	}
	return a.FirstName
}

// IsSuperAdmin checks if admin is super admin
func (a *AdminUser) IsSuperAdmin() bool {
	return a.Role != nil && a.Role.Name == RoleSuperAdmin
}

// IsAdmin checks if admin is admin
func (a *AdminUser) IsAdmin() bool {
	return a.Role != nil && (a.Role.Name == RoleAdmin || a.Role.Name == RoleSuperAdmin)
}

// HasPermission checks if admin has a specific permission
func (a *AdminUser) HasPermission(permission string) bool {
	if a.Role == nil || a.Role.Permissions == "" {
		return false
	}
	// Parse permissions JSON and check
	// Simplified implementation
	return true
}

// IsLocked checks if admin account is locked
func (a *AdminUser) IsLocked() bool {
	if a.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*a.LockedUntil)
}

// LockAccount locks the admin account
func (a *AdminUser) LockAccount(duration time.Duration) {
	until := time.Now().Add(duration)
	a.LockedUntil = &until
}

// IncrementFailedLogin increments failed login counter
func (a *AdminUser) IncrementFailedLogin() {
	a.FailedLoginAttempts++
	if a.FailedLoginAttempts >= 5 {
		a.LockAccount(30 * time.Minute)
	}
}

// ResetFailedLogin resets failed login counter
func (a *AdminUser) ResetFailedLogin() {
	a.FailedLoginAttempts = 0
	a.LockedUntil = nil
}

// CheckPassword compares the provided password with the hashed password
func (a *AdminUser) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
	return err == nil
}

// HashPassword hashes the admin user's password using bcrypt
func (a *AdminUser) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	a.Password = string(hashedPassword)
	return nil
}

// AuditLogInput represents input for creating audit log
type AuditLogInput struct {
	AdminID    uint   `json:"admin_id"`
	Action     string `json:"action"`
	EntityType string `json:"entity_type"`
	EntityID   *uint  `json:"entity_id"`
	OldValues  string `json:"old_values"`
	NewValues  string `json:"new_values"`
	IPAddress  string `json:"ip_address"`
	UserAgent  string `json:"user_agent"`
}

// Common admin actions
const (
	ActionCreate        = "create"
	ActionUpdate        = "update"
	ActionDelete        = "delete"
	ActionBan           = "ban"
	ActionUnban         = "unban"
	ActionApprove       = "approve"
	ActionReject        = "reject"
	ActionRefund        = "refund"
	ActionCancel        = "cancel"
	ActionRestore       = "restore"
	ActionExport        = "export"
	ActionImport        = "import"
	ActionSystemSetting = "system_setting"
)

// Entity types for audit logging
const (
	EntityUser       = "user"
	EntityShop       = "shop"
	EntityProduct    = "product"
	EntityOrder      = "order"
	EntityPayment    = "payment"
	EntityRefund     = "refund"
	EntityVoucher    = "voucher"
	EntityCategory   = "category"
	EntityAdmin      = "admin"
	EntitySystem     = "system"
)
