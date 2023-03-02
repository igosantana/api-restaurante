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
	Name        string                `json:"name" form:"name" binding:"required,max=110"`
	Description string                `json:"description" form:"description" binding:"required,max=255"`
	Price       string                `json:"price" form:"price" binding:"required"`
	Category    string                `json:"category" form:"category" binding:"required,max=110"`
	Quantity    string                `json:"quantity" form:"quantity" binding:"required"`
	File        *multipart.FileHeader `json:"file" form:"file" binding:"required"`
}

type UpdateProduct struct {
	Name        string  `json:"name,omitempty" binding:"max=110"`
	Description string  `json:"description,omitempty" binding:"max=250"`
	Price       float64 `json:"price,omitempty"`
	Category    string  `json:"category,omitempty" binding:"max=110"`
	Quantity    int     `json:"quantity,omitempty"`
}

type ToGetAllProducts struct {
	ID          string  `json:"id"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Category    string  `json:"category,omitempty"`
	Quantity    int     `json:"quantity,omitempty"`
	ImageUrl    string  `json:"imageUrl,omitempty"`
}

func (up *UpdateProduct) ToUpdateProductModel() *Product {
	return &Product{
		Name:        up.Name,
		Description: up.Description,
		Price:       up.Price,
		Category:    up.Category,
		Quantity:    up.Quantity,
	}
}

func (p *Product) ToGetAll() ToGetAllProducts {
	return ToGetAllProducts{
		ID:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Category:    p.Category,
		Quantity:    p.Quantity,
		ImageUrl:    p.Image,
	}
}
