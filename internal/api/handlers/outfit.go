package handlers

import (
	"net/http"

	"cotton-cloud-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OutfitHandler handles outfit record-related requests
type OutfitHandler struct {
	db *gorm.DB
}

// NewOutfitHandler creates a new OutfitHandler
func NewOutfitHandler(db *gorm.DB) *OutfitHandler {
	return &OutfitHandler{db: db}
}

// List returns all outfit records for the current user
func (h *OutfitHandler) List(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		userID = "demo-user"
	}

	var records []models.OutfitRecord
	if err := h.db.Where("user_id = ?", userID).Order("date DESC").Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch records"})
		return
	}

	c.JSON(http.StatusOK, records)
}

// GetByDate returns an outfit record for a specific date
func (h *OutfitHandler) GetByDate(c *gin.Context) {
	date := c.Param("date")
	userID := c.Query("user_id")
	if userID == "" {
		userID = "demo-user"
	}

	var record models.OutfitRecord
	if err := h.db.Where("user_id = ? AND date = ?", userID, date).First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "No outfit recorded for this date"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch record"})
		return
	}

	c.JSON(http.StatusOK, record)
}

// Create creates a new outfit record (or updates existing for same date)
func (h *OutfitHandler) Create(c *gin.Context) {
	var req models.CreateOutfitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.Query("user_id")
	if userID == "" {
		userID = "demo-user"
	}

	// Check if record exists for this date
	var existing models.OutfitRecord
	if err := h.db.Where("user_id = ? AND date = ?", userID, req.Date).First(&existing).Error; err == nil {
		// Update existing record
		existing.Items = req.Items
		existing.CollageURL = req.CollageURL
		if err := h.db.Save(&existing).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
			return
		}
		c.JSON(http.StatusOK, existing)
		return
	}

	// Create new record
	record := models.OutfitRecord{
		UserID:     userID,
		Date:       req.Date,
		Items:      req.Items,
		CollageURL: req.CollageURL,
	}

	if err := h.db.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create record"})
		return
	}

	c.JSON(http.StatusCreated, record)
}

// Update updates an existing outfit record
func (h *OutfitHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var record models.OutfitRecord
	if err := h.db.First(&record, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch record"})
		return
	}

	var req models.UpdateOutfitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Items != nil {
		record.Items = req.Items
	}
	if req.CollageURL != nil {
		record.CollageURL = req.CollageURL
	}

	if err := h.db.Save(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
		return
	}

	c.JSON(http.StatusOK, record)
}

// Delete removes an outfit record
func (h *OutfitHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	result := h.db.Delete(&models.OutfitRecord{}, "id = ?", id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete record"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Record deleted"})
}
