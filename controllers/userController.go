package controllers

import (
	"api-restaurante/initializers"
	"api-restaurante/models"
	"api-restaurante/utils"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type toCreateUser struct {
	Name     string `json:"name" binding:"max=110"`
	Email    string `json:"email" binding:"max=110,email"`
	Password string `json:"password"`
}

func Signup(c *gin.Context) {
	var user models.User
	var toCreateUser toCreateUser
	err := c.ShouldBind(&toCreateUser)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]utils.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = utils.ErrorMsg{Field: fe.Field(), Message: utils.GetErrorMsg(fe)}
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
		}
		return
	}
	result := initializers.DB.First(&user, "email = ?", toCreateUser.Email)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email already exists.",
		})
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(toCreateUser.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}
	user.Name = toCreateUser.Name
	user.Email = toCreateUser.Email
	user.Password = string(hash)
	initializers.DB.Create(&user)
	c.JSON(http.StatusOK, gin.H{
		"message": "User created!",
	})
}

func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := c.ShouldBind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to login",
		})
		return
	}
	var user models.User
	result := initializers.DB.First(&user, "email = ?", body.Email)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid email or password",
		})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	secretKey := os.Getenv("SECRET_KEY")
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create  token",
		})
	}
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
}

type userUpdate struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (u *userUpdate) toModel() *models.User {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		log.Println("Decrypt error update")
	}
	return &models.User{
		Name:     u.Name,
		Email:    u.Email,
		Password: string(hash),
	}
}

func UpdateUser(c *gin.Context) {
	var user models.User
	var userUpdate userUpdate
	id := c.Param("id")
	result := initializers.DB.First(&user, "id = ?", id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}
	err := c.ShouldBind(&userUpdate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fail to read body",
		})
	}
	initializers.DB.Model(&user).Updates(userUpdate.toModel())
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func DeleteUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	result := initializers.DB.First(&user, "id = ?", id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
		})
		return
	}
	initializers.DB.Delete(&user)
}

func GetUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	result := initializers.DB.First(&user, "id = ?", id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": user.UserToUser(),
	})
}
