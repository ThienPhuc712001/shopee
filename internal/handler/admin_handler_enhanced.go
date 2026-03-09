package handler

import (
	"ecommerce/internal/domain/model"
	"ecommerce/internal/service"
	"ecommerce/pkg/response"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// AdminHandlerEnhanced handles admin-related requests
type AdminHandlerEnhanced struct {
	adminService service.AdminServiceEnhanced
}

// NewAdminHandlerEnhanced creates a new enhanced admin handler
func NewAdminHandlerEnhanced(adminService service.AdminServiceEnhanced) *AdminHandlerEnhanced {
	return &AdminHandlerEnhanced{
		adminService: adminService,
	}
}

// AdminLoginRequest represents admin login request
type AdminLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AdminLogin handles admin authentication
// @Summary Admin login
// @Description Authenticate admin user
// @Tags admin
// @Accept json
// @Produce json
// @Param request body AdminLoginRequest true "Login credentials"
// @Success 200 {object} response.Response
// @Router /api/admin/auth/login [post]
func (h *AdminHandlerEnhanced) AdminLogin(c *gin.Context) {
	var req AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	ip := c.ClientIP()
	admin, token, err := h.adminService.AdminLogin(req.Email, req.Password, ip)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Unauthorized(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"admin": admin,
		"token": token,
	}, "Login successful"))
}

// GetUsers handles getting all users
// @Summary Get all users
// @Description Get all platform users (Admin only)
// @Tags admin/users
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} response.PaginatedResponse
// @Router /api/admin/users [get]
func (h *AdminHandlerEnhanced) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	users, total, err := h.adminService.GetUsers(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalError("Failed to get users"))
		return
	}

	c.JSON(http.StatusOK, response.Paginated(gin.H{
		"users": users,
	}, total, page, limit, ""))
}

// BanUser handles banning a user
// @Summary Ban user
// @Description Ban a user account (Admin only)
// @Tags admin/users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body BanUserRequest true "Ban data"
// @Success 200 {object} response.Response
// @Router /api/admin/users/ban [post]
func (h *AdminHandlerEnhanced) BanUser(c *gin.Context) {
	adminID, _ := c.Get("admin_id")

	var req model.BanUserInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	user, err := h.adminService.BanUser(adminID.(uint), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"user": user,
	}, "User banned successfully"))
}

// ApproveSeller handles approving a seller
// @Summary Approve seller
// @Description Approve a seller application (Admin only)
// @Tags admin/sellers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ApproveSellerRequest true "Approval data"
// @Success 200 {object} response.Response
// @Router /api/admin/sellers/approve [post]
func (h *AdminHandlerEnhanced) ApproveSeller(c *gin.Context) {
	adminID, _ := c.Get("admin_id")

	var req model.ApproveSellerInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	shop, err := h.adminService.ApproveSeller(adminID.(uint), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"shop": shop,
	}, "Seller approved successfully"))
}

// DeleteProduct handles deleting a product
// @Summary Delete product
// @Description Delete a product (Admin only)
// @Tags admin/products
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} response.Response
// @Router /api/admin/products/{id} [delete]
func (h *AdminHandlerEnhanced) DeleteProduct(c *gin.Context) {
	adminID, _ := c.Get("admin_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest("Invalid product ID"))
		return
	}

	if err := h.adminService.DeleteProduct(adminID.(uint), uint(id), "Admin deletion"); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithMessage("Product deleted successfully"))
}

// GetOrders handles getting all orders
// @Summary Get all orders
// @Description Get all platform orders (Admin only)
// @Tags admin/orders
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} response.PaginatedResponse
// @Router /api/admin/orders [get]
func (h *AdminHandlerEnhanced) GetOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	orders, total, err := h.adminService.GetOrders(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalError("Failed to get orders"))
		return
	}

	c.JSON(http.StatusOK, response.Paginated(gin.H{
		"orders": orders,
	}, total, page, limit, ""))
}

// RefundOrder handles refunding an order
// @Summary Refund order
// @Description Process a refund for an order (Admin only)
// @Tags admin/orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RefundOrderRequest true "Refund data"
// @Success 200 {object} response.Response
// @Router /api/admin/orders/refund [post]
func (h *AdminHandlerEnhanced) RefundOrder(c *gin.Context) {
	adminID, _ := c.Get("admin_id")

	var req model.RefundOrderInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	refund, err := h.adminService.RefundOrder(adminID.(uint), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"refund": refund,
	}, "Refund processed successfully"))
}

