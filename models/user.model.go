package models

import (
	"api-restaurante/app"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string         `gorm:"type:varchar(110)"`
	Email     string         `gorm:"uniqueIndex;type:varchar(110)"`
	Password  string         `gorm:"type:varchar(110)"`
	Roles     pq.StringArray `gorm:"type:text[]"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ToCreateUser struct {
	Name     string `json:"name" binding:"max=110,required,min=3"`
	Email    string `json:"email" binding:"max=110,email,required"`
	Password string `json:"password" binding:"min=6,required,max=8"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"email,required,max=110"`
	Password string `json:"password" binding:"required,min=6,max=8"`
}

type UserUpdate struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (u *UserUpdate) ToUpdateUserModel() *User {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		log.Println("Decrypt error update")
	}
	return &User{
		Name:     u.Name,
		Email:    u.Email,
		Password: string(hash),
	}
}

func (u *User) UserToUser() app.User {
	return app.User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}
