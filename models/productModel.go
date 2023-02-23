package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string    `gorm:"type:unique;type:varchar(110)" binding:"required"`
	Description string    `gorm:"type:varchar(200)" binding:"required"`
	Price       float64   `gorm:"type:decimal(7,6)" binding:"required"`
	Category    string    `gorm:"type:varchar(100)" binding:"required"`
	Quantity    int       `gorm:"type:int"`
	Image       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
