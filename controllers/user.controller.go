package controllers

import (
	"api-restaurante/initializers"
	"api-restaurante/models"
	"api-restaurante/utils"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	Admin    = "Admin"
	Owner    = "Owner"
	Customer = "Customer"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(DB *gorm.DB) UserController {
	return UserController{DB}
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var user models.User
	var toCreateUser models.ToCreateUser
	err := c.BindJSON(&toCreateUser)
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
		var errors []models.ErrorDetail = make([]models.ErrorDetail, 0, 1)
		errors = append(errors, models.ErrorDetail{
			ErrorType:    models.ErrorTypeValidation,
			ErrorMessage: fmt.Sprintf("Email already exists"),
		})
		utils.BadRequest(c, http.StatusBadRequest, "Email already exists", errors)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(toCreateUser.Password), 10)
	if err != nil {
		var errors []models.ErrorDetail = make([]models.ErrorDetail, 0, 1)
		errors = append(errors, models.ErrorDetail{
			ErrorType:    models.ErrorTypeError,
			ErrorMessage: fmt.Sprintf("%v", err),
		})
		utils.BadRequest(c, http.StatusBadRequest, "Failed to hash password", errors)
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
	create := initializers.DB.Create(&user)
	if create.Error != nil {
		var errors []models.ErrorDetail = make([]models.ErrorDetail, 0, 1)
		errors = append(errors, models.ErrorDetail{
			ErrorType:    models.ErrorTypeError,
			ErrorMessage: fmt.Sprintf("%v", create.Error),
		})
		utils.BadRequest(c, http.StatusBadRequest, "Failed to create user", errors)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created!",
	})
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	var user models.User
	var userUpdate models.UserUpdate
	id := c.Param("id")
	result := initializers.DB.First(&user, "id = ?", id)
	if result.RowsAffected == 0 {
		var errors []models.ErrorDetail = make([]models.ErrorDetail, 0, 1)
		errors = append(errors, models.ErrorDetail{
			ErrorType:    models.ErrorTypeError,
			ErrorMessage: fmt.Sprintf("User not found"),
		})
		utils.BadRequest(c, http.StatusBadRequest, "User not found", errors)
	}
	err := c.BindJSON(&userUpdate)
	if err != nil {
		var errors []models.ErrorDetail = make([]models.ErrorDetail, 0, 1)
		errors = append(errors, models.ErrorDetail{
			ErrorType:    models.ErrorTypeError,
			ErrorMessage: fmt.Sprintf("%v", err),
		})
		utils.BadRequest(c, http.StatusBadRequest, "Error to read body", errors)
	}
	initializers.DB.Model(&user).Updates(userUpdate.ToUpdateUserModel())
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	result := initializers.DB.First(&user, "id = ?", id)
	if result.RowsAffected == 0 {
		var errors []models.ErrorDetail = make([]models.ErrorDetail, 0, 1)
		errors = append(errors, models.ErrorDetail{
			ErrorType:    models.ErrorTypeError,
			ErrorMessage: fmt.Sprintf("User not found"),
		})
		utils.BadRequest(c, http.StatusBadRequest, "User not found", errors)
	}
	initializers.DB.Delete(&user)
}

func (uc *UserController) GetUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	result := initializers.DB.First(&user, "id = ?", id)
	if result.RowsAffected == 0 {
		var errors []models.ErrorDetail = make([]models.ErrorDetail, 0, 1)
		errors = append(errors, models.ErrorDetail{
			ErrorType:    models.ErrorTypeError,
			ErrorMessage: fmt.Sprintf("User not found"),
		})
		utils.BadRequest(c, http.StatusBadRequest, "User not found", errors)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": user.UserToUser(),
	})
}

func (uc *UserController) GetAllUsers(c *gin.Context) {
	var users []models.User
	result := initializers.DB.Find(&users)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": &users,
	})
}
