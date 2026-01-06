package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user account
type User struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"uniqueIndex"`
	Nickname  string    `json:"nickname"`
	Password  string    `json:"-"` // Never expose password in JSON
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// Relationships
	ClothingItems []ClothingItem  `json:"clothingItems,omitempty" gorm:"foreignKey:UserID"`
	Avatars       []AvatarProfile `json:"avatars,omitempty" gorm:"foreignKey:UserID"`
	OutfitRecords []OutfitRecord  `json:"outfitRecords,omitempty" gorm:"foreignKey:UserID"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}
