package repository

import (
	"ecommerce/internal/domain/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

// AdminRepositoryEnhanced defines the enhanced interface for admin data operations
type AdminRepositoryEnhanced interface {
	// Admin User CRUD
	CreateAdminUser(admin *model.AdminUser) error
	GetAdminUserByID(id uint) (*model.AdminUser, error)
	GetAdminUserByEmail(email string) (*model.AdminUser, error)
	UpdateAdminUser(admin *model.AdminUser) error
	DeleteAdminUser(id uint) error

	// Admin User Queries
	GetAdminUsers(limit, offset int) ([]model.AdminUser, int64, error)
	GetAdminUsersByRole(roleID uint, limit, offset int) ([]model.AdminUser, int64, error)
	GetAdminUsersByStatus(status model.AdminStatus, limit, offset int) ([]model.AdminUser, int64, error)

	// Admin User Status
	UpdateAdminStatus(id uint, status model.AdminStatus) error
	UpdateAdminLastLogin(id uint, ip string) error
	UpdateAdminPassword(id uint, hashedPassword string) error

	// Admin Roles
	GetAdminRoleByID(id uint) (*model.AdminRole, error)
	GetAdminRoleByName(name model.AdminRoleType) (*model.AdminRole, error)
	GetAllAdminRoles() ([]model.AdminRole, error)
	CreateAdminRole(role *model.AdminRole) error
	UpdateAdminRole(role *model.AdminRole) error

	// Audit Logs
	CreateAuditLog(log *model.AuditLog) error
	GetAuditLogs(limit, offset int) ([]model.AuditLog, int64, error)
	GetAuditLogsByAdminID(adminID uint, limit, offset int) ([]model.AuditLog, int64, error)
	GetAuditLogsByAction(action string, limit, offset int) ([]model.AuditLog, int64, error)
	GetAuditLogsByEntityType(entityType string, limit, offset int) ([]model.AuditLog, int64, error)
	GetAuditLogsByDateRange(startDate, endDate time.Time, limit, offset int) ([]model.AuditLog, int64, error)

	// System Settings
	GetSystemSetting(key string) (*model.SystemSetting, error)
	GetAllSystemSettings() ([]model.SystemSetting, error)
	GetPublicSystemSettings() ([]model.SystemSetting, error)
	UpdateSystemSetting(key string, value string, updatedBy *uint) error
	CreateSystemSetting(setting *model.SystemSetting) error

	// Review Management
	GetAllReviews(limit, offset int) ([]model.Review, int64, error)
	DeleteSystemSetting(key string) error

	// Refunds
	CreateRefund(refund *model.Refund) error

	// Analytics
	GetAdminStats() (*model.AdminStats, error)
	GetSalesAnalytics(startDate, endDate time.Time) (*model.SalesAnalytics, error)
	GetUserAnalytics() (*model.UserAnalytics, error)
	GetProductAnalytics(limit int) (*model.ProductAnalytics, error)

	// Cleanup
	DeleteOldAuditLogs(olderThan time.Duration) (int64, error)
}

type adminRepositoryEnhanced struct {
	db *gorm.DB
}

// NewAdminRepositoryEnhanced creates a new enhanced admin repository
func NewAdminRepositoryEnhanced(db *gorm.DB) AdminRepositoryEnhanced {
	return &adminRepositoryEnhanced{db: db}
}

// ==================== ADMIN USER CRUD ====================

func (r *adminRepositoryEnhanced) CreateAdminUser(admin *model.AdminUser) error {
	return r.db.Create(admin).Error
}

func (r *adminRepositoryEnhanced) GetAdminUserByID(id uint) (*model.AdminUser, error) {
	var admin model.AdminUser
	err := r.db.Preload("Role").First(&admin, id).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminRepositoryEnhanced) GetAdminUserByEmail(email string) (*model.AdminUser, error) {
	var admin model.AdminUser
	err := r.db.Preload("Role").Where("email = ?", email).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminRepositoryEnhanced) UpdateAdminUser(admin *model.AdminUser) error {
	return r.db.Save(admin).Error
}

func (r *adminRepositoryEnhanced) DeleteAdminUser(id uint) error {
	return r.db.Delete(&model.AdminUser{}, id).Error
}

// ==================== ADMIN USER QUERIES ====================

func (r *adminRepositoryEnhanced) GetAdminUsers(limit, offset int) ([]model.AdminUser, int64, error) {
	var admins []model.AdminUser
	var total int64

	if err := r.db.Model(&model.AdminUser{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("Role").
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&admins).Error

	return admins, total, err
}

func (r *adminRepositoryEnhanced) GetAdminUsersByRole(roleID uint, limit, offset int) ([]model.AdminUser, int64, error) {
	var admins []model.AdminUser
	var total int64

	if err := r.db.Model(&model.AdminUser{}).Where("role_id = ?", roleID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where("role_id = ?", roleID).
		Preload("Role").
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&admins).Error

	return admins, total, err
}

func (r *adminRepositoryEnhanced) GetAdminUsersByStatus(status model.AdminStatus, limit, offset int) ([]model.AdminUser, int64, error) {
	var admins []model.AdminUser
	var total int64

	if err := r.db.Model(&model.AdminUser{}).Where("status = ?", status).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where("status = ?", status).
		Preload("Role").
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&admins).Error

	return admins, total, err
}

// ==================== ADMIN USER STATUS ====================

func (r *adminRepositoryEnhanced) UpdateAdminStatus(id uint, status model.AdminStatus) error {
	return r.db.Model(&model.AdminUser{}).Where("id = ?", id).Update("status", status).Error
}

func (r *adminRepositoryEnhanced) UpdateAdminLastLogin(id uint, ip string) error {
	now := time.Now()
	return r.db.Model(&model.AdminUser{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_login_at": now,
		"last_login_ip": ip,
		"failed_login_attempts": 0,
		"locked_until": nil,
	}).Error
}

func (r *adminRepositoryEnhanced) UpdateAdminPassword(id uint, hashedPassword string) error {
	return r.db.Model(&model.AdminUser{}).Where("id = ?", id).Update("password", hashedPassword).Error
}

// ==================== ADMIN ROLES ====================

func (r *adminRepositoryEnhanced) GetAdminRoleByID(id uint) (*model.AdminRole, error) {
	var role model.AdminRole
	err := r.db.First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *adminRepositoryEnhanced) GetAdminRoleByName(name model.AdminRoleType) (*model.AdminRole, error) {
	var role model.AdminRole
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *adminRepositoryEnhanced) GetAllAdminRoles() ([]model.AdminRole, error) {
	var roles []model.AdminRole
	err := r.db.Order("id ASC").Find(&roles).Error
	return roles, err
}

func (r *adminRepositoryEnhanced) CreateAdminRole(role *model.AdminRole) error {
	return r.db.Create(role).Error
}

func (r *adminRepositoryEnhanced) UpdateAdminRole(role *model.AdminRole) error {
	return r.db.Save(role).Error
}

// ==================== AUDIT LOGS ====================

func (r *adminRepositoryEnhanced) CreateAuditLog(log *model.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *adminRepositoryEnhanced) GetAuditLogs(limit, offset int) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64

	if err := r.db.Model(&model.AuditLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("Admin").
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&logs).Error

	return logs, total, err
}

