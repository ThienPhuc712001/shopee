package service

import (
	"ecommerce/internal/domain/model"
	"ecommerce/internal/repository"
	"errors"
	"fmt"
	"time"
)

// AdminServiceEnhanced defines the enhanced admin service interface
type AdminServiceEnhanced interface {
	// Admin Authentication
	AdminLogin(email, password, ip string) (*model.AdminUser, string, error)
	AdminLogout(adminID uint) error

	// Admin User Management
	CreateAdminUser(input *model.AdminCreateInput, creatorID uint) (*model.AdminUser, error)
	UpdateAdminUser(id uint, input *model.AdminUpdateInput, updaterID uint) (*model.AdminUser, error)
	DeleteAdminUser(id uint, deleterID uint) error
	GetAdminUser(id uint) (*model.AdminUser, error)
	GetAdminUsers(page, limit int) ([]model.AdminUser, int64, error)

	// User Management (Platform Users)
	BanUser(adminID uint, input *model.BanUserInput) (*model.User, error)
	UnbanUser(adminID uint, userID uint) (*model.User, error)
	GetUsers(page, limit int) ([]model.User, int64, error)

	// Seller Management
	ApproveSeller(adminID uint, input *model.ApproveSellerInput) (*model.Shop, error)
	RejectSeller(adminID uint, shopID uint, reason string) error
	SuspendSeller(adminID uint, shopID uint, reason string) (*model.Shop, error)
	GetPendingSellers(page, limit int) ([]model.Shop, int64, error)

	// Product Management
	DeleteProduct(adminID uint, productID uint, reason string) error
	RestoreProduct(adminID uint, productID uint) error
	GetProductsForModeration(page, limit int) ([]model.Product, int64, error)

	// Order Management
	GetOrders(page, limit int) ([]model.Order, int64, error)
	GetOrder(adminID uint, orderID uint) (*model.Order, error)
	RefundOrder(adminID uint, input *model.RefundOrderInput) (*model.Refund, error)
	CancelOrder(adminID uint, orderID uint, reason string) (*model.Order, error)

	// Analytics
	GetAdminStats() (*model.AdminStats, error)
	GetSalesAnalytics(startDate, endDate time.Time) (*model.SalesAnalytics, error)
	GetUserAnalytics() (*model.UserAnalytics, error)
	GetProductAnalytics(limit int) (*model.ProductAnalytics, error)

	// System Settings
	GetSystemSetting(key string) (*model.SystemSetting, error)
	UpdateSystemSetting(key, value string, adminID uint) (*model.SystemSetting, error)
	GetAllSystemSettings() ([]model.SystemSetting, error)

	// Audit Logs
	GetAuditLogs(page, limit int) ([]model.AuditLog, int64, error)
	GetAuditLogsByAdminID(adminID uint, page, limit int) ([]model.AuditLog, int64, error)
	CreateAuditLog(input *model.AuditLogInput) error

	// Review Management
	GetAllReviews(page, limit int) ([]model.Review, int64, error)
}

type adminServiceEnhanced struct {
	adminRepo  repository.AdminRepositoryEnhanced
	userRepo   repository.UserRepositoryEnhanced
	shopRepo   repository.ShopRepositoryEnhanced
	productRepo repository.ProductRepositoryEnhanced
	orderRepo  repository.OrderRepositoryEnhanced
}

// NewAdminServiceEnhanced creates a new enhanced admin service
func NewAdminServiceEnhanced(
	adminRepo repository.AdminRepositoryEnhanced,
	userRepo repository.UserRepositoryEnhanced,
	shopRepo repository.ShopRepositoryEnhanced,
	productRepo repository.ProductRepositoryEnhanced,
	orderRepo repository.OrderRepositoryEnhanced,
) AdminServiceEnhanced {
	return &adminServiceEnhanced{
		adminRepo:   adminRepo,
		userRepo:    userRepo,
		shopRepo:    shopRepo,
		productRepo: productRepo,
		orderRepo:   orderRepo,
	}
}

