package handlers

import (
	"net/http"

	"cotton-cloud-backend/internal/models"
	"cotton-cloud-backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	db   *gorm.DB
	auth *services.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		db:   db,
		auth: services.NewAuthService(),
	}
}

// RegisterRequest is the request body for user registration
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname"`
}

// LoginRequest is the request body for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse is the response for successful authentication
type AuthResponse struct {
	Token   string `json:"token"`
	UserID  string `json:"userId"`
	Email   string `json:"email"`
	Message string `json:"message,omitempty"`
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email already exists
	var existingUser models.User
	if err := h.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Hash password
	hashedPassword, err := h.auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	// Create user
	user := models.User{
		Email:    req.Email,
		Password: hashedPassword,
		Nickname: req.Nickname,
	}

	if user.Nickname == "" {
		user.Nickname = "Fashion Lover"
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate token
	token, err := h.auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		Token:   token,
		UserID:  user.ID,
		Email:   user.Email,
		Message: "Registration successful",
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	var user models.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check password
	if !h.auth.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate token
	token, err := h.auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token:  token,
		UserID: user.ID,
		Email:  user.Email,
	})
}

// RefreshToken refreshes an existing valid token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Get token from header
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
		return
	}

	tokenString := authHeader[7:]

	// Refresh token
	newToken, err := h.auth.RefreshToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": newToken})
}