func (r *adminRepositoryEnhanced) GetAuditLogsByAdminID(adminID uint, limit, offset int) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64

	if err := r.db.Model(&model.AuditLog{}).Where("admin_id = ?", adminID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where("admin_id = ?", adminID).
		Preload("Admin").
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&logs).Error

	return logs, total, err
}

func (r *adminRepositoryEnhanced) GetAuditLogsByAction(action string, limit, offset int) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64

	if err := r.db.Model(&model.AuditLog{}).Where("action = ?", action).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where("action = ?", action).
		Preload("Admin").
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&logs).Error

	return logs, total, err
}

func (r *adminRepositoryEnhanced) GetAuditLogsByEntityType(entityType string, limit, offset int) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64

	if err := r.db.Model(&model.AuditLog{}).Where("entity_type = ?", entityType).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where("entity_type = ?", entityType).
		Preload("Admin").
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&logs).Error

	return logs, total, err
}

func (r *adminRepositoryEnhanced) GetAuditLogsByDateRange(startDate, endDate time.Time, limit, offset int) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64

	if err := r.db.Model(&model.AuditLog{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Preload("Admin").
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&logs).Error

	return logs, total, err
}

// ==================== SYSTEM SETTINGS ====================

func (r *adminRepositoryEnhanced) GetSystemSetting(key string) (*model.SystemSetting, error) {
	var setting model.SystemSetting
	err := r.db.Where("key = ?", key).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *adminRepositoryEnhanced) GetAllSystemSettings() ([]model.SystemSetting, error) {
	var settings []model.SystemSetting
	err := r.db.Order("key ASC").Find(&settings).Error
	return settings, err
}

