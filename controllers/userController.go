package controllers

import (
	"api-restaurante/initializers"
	"api-restaurante/models"
	"api-restaurante/token"
	"api-restaurante/utils"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

const (
	Admin    = "Admin"
	Owner    = "Owner"
	Customer = "Customer"
)

func CreateUser(c *gin.Context) {
	var user models.User
	var toCreateUser models.ToCreateUser
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
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(toCreateUser.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}
	if toCreateUser.Email == "igo@admin.com" {
		user.Roles = append(user.Roles, Admin)
	} else if toCreateUser.Email == "leo@email.com" {
		user.Roles = append(user.Roles, Owner)
	} else {
		user.Roles = append(user.Roles, Customer)
	}
	user.Name = toCreateUser.Name
	user.Email = toCreateUser.Email
	user.Password = string(hash)
	initializers.DB.Create(&user)
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created!",
	})
}

func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBind(&body); err != nil {
		var errors []models.ErrorDetail = make([]models.ErrorDetail, 0, 1)
		errors = append(errors, models.ErrorDetail{
			ErrorType:    models.ErrorTypeValidation,
			ErrorMessage: fmt.Sprintf("%v", err),
		})
		badRequest(c, http.StatusBadRequest, "invalid request", errors)
	}
	var user models.User
	result := initializers.DB.First(&user, "email = ?", body.Email)
	if result.RowsAffected == 0 {
		var errors []models.ErrorDetail = make([]models.ErrorDetail, 0, 1)
		errors = append(errors, models.ErrorDetail{
			ErrorType:    models.ErrorTypeValidation,
			ErrorMessage: fmt.Sprintf("%v", "Invalid email or password"),
		})
		badRequest(c, http.StatusBadRequest, "Invalid email or password", errors)
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		var errors []models.ErrorDetail = make([]models.ErrorDetail, 0, 1)
		errors = append(errors, models.ErrorDetail{
			ErrorType:    models.ErrorTypeValidation,
			ErrorMessage: fmt.Sprintf("%v", err),
		})
		badRequest(c, http.StatusBadRequest, "Invalid email or password", errors)
	}
	var claims = &models.JwtClaims{}
	claims.UserId = user.ID.String()
	claims.Roles = user.Roles
	var tokenCreationTime = time.Now().UTC()
	var expirationTime = tokenCreationTime.Add(time.Duration(2) * time.Hour)
	tokenString, err := token.GenerateToken(claims, expirationTime)
	if err != nil {
		badRequest(c, http.StatusBadRequest, "Error in generating token", []models.ErrorDetail{
			{
				ErrorType:    models.ErrorTypeError,
				ErrorMessage: err.Error(),
			},
		})
	}
	c.SetCookie("Authorization", tokenString, 3600*2, "", "", false, true)
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
