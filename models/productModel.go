package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string    `gorm:"type:unique;type:varchar(110)"`
	Description string    `gorm:"type:varchar(200)"`
	Price       float64   `gorm:"type:float"`
	Category    string    `gorm:"type:varchar(100)"`
	Quantity    int       `gorm:"type:int"`
	Image       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
