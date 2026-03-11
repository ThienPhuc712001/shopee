package handler

import (
	"errors"
	"fmt"
	"mime/multipart"
	"ecommerce/internal/service"
	"ecommerce/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UploadHandler handles image upload HTTP requests
type UploadHandler struct {
	uploadService *service.UploadService
}

// NewUploadHandler creates a new upload handler
func NewUploadHandler(uploadService *service.UploadService) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
	}
}

// UploadProductImageRequest represents the request body for product image upload
type UploadProductImageRequest struct {
	ProductID int64 `form:"product_id" binding:"required,min=1"`
	IsPrimary bool  `form:"is_primary"`
}

// UploadReviewImageRequest represents the request body for review image upload
type UploadReviewImageRequest struct {
	ReviewID int64 `form:"review_id" binding:"required,min=1"`
}

// UploadAvatarRequest represents the request body for avatar upload
type UploadAvatarRequest struct {
	UserID int64 `form:"user_id"`
}

// UploadProductImage handles product image upload
// @Summary Upload product image
// @Description Upload an image for a product
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Image file"
// @Param product_id formData int true "Product ID"
// @Param is_primary formData bool false "Is primary image"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/upload/product [POST]
func (h *UploadHandler) UploadProductImage(c *gin.Context) {
	// Parse form data
	var req UploadProductImageRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request parameters",
			"error":   err.Error(),
		})
		return
	}

	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "No file uploaded",
			"error":   err.Error(),
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		userID = int64(0) // Anonymous upload or handle as error
	}

	// Get client IP
	ipAddress := c.ClientIP()

	// Upload image
	result, err := h.uploadService.UploadProductImage(
		file,
		req.ProductID,
		userID.(int64),
		ipAddress,
		req.IsPrimary,
	)
	if err != nil {
		var fileErr *utils.FileError
		if errors.As(err, &fileErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": fileErr.Message,
				"error":   fileErr.Code,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to upload image",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Image uploaded successfully",
		"data":    result,
	})
}

// UploadMultipleProductImages handles multiple product image upload
// @Summary Upload multiple product images
// @Description Upload multiple images for a product
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param files formData []file true "Image files"
// @Param product_id formData int true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/upload/product/multiple [POST]
func (h *UploadHandler) UploadMultipleProductImages(c *gin.Context) {
	// Parse form data
	productIDStr := c.PostForm("product_id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil || productID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid product ID",
		})
		return
	}

	// Get uploaded files
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid form data",
			"error":   err.Error(),
		})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		files = make([]*multipart.FileHeader, 0)
		for _, fileHeaders := range form.File {
			files = append(files, fileHeaders...)
		}
	}

	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "No files uploaded",
		})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		userID = int64(0)
	}

	ipAddress := c.ClientIP()

	// Upload images
	results, err := h.uploadService.UploadMultipleProductImages(
		files,
		productID,
		userID.(int64),
		ipAddress,
	)
	if err != nil {
		var fileErr *utils.FileError
		if errors.As(err, &fileErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": fileErr.Message,
				"error":   fileErr.Code,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to upload images",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Uploaded %d images successfully", len(results)),
		"data": gin.H{
			"images": results,
			"count":  len(results),
		},
	})
}

// UploadReviewImage handles review image upload
// @Summary Upload review image
// @Description Upload an image for a review
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Image file"
// @Param review_id formData int true "Review ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/upload/review [POST]
func (h *UploadHandler) UploadReviewImage(c *gin.Context) {
	// Parse form data
	var req UploadReviewImageRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request parameters",
			"error":   err.Error(),
		})
		return
	}

	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "No file uploaded",
			"error":   err.Error(),
		})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		userID = int64(0)
	}

	ipAddress := c.ClientIP()

	// Upload image
	result, err := h.uploadService.UploadReviewImage(
		file,
		req.ReviewID,
		userID.(int64),
		ipAddress,
	)
	if err != nil {
		var fileErr *utils.FileError
		if errors.As(err, &fileErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": fileErr.Message,
				"error":   fileErr.Code,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to upload image",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Image uploaded successfully",
		"data":    result,
	})
}

// UploadAvatar handles user avatar upload
// @Summary Upload user avatar
// @Description Upload a user's avatar image
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Avatar image file"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/upload/avatar [POST]
func (h *UploadHandler) UploadAvatar(c *gin.Context) {
	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "No file uploaded",
			"error":   err.Error(),
		})
		return
	}

	// Get user ID from context (required for avatar upload)
	userID, exists := c.Get("user_id")
	if !exists || userID.(int64) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Authentication required",
		})
		return
	}

	ipAddress := c.ClientIP()

	// Upload avatar
	result, err := h.uploadService.UploadUserAvatar(
		file,
		userID.(int64),
		ipAddress,
	)
	if err != nil {
		var fileErr *utils.FileError
		if errors.As(err, &fileErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": fileErr.Message,
				"error":   fileErr.Code,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to upload avatar",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Avatar uploaded successfully",
		"data":    result,
	})
}

// DeleteProductImage handles product image deletion
// @Summary Delete product image
// @Description Delete a product image by ID
// @Tags upload
// @Produce json
// @Param id path int true "Image ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/upload/product/:id [DELETE]
func (h *UploadHandler) DeleteProductImage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid image ID",
		})
		return
	}

	if err := h.uploadService.DeleteProductImage(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete image",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Image deleted successfully",
	})
}

// GetProductImages handles getting all images for a product
// @Summary Get product images
// @Description Get all images for a product
// @Tags upload
// @Produce json
// @Param product_id query int true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/upload/product/images [GET]
func (h *UploadHandler) GetProductImages(c *gin.Context) {
	productIDStr := c.Query("product_id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil || productID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid product ID",
		})
		return
	}

	images, err := h.uploadService.GetProductImages(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get images",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"images": images,
			"count":  len(images),
		},
	})
}