func (r *adminRepositoryEnhanced) GetPublicSystemSettings() ([]model.SystemSetting, error) {
	var settings []model.SystemSetting
	err := r.db.Where("is_public = ?", true).Order("key ASC").Find(&settings).Error
	return settings, err
}

func (r *adminRepositoryEnhanced) UpdateSystemSetting(key string, value string, updatedBy *uint) error {
	now := time.Now()
	return r.db.Model(&model.SystemSetting{}).Where("key = ?", key).Updates(map[string]interface{}{
		"value":      value,
		"updated_by": updatedBy,
		"updated_at": now,
	}).Error
}

func (r *adminRepositoryEnhanced) CreateSystemSetting(setting *model.SystemSetting) error {
	return r.db.Create(setting).Error
}

func (r *adminRepositoryEnhanced) DeleteSystemSetting(key string) error {
	return r.db.Where("key = ?", key).Delete(&model.SystemSetting{}).Error
}

func (r *adminRepositoryEnhanced) CreateRefund(refund *model.Refund) error {
	return r.db.Create(refund).Error
}

// ==================== ANALYTICS ====================

func (r *adminRepositoryEnhanced) GetAdminStats() (*model.AdminStats, error) {
	stats := &model.AdminStats{}

	// Total users
	r.db.Model(&model.User{}).Where("deleted_at IS NULL").Count(&stats.TotalUsers)

	// Total sellers
	r.db.Model(&model.Shop{}).Where("deleted_at IS NULL").Count(&stats.TotalSellers)

	// Total products
	r.db.Model(&model.Product{}).Where("status = ? AND deleted_at IS NULL", model.ProductStatusActive).Count(&stats.TotalProducts)

	// Total orders
	r.db.Model(&model.Order{}).Where("deleted_at IS NULL").Count(&stats.TotalOrders)

	// Total revenue
	r.db.Model(&model.Order{}).
		Where("status NOT IN ?", []model.OrderStatus{model.OrderStatusCancelled, model.OrderStatusRefunded}).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&stats.TotalRevenue)

	// Pending orders
	r.db.Model(&model.Order{}).Where("status = ?", model.OrderStatusPending).Count(&stats.PendingOrders)

	// Pending refunds
	r.db.Model(&model.Refund{}).Where("status = ?", "pending").Count(&stats.PendingRefunds)

	// Active users 24h
	r.db.Model(&model.User{}).
		Where("last_login_at >= ?", time.Now().Add(-24*time.Hour)).
		Count(&stats.ActiveUsers24h)

	// New users today
	r.db.Model(&model.User{}).
		Where("created_at >= ?", time.Now().Truncate(24*time.Hour)).
		Count(&stats.NewUsersToday)

	return stats, nil
}

func (r *adminRepositoryEnhanced) GetSalesAnalytics(startDate, endDate time.Time) (*model.SalesAnalytics, error) {
	analytics := &model.SalesAnalytics{}

	// Total sales
	r.db.Model(&model.Order{}).
		Where("status NOT IN ?", []model.OrderStatus{model.OrderStatusCancelled, model.OrderStatusRefunded}).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&analytics.TotalSales)

	// Total orders
	r.db.Model(&model.Order{}).
		Where("status NOT IN ?", []model.OrderStatus{model.OrderStatusCancelled, model.OrderStatusRefunded}).
		Count(&analytics.TotalOrders)

	// Average order value
	if analytics.TotalOrders > 0 {
		analytics.AverageOrderValue = analytics.TotalSales / float64(analytics.TotalOrders)
	}

	// Today sales
	today := time.Now().Truncate(24 * time.Hour)
	r.db.Model(&model.Order{}).
		Where("created_at >= ?", today).
		Where("status NOT IN ?", []model.OrderStatus{model.OrderStatusCancelled, model.OrderStatusRefunded}).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&analytics.TodaySales)

	// Week sales
	weekAgo := time.Now().Add(-7 * 24 * time.Hour)
	r.db.Model(&model.Order{}).
		Where("created_at >= ?", weekAgo).
		Where("status NOT IN ?", []model.OrderStatus{model.OrderStatusCancelled, model.OrderStatusRefunded}).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&analytics.WeekSales)

	// Month sales
	monthAgo := time.Now().Add(-30 * 24 * time.Hour)
	r.db.Model(&model.Order{}).
		Where("created_at >= ?", monthAgo).
		Where("status NOT IN ?", []model.OrderStatus{model.OrderStatusCancelled, model.OrderStatusRefunded}).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&analytics.MonthSales)

	return analytics, nil
}

