package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ClothingItem represents a piece of clothing in a user's wardrobe
type ClothingItem struct {
	ID                string     `json:"id" gorm:"primaryKey"`
	UserID            string     `json:"userId" gorm:"index"`
	ImageURL          string     `json:"imageUrl"`
	OriginalImageURL  *string    `json:"originalImageUrl,omitempty"`
	ProcessedImageURL *string    `json:"processedImageUrl,omitempty"`
	Category          string     `json:"category"`
	Color             string     `json:"color"`
	Material          *string    `json:"material,omitempty"`
	Description       *string    `json:"description,omitempty"`
	Tags              StringList `json:"tags" gorm:"type:text"`
	Style             StringList `json:"style" gorm:"type:text"`
	Season            StringList `json:"season" gorm:"type:text"`
	WearCount         int        `json:"wearCount" gorm:"default:0"`
	MaxWearCount      int        `json:"maxWearCount" gorm:"default:5"`
	LastWashedAt      *time.Time `json:"lastWashedAt,omitempty"`
	CreatedAt         time.Time  `json:"addedAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}

func (c *ClothingItem) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

// NeedsCare returns true if the item needs to be washed
func (c *ClothingItem) NeedsCare() bool {
	return c.WearCount >= c.MaxWearCount
}

// CreateClothingItemRequest is the request body for creating a clothing item
type CreateClothingItemRequest struct {
	ImageURL          string   `json:"imageUrl" binding:"required"`
	OriginalImageURL  *string  `json:"originalImageUrl,omitempty"`
	ProcessedImageURL *string  `json:"processedImageUrl,omitempty"`
	Category          string   `json:"category" binding:"required"`
	Color             string   `json:"color" binding:"required"`
	Material          *string  `json:"material,omitempty"`
	Description       *string  `json:"description,omitempty"`
	Tags              []string `json:"tags,omitempty"`
	Style             []string `json:"style,omitempty"`
	Season            []string `json:"season,omitempty"`
	MaxWearCount      *int     `json:"maxWearCount,omitempty"`
}

// UpdateClothingItemRequest is the request body for updating a clothing item
type UpdateClothingItemRequest struct {
	ImageURL     *string  `json:"imageUrl,omitempty"`
	Category     *string  `json:"category,omitempty"`
	Color        *string  `json:"color,omitempty"`
	Material     *string  `json:"material,omitempty"`
	Description  *string  `json:"description,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	Style        []string `json:"style,omitempty"`
	Season       []string `json:"season,omitempty"`
	MaxWearCount *int     `json:"maxWearCount,omitempty"`
}
