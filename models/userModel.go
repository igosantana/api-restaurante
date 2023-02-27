package models

import (
	"api-restaurante/app"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string         `gorm:"type:varchar(110)"`
	Email     string         `gorm:"unique;type:varchar(110)"`
	Password  string         `gorm:"type:varchar(110)"`
	Roles     pq.StringArray `gorm:"type:text[]"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type ToCreateUser struct {
	Name     string `json:"name" binding:"max=110"`
	Email    string `json:"email" binding:"max=110,email"`
	Password string `json:"password"`
}

func (u *User) UserToUser() app.User {
	return app.User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}
