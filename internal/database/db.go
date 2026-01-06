package database

import (
	"cotton-cloud-backend/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	// Use pure Go SQLite driver (no CGO required)
	_ "modernc.org/sqlite"
)

// InitDB initializes the database connection
func InitDB() (*gorm.DB, error) {
	// Use the pure Go SQLite driver via GORM
	db, err := gorm.Open(sqlite.Dialector{
		DriverName: "sqlite",
		DSN:        "cotton_cloud.db",
	}, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// AutoMigrate runs database migrations for all models
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.ClothingItem{},
		&models.AvatarProfile{},
		&models.OutfitRecord{},
	)
}