// ==================== ADMIN AUTHENTICATION ====================

func (s *adminServiceEnhanced) AdminLogin(email, password, ip string) (*model.AdminUser, string, error) {
	admin, err := s.adminRepo.GetAdminUserByEmail(email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// Check if account is locked
	if admin.IsLocked() {
		return nil, "", errors.New("account is locked")
	}

	// Check status
	if admin.Status != model.AdminStatusActive {
		return nil, "", errors.New("account is inactive")
	}

	// Verify password
	if !admin.CheckPassword(password) {
		admin.IncrementFailedLogin()
		s.adminRepo.UpdateAdminUser(admin)
		return nil, "", ErrInvalidCredentials
	}

	// Reset failed login
	admin.ResetFailedLogin()
	s.adminRepo.UpdateAdminLastLogin(admin.ID, ip)

	// Generate JWT token (simplified - use actual JWT service in production)
	token := fmt.Sprintf("admin_token_%d_%d", admin.ID, time.Now().Unix())

	return admin, token, nil
}

func (s *adminServiceEnhanced) AdminLogout(adminID uint) error {
	// In production, invalidate token
	return nil
}

// ==================== ADMIN USER MANAGEMENT ====================

func (s *adminServiceEnhanced) CreateAdminUser(input *model.AdminCreateInput, creatorID uint) (*model.AdminUser, error) {
	// Check if email already exists
	existing, _ := s.adminRepo.GetAdminUserByEmail(input.Email)
	if existing != nil {
		return nil, ErrAdminAlreadyExists
	}

	// Validate role
	_, err := s.adminRepo.GetAdminRoleByID(input.RoleID)
	if err != nil {
		return nil, ErrInvalidRole
	}

	// Create admin user
	admin := &model.AdminUser{
		Email:     input.Email,
		RoleID:    input.RoleID,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.Phone,
		Status:    model.AdminStatusActive,
	}

	// Hash password
	if err := admin.HashPassword(input.Password); err != nil {
		return nil, err
	}

	if err := s.adminRepo.CreateAdminUser(admin); err != nil {
		return nil, err
	}

	// Create audit log
	s.CreateAuditLog(&model.AuditLogInput{
		AdminID:    creatorID,
		Action:     model.ActionCreate,
		EntityType: model.EntityAdmin,
		EntityID:   &admin.ID,
		NewValues:  fmt.Sprintf(`{"email": "%s", "role_id": %d}`, admin.Email, admin.RoleID),
	})

	return admin, nil
}

func (s *adminServiceEnhanced) UpdateAdminUser(id uint, input *model.AdminUpdateInput, updaterID uint) (*model.AdminUser, error) {
	admin, err := s.adminRepo.GetAdminUserByID(id)
	if err != nil {
		return nil, ErrAdminNotFound
	}

	oldValues := fmt.Sprintf(`{"role_id": %d, "status": "%s"}`, admin.RoleID, admin.Status)

	if input.FirstName != "" {
		admin.FirstName = input.FirstName
	}
	if input.LastName != "" {
		admin.LastName = input.LastName
	}
	if input.Phone != "" {
		admin.Phone = input.Phone
	}
	if input.AvatarURL != "" {
		admin.AvatarURL = input.AvatarURL
	}
	if input.RoleID > 0 {
		admin.RoleID = input.RoleID
	}
	if input.Status != "" {
		admin.Status = input.Status
	}

	if err := s.adminRepo.UpdateAdminUser(admin); err != nil {
		return nil, err
	}

	// Create audit log
	s.CreateAuditLog(&model.AuditLogInput{
		AdminID:    updaterID,
		Action:     model.ActionUpdate,
		EntityType: model.EntityAdmin,
		EntityID:   &admin.ID,
		OldValues:  oldValues,
		NewValues:  fmt.Sprintf(`{"role_id": %d, "status": "%s"}`, admin.RoleID, admin.Status),
	})

	return admin, nil
}

func (s *adminServiceEnhanced) DeleteAdminUser(id uint, deleterID uint) error {
	admin, err := s.adminRepo.GetAdminUserByID(id)
	if err != nil {
		return ErrAdminNotFound
	}

	// Prevent deleting self
	if id == deleterID {
		return errors.New("cannot delete your own account")
	}

	if err := s.adminRepo.DeleteAdminUser(id); err != nil {
		return err
	}

	// Create audit log
	s.CreateAuditLog(&model.AuditLogInput{
		AdminID:    deleterID,
		Action:     model.ActionDelete,
		EntityType: model.EntityAdmin,
		EntityID:   &admin.ID,
	})

	return nil
}

func (s *adminServiceEnhanced) GetAdminUser(id uint) (*model.AdminUser, error) {
	return s.adminRepo.GetAdminUserByID(id)
}

func (s *adminServiceEnhanced) GetAdminUsers(page, limit int) ([]model.AdminUser, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	return s.adminRepo.GetAdminUsers(limit, offset)
}

// ==================== USER MANAGEMENT ====================

func (s *adminServiceEnhanced) BanUser(adminID uint, input *model.BanUserInput) (*model.User, error) {
	user, err := s.userRepo.FindByID(input.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if err := s.userRepo.UpdateStatus(input.UserID, model.StatusBanned); err != nil {
		return nil, err
	}

	// Create audit log
	s.CreateAuditLog(&model.AuditLogInput{
		AdminID:    adminID,
		Action:     model.ActionBan,
		EntityType: model.EntityUser,
		EntityID:   &input.UserID,
		NewValues:  fmt.Sprintf(`{"status": "banned", "reason": "%s"}`, input.Reason),
	})

	return user, nil
}

func (s *adminServiceEnhanced) UnbanUser(adminID uint, userID uint) (*model.User, error) {
	if err := s.userRepo.UpdateStatus(userID, model.StatusActive); err != nil {
		return nil, err
	}

	// Create audit log
	s.CreateAuditLog(&model.AuditLogInput{
		AdminID:    adminID,
		Action:     model.ActionUnban,
		EntityType: model.EntityUser,
		EntityID:   &userID,
		NewValues:  `{"status": "active"}`,
	})

	return s.userRepo.FindByID(userID)
}

func (s *adminServiceEnhanced) GetUsers(page, limit int) ([]model.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	return s.userRepo.FindAll(limit, offset)
}

// ==================== SELLER MANAGEMENT ====================

func (s *adminServiceEnhanced) ApproveSeller(adminID uint, input *model.ApproveSellerInput) (*model.Shop, error) {
	shop, err := s.shopRepo.FindByID(input.ShopID)
	if err != nil {
		return nil, errors.New("shop not found")
	}

	shop.Status = model.ShopStatusActive
	shop.VerificationStatus = "verified"
	if err := s.shopRepo.Update(shop); err != nil {
		return nil, err
	}

	// Create audit log
	s.CreateAuditLog(&model.AuditLogInput{
		AdminID:    adminID,
		Action:     model.ActionApprove,
		EntityType: model.EntityShop,
		EntityID:   &input.ShopID,
		NewValues:  fmt.Sprintf(`{"status": "active", "verification": "verified", "notes": "%s"}`, input.Notes),
	})

	return shop, nil
}

func (s *adminServiceEnhanced) RejectSeller(adminID uint, shopID uint, reason string) error {
	shop, err := s.shopRepo.FindByID(shopID)
	if err != nil {
		return errors.New("shop not found")
	}

	shop.Status = model.ShopStatusInactive
	if err := s.shopRepo.Update(shop); err != nil {
		return err
	}

	// Create audit log
	s.CreateAuditLog(&model.AuditLogInput{
		AdminID:    adminID,
		Action:     model.ActionReject,
		EntityType: model.EntityShop,
		EntityID:   &shopID,
		NewValues:  fmt.Sprintf(`{"status": "inactive", "reason": "%s"}`, reason),
	})

	return nil
}

func (s *adminServiceEnhanced) SuspendSeller(adminID uint, shopID uint, reason string) (*model.Shop, error) {
	shop, err := s.shopRepo.FindByID(shopID)
	if err != nil {
		return nil, errors.New("shop not found")
	}

	shop.Status = model.ShopStatusSuspended
	if err := s.shopRepo.Update(shop); err != nil {
		return nil, err
	}

	// Create audit log
	s.CreateAuditLog(&model.AuditLogInput{
		AdminID:    adminID,
		Action:     "suspend",
		EntityType: model.EntityShop,
		EntityID:   &shopID,
		NewValues:  fmt.Sprintf(`{"status": "suspended", "reason": "%s"}`, reason),
	})

	return shop, nil
}

func (s *adminServiceEnhanced) GetPendingSellers(page, limit int) ([]model.Shop, int64, error) {
	// In production, add a method to repository to filter by status
	return s.shopRepo.FindAll(limit, (page-1)*limit)
}

// ==================== PRODUCT MANAGEMENT ====================

func (s *adminServiceEnhanced) DeleteProduct(adminID uint, productID uint, reason string) error {
	_, err := s.productRepo.FindByID(productID)
	if err != nil {
		return errors.New("product not found")
	}

	if err := s.productRepo.Delete(productID); err != nil {
		return err
	}

	// Create audit log
	s.CreateAuditLog(&model.AuditLogInput{
		AdminID:    adminID,
		Action:     model.ActionDelete,
		EntityType: model.EntityProduct,
		EntityID:   &productID,
		NewValues:  fmt.Sprintf(`{"product_id": %d, "reason": "%s"}`, productID, reason),
	})

	return nil
}

func (s *adminServiceEnhanced) RestoreProduct(adminID uint, productID uint) error {
	// In production, implement soft delete restore
	return nil
}

func (s *adminServiceEnhanced) GetProductsForModeration(page, limit int) ([]model.Product, int64, error) {
	return s.productRepo.FindAll(limit, (page-1)*limit)
}

// ==================== ORDER MANAGEMENT ====================

func (s *adminServiceEnhanced) GetOrders(page, limit int) ([]model.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	return s.orderRepo.FindAll(limit, offset)
}

func (s *adminServiceEnhanced) GetOrder(adminID uint, orderID uint) (*model.Order, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	// Create audit log for viewing
	s.CreateAuditLog(&model.AuditLogInput{
		AdminID:    adminID,
		Action:     "view",
		EntityType: model.EntityOrder,
		EntityID:   &orderID,
	})

	return order, nil
}

func (s *adminServiceEnhanced) RefundOrder(adminID uint, input *model.RefundOrderInput) (*model.Refund, error) {
	order, err := s.orderRepo.FindByID(input.OrderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	refundAmount := input.Amount
	if refundAmount <= 0 {
		refundAmount = order.TotalAmount
	}

	refund := &model.Refund{
		PaymentID: 0, // Get from payment service
		OrderID:   input.OrderID,
		Amount:    refundAmount,
		Reason:    input.Reason,
		Status:    "pending",
		Type:      "full",
	}

	if err := s.adminRepo.CreateRefund(refund); err != nil {
		return nil, err
	}

	// Create audit log
	s.CreateAuditLog(&model.AuditLogInput{
		AdminID:    adminID,
		Action:     model.ActionRefund,
		EntityType: model.EntityRefund,
		EntityID:   &refund.ID,
		NewValues:  fmt.Sprintf(`{"order_id": %d, "amount": %.2f, "reason": "%s"}`, input.OrderID, refundAmount, input.Reason),
	})

	return refund, nil
}

func (s *adminServiceEnhanced) CancelOrder(adminID uint, orderID uint, reason string) (*model.Order, error) {
	if err := s.orderRepo.UpdateStatus(orderID, model.OrderStatusCancelled); err != nil {
		return nil, err
	}

	// Create audit log
	s.CreateAuditLog(&model.AuditLogInput{
		AdminID:    adminID,
		Action:     model.ActionCancel,
		EntityType: model.EntityOrder,
		EntityID:   &orderID,
		NewValues:  fmt.Sprintf(`{"status": "cancelled", "reason": "%s"}`, reason),
	})

	return s.orderRepo.FindByID(orderID)
}

// ==================== ANALYTICS ====================

func (s *adminServiceEnhanced) GetAdminStats() (*model.AdminStats, error) {
	return s.adminRepo.GetAdminStats()
}

func (s *adminServiceEnhanced) GetSalesAnalytics(startDate, endDate time.Time) (*model.SalesAnalytics, error) {
	return s.adminRepo.GetSalesAnalytics(startDate, endDate)
}

func (s *adminServiceEnhanced) GetUserAnalytics() (*model.UserAnalytics, error) {
	return s.adminRepo.GetUserAnalytics()
}

func (s *adminServiceEnhanced) GetProductAnalytics(limit int) (*model.ProductAnalytics, error) {
	return s.adminRepo.GetProductAnalytics(limit)
}

// ==================== SYSTEM SETTINGS ====================

func (s *adminServiceEnhanced) GetSystemSetting(key string) (*model.SystemSetting, error) {
	return s.adminRepo.GetSystemSetting(key)
}

func (s *adminServiceEnhanced) UpdateSystemSetting(key, value string, adminID uint) (*model.SystemSetting, error) {
	setting, err := s.adminRepo.GetSystemSetting(key)
	if err != nil {
		// Create new setting
		setting = &model.SystemSetting{
			Key:   key,
			Value: value,
			Type:  "string",
		}
		if err := s.adminRepo.CreateSystemSetting(setting); err != nil {
			return nil, err
		}
	} else {
		oldValue := setting.Value
		if err := s.adminRepo.UpdateSystemSetting(key, value, &adminID); err != nil {
			return nil, err
		}

		// Create audit log
		s.CreateAuditLog(&model.AuditLogInput{
			AdminID:    adminID,
			Action:     model.ActionSystemSetting,
			EntityType: model.EntitySystem,
			NewValues:  fmt.Sprintf(`{"key": "%s", "old_value": "%s", "new_value": "%s"}`, key, oldValue, value),
		})
	}

	return s.adminRepo.GetSystemSetting(key)
}

func (s *adminServiceEnhanced) GetAllSystemSettings() ([]model.SystemSetting, error) {
	return s.adminRepo.GetAllSystemSettings()
}

// ==================== AUDIT LOGS ====================

func (s *adminServiceEnhanced) GetAuditLogs(page, limit int) ([]model.AuditLog, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	return s.adminRepo.GetAuditLogs(limit, offset)
}

func (s *adminServiceEnhanced) GetAuditLogsByAdminID(adminID uint, page, limit int) ([]model.AuditLog, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	return s.adminRepo.GetAuditLogsByAdminID(adminID, limit, offset)
}

func (s *adminServiceEnhanced) CreateAuditLog(input *model.AuditLogInput) error {
	log := &model.AuditLog{
		AdminID:    input.AdminID,
		Action:     input.Action,
		EntityType: input.EntityType,
		EntityID:   input.EntityID,
		OldValues:  input.OldValues,
		NewValues:  input.NewValues,
		IPAddress:  input.IPAddress,
		UserAgent:  input.UserAgent,
	}

	return s.adminRepo.CreateAuditLog(log)
}

// GetAllReviews retrieves all reviews for moderation
func (s *adminServiceEnhanced) GetAllReviews(page, limit int) ([]model.Review, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	// Get reviews from repository
	// Note: This assumes the adminRepo has a method to get all reviews
	// If not, we'll need to add it or use a different repository
	return s.adminRepo.GetAllReviews(limit, offset)
}
