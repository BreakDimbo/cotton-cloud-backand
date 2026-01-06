package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cotton-cloud-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// AIHandler handles AI-related proxy requests to Gemini
type AIHandler struct {
	gemini *services.GeminiService
}

// NewAIHandler creates a new AIHandler
func NewAIHandler() *AIHandler {
	gemini, err := services.NewGeminiService()
	if err != nil {
		// Log error but continue - AI features will return mock data
		println("Warning: Failed to initialize Gemini service:", err.Error())
		return &AIHandler{gemini: nil}
	}
	return &AIHandler{gemini: gemini}
}

// AnalyzeClothingRequest is the request body for clothing analysis
type AnalyzeClothingRequest struct {
	ImageBase64 string `json:"imageBase64" binding:"required"`
	MimeType    string `json:"mimeType" binding:"required"`
}

// RefineAnalysisRequest is the request body for refining analysis
type RefineAnalysisRequest struct {
	ImageBase64  string `json:"imageBase64" binding:"required"`
	MimeType     string `json:"mimeType" binding:"required"`
	UserFeedback string `json:"userFeedback" binding:"required"`
}

// GenerateCutoutRequest is the request body for cutout generation
type GenerateCutoutRequest struct {
	ImageBase64 string `json:"imageBase64" binding:"required"`
	MimeType    string `json:"mimeType" binding:"required"`
}

// GenerateAvatarRequest is the request body for avatar generation
type GenerateAvatarRequest struct {
	FaceImageBase64 string `json:"faceImageBase64" binding:"required"`
	MimeType        string `json:"mimeType" binding:"required"`
	Gender          string `json:"gender" binding:"required"`
	Height          string `json:"height" binding:"required"`
	Weight          string `json:"weight" binding:"required"`
	Bust            string `json:"bust"`
	Waist           string `json:"waist"`
	Hips            string `json:"hips"`
	Thigh           string `json:"thigh"`
	Calf            string `json:"calf"`
	Features        string `json:"features"`
}

// GenerateCollageRequest is the request body for collage generation
type GenerateCollageRequest struct {
	ItemImages []string `json:"itemImages" binding:"required"` // Base64 images
}

// VirtualTryOnRequest is the request body for virtual try-on
type VirtualTryOnRequest struct {
	AvatarImageBase64 string   `json:"avatarImageBase64" binding:"required"`
	ItemImages        []string `json:"itemImages" binding:"required"` // Base64 images
}

// AnalyzeClothing analyzes a clothing image using Gemini AI
func (h *AIHandler) AnalyzeClothing(c *gin.Context) {
	var req AnalyzeClothingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("[HANDLER] AnalyzeClothing MIME: %s\n", req.MimeType)

	// If Gemini service not available, return mock data
	if h.gemini == nil {
		c.JSON(http.StatusOK, gin.H{
			"category":    "Tops",
			"color":       "White",
			"material":    "Cotton",
			"description": "A soft, cloudlike piece perfect for everyday elegance.",
			"tags":        []string{"casual", "everyday", "basic"},
			"style":       []string{"Casual", "Minimalist"},
			"season":      []string{"Spring", "Summer", "All Season"},
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	analysis, err := h.gemini.AnalyzeClothing(ctx, req.ImageBase64, req.MimeType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analysis)
}

// RefineAnalysis refines clothing analysis based on user feedback
func (h *AIHandler) RefineAnalysis(c *gin.Context) {
	var req RefineAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.gemini == nil {
		c.JSON(http.StatusOK, gin.H{
			"category":    "Tops",
			"color":       "White",
			"material":    "Cotton",
			"description": "A refined piece based on your feedback.",
			"tags":        []string{"refined", "custom"},
			"style":       []string{"Casual"},
			"season":      []string{"All Season"},
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	analysis, err := h.gemini.RefineClothingAnalysis(ctx, req.ImageBase64, req.UserFeedback, req.MimeType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analysis)
}

// GenerateCutout generates a clothing cutout using Gemini AI
func (h *AIHandler) GenerateCutout(c *gin.Context) {
	var req GenerateCutoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("[HANDLER] GenerateCutout MIME: %s\n", req.MimeType)

	if h.gemini == nil {
		c.JSON(http.StatusOK, gin.H{
			"message":  "Cutout generation - Gemini not configured",
			"imageUrl": "https://picsum.photos/400/600",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
	defer cancel()

	imageBase64, err := h.gemini.GenerateCutout(ctx, req.ImageBase64, req.MimeType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"imageBase64": imageBase64,
		"message":     "Cutout generated successfully",
	})
}

// GenerateAvatar generates a full-body avatar using Gemini AI
func (h *AIHandler) GenerateAvatar(c *gin.Context) {
	var req GenerateAvatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.gemini == nil {
		c.JSON(http.StatusOK, gin.H{
			"message":  "Avatar generation - Gemini not configured",
			"imageUrl": "https://picsum.photos/400/600",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 90*time.Second)
	defer cancel()

	// Convert request to AvatarMetrics
	metrics := services.AvatarMetrics{
		Gender:   req.Gender,
		Height:   req.Height,
		Weight:   req.Weight,
		Bust:     req.Bust,
		Waist:    req.Waist,
		Hips:     req.Hips,
		Thigh:    req.Thigh,
		Calf:     req.Calf,
		Features: req.Features,
	}

	imageBase64, err := h.gemini.GenerateAvatar(ctx, req.FaceImageBase64, req.MimeType, metrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"imageBase64": imageBase64,
		"message":     "Avatar generated successfully",
	})
}

// GenerateCollage generates an outfit collage using Gemini AI
func (h *AIHandler) GenerateCollage(c *gin.Context) {
	var req GenerateCollageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.gemini == nil {
		c.JSON(http.StatusOK, gin.H{
			"message":    "Collage generation - Gemini not configured",
			"collageUrl": "https://picsum.photos/600/800",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
	defer cancel()

	imageBase64, err := h.gemini.GenerateCollage(ctx, req.ItemImages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"imageBase64": imageBase64,
		"message":     "Collage generated successfully",
	})
}

// VirtualTryOn performs virtual try-on using Gemini AI
func (h *AIHandler) VirtualTryOn(c *gin.Context) {
	var req VirtualTryOnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.gemini == nil {
		c.JSON(http.StatusOK, gin.H{
			"message":  "Virtual try-on - Gemini not configured",
			"tryOnUrl": "https://picsum.photos/400/600",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 90*time.Second)
	defer cancel()

	imageBase64, err := h.gemini.VirtualTryOn(ctx, req.AvatarImageBase64, req.ItemImages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"imageBase64": imageBase64,
		"message":     "Virtual try-on generated successfully",
	})
}
