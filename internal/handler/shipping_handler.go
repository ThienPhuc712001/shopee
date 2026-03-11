package handler

import (
	"ecommerce/internal/domain/model"
	"ecommerce/internal/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ShippingHandler handles shipping HTTP requests
type ShippingHandler struct {
	shippingService *service.ShippingService
}

// NewShippingHandler creates a new shipping handler
func NewShippingHandler(shippingService *service.ShippingService) *ShippingHandler {
	return &ShippingHandler{shippingService: shippingService}
}

// ==================== SHIPPING ADDRESS ====================

// CreateAddress handles POST /api/shipping/addresses
func (h *ShippingHandler) CreateAddress(c *gin.Context) {
	var input model.ShippingAddressInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid input",
			"error":   err.Error(),
		})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Authentication required",
		})
		return
	}

	address, err := h.shippingService.CreateAddress(c.Request.Context(), userID.(uint), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to create address",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Address created successfully",
		"data":    address,
	})
}

// GetAddresses handles GET /api/shipping/addresses
func (h *ShippingHandler) GetAddresses(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Authentication required",
		})
		return
	}

	addresses, err := h.shippingService.GetAddressesByUser(c.Request.Context(), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get addresses",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"addresses": addresses,
			"count":     len(addresses),
		},
	})
}

// GetAddressByID handles GET /api/shipping/addresses/:id
func (h *ShippingHandler) GetAddressByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid address ID",
		})
		return
	}

	address, err := h.shippingService.GetAddressByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, service.ErrAddressNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Address not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get address",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    address,
	})
}

// UpdateAddress handles PUT /api/shipping/addresses/:id
func (h *ShippingHandler) UpdateAddress(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid address ID",
		})
		return
	}

	var input model.ShippingAddressInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid input",
			"error":   err.Error(),
		})
		return
	}

	address, err := h.shippingService.UpdateAddress(c.Request.Context(), uint(id), input)
	if err != nil {
		if errors.Is(err, service.ErrAddressNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Address not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update address",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Address updated successfully",
		"data":    address,
	})
}

// DeleteAddress handles DELETE /api/shipping/addresses/:id
func (h *ShippingHandler) DeleteAddress(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid address ID",
		})
		return
	}

	if err := h.shippingService.DeleteAddress(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, service.ErrAddressNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Address not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete address",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Address deleted successfully",
	})
}

// SetDefaultAddress handles POST /api/shipping/addresses/:id/set-default
func (h *ShippingHandler) SetDefaultAddress(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid address ID",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Authentication required",
		})
		return
	}

	if err := h.shippingService.SetDefaultAddress(c.Request.Context(), userID.(uint), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to set default address",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Default address updated successfully",
	})
}

// ==================== SHIPMENT ====================

// CreateShipment handles POST /api/shipments
func (h *ShippingHandler) CreateShipment(c *gin.Context) {
	var input service.CreateShipmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid input",
			"error":   err.Error(),
		})
		return
	}

	shipment, err := h.shippingService.CreateShipment(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, service.ErrShipmentExists) {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "Shipment already exists for this order",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to create shipment",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Shipment created successfully",
		"data":    shipment,
	})
}

// GetShipmentByID handles GET /api/shipments/:id
func (h *ShippingHandler) GetShipmentByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid shipment ID",
		})
		return
	}

	shipment, err := h.shippingService.GetShipmentByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, service.ErrShipmentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Shipment not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get shipment",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    shipment,
	})
}

// GetShipmentByOrderID handles GET /api/orders/:id/shipment
func (h *ShippingHandler) GetShipmentByOrderID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid order ID",
		})
		return
	}

	shipment, err := h.shippingService.GetShipmentByOrderID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, service.ErrShipmentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Shipment not found for this order",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get shipment",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    shipment,
	})
}

// UpdateShipmentStatus handles PUT /api/shipments/:id/status
func (h *ShippingHandler) UpdateShipmentStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid shipment ID",
		})
		return
	}

	var input struct {
		Status string `json:"status" binding:"required"`
		Notes  string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid input",
			"error":   err.Error(),
		})
		return
	}

	status := model.ShipmentStatus(input.Status)
	shipment, err := h.shippingService.UpdateShipmentStatus(c.Request.Context(), uint(id), status, input.Notes)
	if err != nil {
		if errors.Is(err, service.ErrShipmentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Shipment not found",
			})
			return
		}
		if errors.Is(err, service.ErrInvalidStatusTransition) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid status transition",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update shipment status",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Shipment status updated successfully",
		"data":    shipment,
	})
}

// AddTrackingEvent handles POST /api/shipments/:id/tracking
func (h *ShippingHandler) AddTrackingEvent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid shipment ID",
		})
		return
	}

	var input model.TrackingEventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid input",
			"error":   err.Error(),
		})
		return
	}

	event, err := h.shippingService.AddTrackingEvent(c.Request.Context(), uint(id), input)
	if err != nil {
		if errors.Is(err, service.ErrShipmentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Shipment not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to add tracking event",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Tracking event added successfully",
		"data":    event,
	})
}

// GetTrackingTimeline handles GET /api/shipments/:id/tracking
func (h *ShippingHandler) GetTrackingTimeline(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid shipment ID",
		})
		return
	}

	timeline, err := h.shippingService.GetTrackingTimeline(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, service.ErrShipmentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Shipment not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get tracking timeline",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    timeline,
	})
}

// GetOrderTracking handles GET /api/orders/:id/tracking
func (h *ShippingHandler) GetOrderTracking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid order ID",
		})
		return
	}

	timeline, err := h.shippingService.GetOrderTracking(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, service.ErrShipmentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "No tracking information available for this order",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get order tracking",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    timeline,
	})
}

// GetShipmentsByStatus handles GET /api/shipments?status=pending
func (h *ShippingHandler) GetShipmentsByStatus(c *gin.Context) {
	status := c.Query("status")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	shipments, total, err := h.shippingService.GetShipmentsByStatus(c.Request.Context(), model.ShipmentStatus(status), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get shipments",
			"error":   err.Error(),
		})
		return
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"shipments": shipments,
			"pagination": gin.H{
				"total":       total,
				"page":        page,
				"limit":       limit,
				"total_pages": totalPages,
			},
		},
	})
}

// GetShipmentStats handles GET /api/shipments/stats
func (h *ShippingHandler) GetShipmentStats(c *gin.Context) {
	stats, err := h.shippingService.GetShipmentStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get shipment statistics",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// ==================== CARRIER ====================

// GetActiveCarriers handles GET /api/shipping/carriers
func (h *ShippingHandler) GetActiveCarriers(c *gin.Context) {
	carriers, err := h.shippingService.GetActiveCarriers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get carriers",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"carriers": carriers,
			"count":    len(carriers),
		},
	})
}
