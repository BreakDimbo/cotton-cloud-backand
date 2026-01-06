package handlers

import (
	"net/http"

	"cotton-cloud-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ClothingHandler handles clothing-related requests
type ClothingHandler struct {
	db *gorm.DB
}

// NewClothingHandler creates a new ClothingHandler
func NewClothingHandler(db *gorm.DB) *ClothingHandler {
	return &ClothingHandler{db: db}
}

// List returns all clothing items for the current user
func (h *ClothingHandler) List(c *gin.Context) {
	// TODO: Get user ID from JWT token
	userID := c.Query("user_id")
	if userID == "" {
		userID = "demo-user" // Demo mode
	}

	var items []models.ClothingItem
	if err := h.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}

	c.JSON(http.StatusOK, items)
}

// Get returns a single clothing item by ID
func (h *ClothingHandler) Get(c *gin.Context) {
	id := c.Param("id")

	var item models.ClothingItem
	if err := h.db.First(&item, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch item"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// Create creates a new clothing item
func (h *ClothingHandler) Create(c *gin.Context) {
	var req models.CreateClothingItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Get user ID from JWT token
	userID := c.Query("user_id")
	if userID == "" {
		userID = "demo-user"
	}

	maxWearCount := 5
	if req.MaxWearCount != nil {
		maxWearCount = *req.MaxWearCount
	}

	item := models.ClothingItem{
		UserID:            userID,
		ImageURL:          req.ImageURL,
		OriginalImageURL:  req.OriginalImageURL,
		ProcessedImageURL: req.ProcessedImageURL,
		Category:          req.Category,
		Color:             req.Color,
		Material:          req.Material,
		Description:       req.Description,
		Tags:              req.Tags,
		Style:             req.Style,
		Season:            req.Season,
		MaxWearCount:      maxWearCount,
	}

	if err := h.db.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// Update updates an existing clothing item
func (h *ClothingHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var item models.ClothingItem
	if err := h.db.First(&item, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch item"})
		return
	}

	var req models.UpdateClothingItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields if provided
	if req.ImageURL != nil {
		item.ImageURL = *req.ImageURL
	}
	if req.Category != nil {
		item.Category = *req.Category
	}
	if req.Color != nil {
		item.Color = *req.Color
	}
	if req.Material != nil {
		item.Material = req.Material
	}
	if req.Description != nil {
		item.Description = req.Description
	}
	if req.Tags != nil {
		item.Tags = req.Tags
	}
	if req.Style != nil {
		item.Style = req.Style
	}
	if req.Season != nil {
		item.Season = req.Season
	}
	if req.MaxWearCount != nil {
		item.MaxWearCount = *req.MaxWearCount
	}

	if err := h.db.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// Delete removes a clothing item
func (h *ClothingHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	result := h.db.Delete(&models.ClothingItem{}, "id = ?", id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
}

// Wash resets the wear count for an item
func (h *ClothingHandler) Wash(c *gin.Context) {
	id := c.Param("id")

	result := h.db.Model(&models.ClothingItem{}).Where("id = ?", id).Updates(map[string]interface{}{
		"wear_count":     0,
		"last_washed_at": gorm.Expr("CURRENT_TIMESTAMP"),
	})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to wash item"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item washed"})
}

// IncrementWear increments the wear count for an item
func (h *ClothingHandler) IncrementWear(c *gin.Context) {
	id := c.Param("id")

	result := h.db.Model(&models.ClothingItem{}).Where("id = ?", id).Update("wear_count", gorm.Expr("wear_count + 1"))

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to increment wear count"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wear count incremented"})
}