// GetAdminStats handles getting admin statistics
// @Summary Get admin statistics
// @Description Get platform statistics (Admin only)
// @Tags admin/analytics
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/admin/analytics/stats [get]
func (h *AdminHandlerEnhanced) GetAdminStats(c *gin.Context) {
	stats, err := h.adminService.GetAdminStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalError("Failed to get statistics"))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"stats": stats,
	}, ""))
}

// GetSalesAnalytics handles getting sales analytics
// @Summary Get sales analytics
// @Description Get sales analytics (Admin only)
// @Tags admin/analytics
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} response.Response
// @Router /api/admin/analytics/sales [get]
func (h *AdminHandlerEnhanced) GetSalesAnalytics(c *gin.Context) {
	startDateStr := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDateStr := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	startDate, _ := time.Parse("2006-01-02", startDateStr)
	endDate, _ := time.Parse("2006-01-02", endDateStr)

	analytics, err := h.adminService.GetSalesAnalytics(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalError("Failed to get sales analytics"))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"sales": analytics,
	}, ""))
}

// GetUserAnalytics handles getting user analytics
// @Summary Get user analytics
// @Description Get user analytics (Admin only)
// @Tags admin/analytics
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/admin/analytics/users [get]
func (h *AdminHandlerEnhanced) GetUserAnalytics(c *gin.Context) {
	analytics, err := h.adminService.GetUserAnalytics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalError("Failed to get user analytics"))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"users": analytics,
	}, ""))
}

// GetProductAnalytics handles getting product analytics
// @Summary Get product analytics
// @Description Get product analytics (Admin only)
// @Tags admin/analytics
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Number of products"
// @Success 200 {object} response.Response
// @Router /api/admin/analytics/products [get]
func (h *AdminHandlerEnhanced) GetProductAnalytics(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	analytics, err := h.adminService.GetProductAnalytics(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalError("Failed to get product analytics"))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"products": analytics,
	}, ""))
}

// GetAuditLogs handles getting audit logs
// @Summary Get audit logs
// @Description Get admin audit logs (Admin only)
// @Tags admin/audit-logs
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} response.PaginatedResponse
// @Router /api/admin/audit-logs [get]
func (h *AdminHandlerEnhanced) GetAuditLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	logs, total, err := h.adminService.GetAuditLogs(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalError("Failed to get audit logs"))
		return
	}

	c.JSON(http.StatusOK, response.Paginated(gin.H{
		"logs": logs,
	}, total, page, limit, ""))
}

// GetSystemSetting handles getting a system setting
// @Summary Get system setting
// @Description Get a system setting by key (Admin only)
// @Tags admin/settings
// @Produce json
// @Security BearerAuth
// @Param key path string true "Setting key"
// @Success 200 {object} response.Response
// @Router /api/admin/settings/{key} [get]
func (h *AdminHandlerEnhanced) GetSystemSetting(c *gin.Context) {
	key := c.Param("key")

	setting, err := h.adminService.GetSystemSetting(key)
	if err != nil {
		c.JSON(http.StatusNotFound, response.NotFound("Setting not found"))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"setting": setting,
	}, ""))
}

// UpdateSystemSetting handles updating a system setting
// @Summary Update system setting
// @Description Update a system setting (Admin only)
// @Tags admin/settings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param key path string true "Setting key"
// @Param request body UpdateSettingRequest true "Setting value"
// @Success 200 {object} response.Response
// @Router /api/admin/settings/{key} [put]
func (h *AdminHandlerEnhanced) UpdateSystemSetting(c *gin.Context) {
	adminID, _ := c.Get("admin_id")
	key := c.Param("key")

	var req struct {
		Value string `json:"value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BadRequest(err.Error()))
		return
	}

	setting, err := h.adminService.UpdateSystemSetting(key, req.Value, adminID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalError("Failed to update setting"))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"setting": setting,
	}, "Setting updated successfully"))
}

// BanUserRequest represents ban user request
type BanUserRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	Reason string `json:"reason" binding:"required"`
}

// ApproveSellerRequest represents approve seller request
type ApproveSellerRequest struct {
	ShopID uint   `json:"shop_id" binding:"required"`
	Notes  string `json:"notes"`
}

// RefundOrderRequest represents refund order request
type RefundOrderRequest struct {
	OrderID uint    `json:"order_id" binding:"required"`
	Amount  float64 `json:"amount"`
	Reason  string  `json:"reason" binding:"required"`
}

// UpdateSettingRequest represents update setting request
type UpdateSettingRequest struct {
	Value string `json:"value" binding:"required"`
}
