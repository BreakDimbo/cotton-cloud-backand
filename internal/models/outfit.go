package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OutfitRecord represents a logged outfit for a specific date
type OutfitRecord struct {
	ID         string     `json:"id" gorm:"primaryKey"`
	UserID     string     `json:"userId" gorm:"index"`
	Date       string     `json:"date" gorm:"index"`      // YYYY-MM-DD format
	Items      StringList `json:"items" gorm:"type:text"` // JSON array of ClothingItem IDs
	CollageURL *string    `json:"collageUrl,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

func (o *OutfitRecord) BeforeCreate(tx *gorm.DB) error {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
	return nil
}

// CreateOutfitRequest is the request body for creating an outfit record
type CreateOutfitRequest struct {
	Date       string   `json:"date" binding:"required"`
	Items      []string `json:"items" binding:"required"`
	CollageURL *string  `json:"collageUrl,omitempty"`
}

// UpdateOutfitRequest is the request body for updating an outfit record
type UpdateOutfitRequest struct {
	Items      []string `json:"items,omitempty"`
	CollageURL *string  `json:"collageUrl,omitempty"`
}
