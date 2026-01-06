package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AvatarMetrics represents the physical measurements for an avatar
type AvatarMetrics struct {
	Gender   string `json:"gender"`
	Height   string `json:"height"`
	Weight   string `json:"weight"`
	Bust     string `json:"bust"`
	Waist    string `json:"waist"`
	Hips     string `json:"hips"`
	Thigh    string `json:"thigh"`
	Calf     string `json:"calf"`
	Features string `json:"features"`
}

// AvatarProfile represents a user's digital avatar
type AvatarProfile struct {
	ID       string `json:"id" gorm:"primaryKey"`
	UserID   string `json:"userId" gorm:"index"`
	Name     string `json:"name"`
	Tag      string `json:"tag"`
	ImageURL string `json:"imageUrl"`
	IsActive bool   `json:"isActive" gorm:"default:false"`

	// Metrics stored as JSON
	MetricsGender   string `json:"metricsGender"`
	MetricsHeight   string `json:"metricsHeight"`
	MetricsWeight   string `json:"metricsWeight"`
	MetricsBust     string `json:"metricsBust"`
	MetricsWaist    string `json:"metricsWaist"`
	MetricsHips     string `json:"metricsHips"`
	MetricsThigh    string `json:"metricsThigh"`
	MetricsCalf     string `json:"metricsCalf"`
	MetricsFeatures string `json:"metricsFeatures"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (a *AvatarProfile) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

// GetMetrics returns the avatar metrics as a struct
func (a *AvatarProfile) GetMetrics() AvatarMetrics {
	return AvatarMetrics{
		Gender:   a.MetricsGender,
		Height:   a.MetricsHeight,
		Weight:   a.MetricsWeight,
		Bust:     a.MetricsBust,
		Waist:    a.MetricsWaist,
		Hips:     a.MetricsHips,
		Thigh:    a.MetricsThigh,
		Calf:     a.MetricsCalf,
		Features: a.MetricsFeatures,
	}
}

// SetMetrics sets the avatar metrics from a struct
func (a *AvatarProfile) SetMetrics(m AvatarMetrics) {
	a.MetricsGender = m.Gender
	a.MetricsHeight = m.Height
	a.MetricsWeight = m.Weight
	a.MetricsBust = m.Bust
	a.MetricsWaist = m.Waist
	a.MetricsHips = m.Hips
	a.MetricsThigh = m.Thigh
	a.MetricsCalf = m.Calf
	a.MetricsFeatures = m.Features
}

// CreateAvatarRequest is the request body for creating an avatar
type CreateAvatarRequest struct {
	Name     string        `json:"name" binding:"required"`
	Tag      string        `json:"tag"`
	ImageURL string        `json:"imageUrl" binding:"required"`
	Metrics  AvatarMetrics `json:"metrics"`
}

// UpdateAvatarRequest is the request body for updating an avatar
type UpdateAvatarRequest struct {
	Name     *string        `json:"name,omitempty"`
	Tag      *string        `json:"tag,omitempty"`
	ImageURL *string        `json:"imageUrl,omitempty"`
	Metrics  *AvatarMetrics `json:"metrics,omitempty"`
}
