package models

import (
	"mime/multipart"
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

type CreateProductForm struct {
	Name        string                `form:"name" binding:"required,max=110"`
	Description string                `form:"description" binding:"required,max=255"`
	Price       string                `form:"price" binding:"required"`
	Category    string                `form:"category" binding:"required,max=110"`
	Quantity    string                `form:"quantity" binding:"required"`
	File        *multipart.FileHeader `form:"file" binding:"required"`
}

type UpdateProduct struct {
	Name        string  `json:"name,omitempty" binding:"max=110"`
	Description string  `json:"description,omitempty" binding:"max=250"`
	Price       float64 `json:"price,omitempty"`
	Category    string  `json:"category,omitempty" binding:"max=110"`
	Quantity    int     `json:"quantity,omitempty"`
}

func (p *UpdateProduct) ToUpdateProductModel() *Product {
	return &Product{
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Category:    p.Category,
		Quantity:    p.Quantity,
	}
}