func (r *adminRepositoryEnhanced) GetUserAnalytics() (*model.UserAnalytics, error) {
	analytics := &model.UserAnalytics{}

	// Total users
	r.db.Model(&model.User{}).Where("deleted_at IS NULL").Count(&analytics.TotalUsers)

	// New users today
	today := time.Now().Truncate(24 * time.Hour)
	r.db.Model(&model.User{}).Where("created_at >= ?", today).Count(&analytics.NewUsersToday)

	// New users week
	weekAgo := time.Now().Add(-7 * 24 * time.Hour)
	r.db.Model(&model.User{}).Where("created_at >= ?", weekAgo).Count(&analytics.NewUsersWeek)

	// New users month
	monthAgo := time.Now().Add(-30 * 24 * time.Hour)
	r.db.Model(&model.User{}).Where("created_at >= ?", monthAgo).Count(&analytics.NewUsersMonth)

	// Active users (logged in within 30 days)
	r.db.Model(&model.User{}).
		Where("last_login_at >= ?", time.Now().Add(-30*24*time.Hour)).
		Count(&analytics.ActiveUsers)

	// Banned users
	r.db.Model(&model.User{}).Where("status = ?", model.StatusBanned).Count(&analytics.BannedUsers)

	return analytics, nil
}

func (r *adminRepositoryEnhanced) GetProductAnalytics(limit int) (*model.ProductAnalytics, error) {
	analytics := &model.ProductAnalytics{}

	// Total products
	r.db.Model(&model.Product{}).Where("deleted_at IS NULL").Count(&analytics.TotalProducts)

	// Active products
	r.db.Model(&model.Product{}).Where("status = ? AND deleted_at IS NULL", model.ProductStatusActive).Count(&analytics.ActiveProducts)

	// Out of stock
	r.db.Model(&model.Product{}).Where("stock <= 0 AND deleted_at IS NULL").Count(&analytics.OutOfStock)

	// Top products by sold count
	var topProducts []model.ProductStat
	r.db.Model(&model.Product{}).
		Where("status = ? AND deleted_at IS NULL", model.ProductStatusActive).
		Order("sold_count DESC").
		Limit(limit).
		Select("id, name, sold_count, (price * sold_count) as revenue").
		Scan(&topProducts)
	analytics.TopProducts = topProducts

	// Low stock products (stock < 10)
	var lowStockProducts []model.ProductStat
	r.db.Model(&model.Product{}).
		Where("stock > 0 AND stock < 10 AND deleted_at IS NULL").
		Order("stock ASC").
		Limit(limit).
		Select("id, name, stock as sold_count, price as revenue").
		Scan(&lowStockProducts)
	analytics.LowStockProducts = lowStockProducts

	return analytics, nil
}

// ==================== CLEANUP ====================

func (r *adminRepositoryEnhanced) DeleteOldAuditLogs(olderThan time.Duration) (int64, error) {
	cutoffTime := time.Now().Add(-olderThan)

	result := r.db.Where("created_at < ?", cutoffTime).Delete(&model.AuditLog{})

	return result.RowsAffected, result.Error
}

// Error definitions
var (
	ErrAdminNotFound    = errors.New("admin user not found")
	ErrRoleNotFound     = errors.New("admin role not found")
	ErrSettingNotFound  = errors.New("system setting not found")
	ErrDuplicateEmail   = errors.New("email already exists")
	ErrInvalidRole      = errors.New("invalid admin role")
)

// ==================== REVIEW MANAGEMENT ====================

func (r *adminRepositoryEnhanced) GetAllReviews(limit, offset int) ([]model.Review, int64, error) {
	var reviews []model.Review
	var total int64

	// Get total count
	r.db.Model(&model.Review{}).Count(&total)

	// Get reviews with pagination
	err := r.db.Model(&model.Review{}).
		Preload("User").
		Preload("Product").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&reviews).Error

	return reviews, total, err
}
