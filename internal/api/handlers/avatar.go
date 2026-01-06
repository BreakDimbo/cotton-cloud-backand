package handlers

import (
	"net/http"

	"cotton-cloud-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AvatarHandler handles avatar-related requests
type AvatarHandler struct {
	db *gorm.DB
}

// NewAvatarHandler creates a new AvatarHandler
func NewAvatarHandler(db *gorm.DB) *AvatarHandler {
	return &AvatarHandler{db: db}
}

// List returns all avatars for the current user
func (h *AvatarHandler) List(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		userID = "demo-user"
	}

	var avatars []models.AvatarProfile
	if err := h.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&avatars).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch avatars"})
		return
	}

	c.JSON(http.StatusOK, avatars)
}

// Get returns a single avatar by ID
func (h *AvatarHandler) Get(c *gin.Context) {
	id := c.Param("id")

	var avatar models.AvatarProfile
	if err := h.db.First(&avatar, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Avatar not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch avatar"})
		return
	}

	c.JSON(http.StatusOK, avatar)
}

// Create creates a new avatar
func (h *AvatarHandler) Create(c *gin.Context) {
	var req models.CreateAvatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.Query("user_id")
	if userID == "" {
		userID = "demo-user"
	}

	avatar := models.AvatarProfile{
		UserID:   userID,
		Name:     req.Name,
		Tag:      req.Tag,
		ImageURL: req.ImageURL,
	}
	avatar.SetMetrics(req.Metrics)

	// Check if this is the first avatar (make it active)
	var count int64
	h.db.Model(&models.AvatarProfile{}).Where("user_id = ?", userID).Count(&count)
	if count == 0 {
		avatar.IsActive = true
	}

	if err := h.db.Create(&avatar).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create avatar"})
		return
	}

	c.JSON(http.StatusCreated, avatar)
}

// Update updates an existing avatar
func (h *AvatarHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var avatar models.AvatarProfile
	if err := h.db.First(&avatar, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Avatar not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch avatar"})
		return
	}

	var req models.UpdateAvatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != nil {
		avatar.Name = *req.Name
	}
	if req.Tag != nil {
		avatar.Tag = *req.Tag
	}
	if req.ImageURL != nil {
		avatar.ImageURL = *req.ImageURL
	}
	if req.Metrics != nil {
		avatar.SetMetrics(*req.Metrics)
	}

	if err := h.db.Save(&avatar).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update avatar"})
		return
	}

	c.JSON(http.StatusOK, avatar)
}

// Delete removes an avatar
func (h *AvatarHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	result := h.db.Delete(&models.AvatarProfile{}, "id = ?", id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete avatar"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Avatar not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Avatar deleted"})
}

// Activate sets an avatar as the active avatar for the user
func (h *AvatarHandler) Activate(c *gin.Context) {
	id := c.Param("id")

	var avatar models.AvatarProfile
	if err := h.db.First(&avatar, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Avatar not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch avatar"})
		return
	}

	// Deactivate all other avatars for this user
	h.db.Model(&models.AvatarProfile{}).Where("user_id = ?", avatar.UserID).Update("is_active", false)

	// Activate this avatar
	avatar.IsActive = true
	if err := h.db.Save(&avatar).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate avatar"})
		return
	}

	c.JSON(http.StatusOK, avatar)
}
