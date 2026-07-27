package main

import (
	"time"

	"gorm.io/gorm"
)

// Button tracks the click metric state in PostgreSQL
type Button struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ClickedCount int       `gorm:"default:0" json:"clicked_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// IncrementClick atomically increments the click count in the database
func IncrementClick(db *gorm.DB) (int, error) {
	var button Button

	// Fetch or create the single counter record (ID = 1)
	err := db.FirstOrCreate(&button, Button{ID: 1}).Error
	if err != nil {
		return 0, err
	}

	// Atomically increment the counter to avoid race conditions under load
	err = db.Model(&button).UpdateColumn("clicked_count", gorm.Expr("clicked_count + ?", 1)).Error
	if err != nil {
		return 0, err
	}

	// Reload updated count
	db.First(&button, 1)
	return button.ClickedCount, nil
}
