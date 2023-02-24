package models

import (
	"api-restaurante/app"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string    `gorm:"type:varchar(110)"`
	Email     string    `gorm:"unique;type:varchar(110)"`
	Password  string    `gorm:"type:varchar(110)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) UserToUser() app.User {
	return app.User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}
