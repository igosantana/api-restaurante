package controllers

import (
	"api-restaurante/data/response"
	"api-restaurante/initializers"
	"api-restaurante/models"
	"errors"
	"fmt"
	"net/http"
	"os"

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

// CreateUser  godoc
// @Summary Create users
// @Description Save users data in DB
// @Param User body models.ToCreateUser true "Create user"
// @Tags User
// @Produce application/json
// @Success 200 {object} response.Response{}
// @Router /user [post]
func (uc *UserController) CreateUser(c *gin.Context) {
	var user models.User
	var toCreateUser models.ToCreateUser
	err := c.BindJSON(&toCreateUser)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]response.ErrorDetail, len(ve))
			for i, fe := range ve {
				out[i] = response.ErrorDetail{Field: fe.Field(), ErrorMessage: response.GetErrorMsg(fe), ErrorType: response.ErrorTypeValidation}
			}
			response.BadRequest(c, http.StatusBadRequest, "validation error", out)
		}
		return
	}
	result := initializers.DB.First(&user, "email = ?", toCreateUser.Email)
	if result.RowsAffected > 0 {
		var errors []response.ErrorDetail = make([]response.ErrorDetail, 0, 1)
		errors = append(errors, response.ErrorDetail{
			ErrorType:    response.ErrorTypeValidation,
			ErrorMessage: fmt.Sprintf("email already exists"),
		})
		response.BadRequest(c, http.StatusBadRequest, "validation error", errors)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(toCreateUser.Password), 10)
	if err != nil {
		var errors []response.ErrorDetail = make([]response.ErrorDetail, 0, 1)
		errors = append(errors, response.ErrorDetail{
			ErrorType:    response.ErrorTypeError,
			ErrorMessage: fmt.Sprintf("%v", err),
		})
		response.BadRequest(c, http.StatusBadRequest, "failed to hash password", errors)
		return
	}
	admin := os.Getenv("ADMIN")
	owner := os.Getenv("OWNER")
	if toCreateUser.Email == admin {
		user.Roles = append(user.Roles, Admin)
	} else if toCreateUser.Email == owner {
		user.Roles = append(user.Roles, Owner)
	} else {
		user.Roles = append(user.Roles, Customer)
	}
	user.Name = toCreateUser.Name
	user.Email = toCreateUser.Email
	user.Password = string(hash)
	create := initializers.DB.Create(&user)
	if create.Error != nil {
		var errors []response.ErrorDetail = make([]response.ErrorDetail, 0, 1)
		errors = append(errors, response.ErrorDetail{
			ErrorType:    response.ErrorTypeError,
			ErrorMessage: fmt.Sprintf("%v", create.Error),
		})
		response.BadRequest(c, http.StatusBadRequest, "failed to create user", errors)
		return
	}
	response.Ok(c, http.StatusOK, "success", nil)
}

// UpdateUser  godoc
// @Summary Update users
// @Description update user and save in DB. Only Admin or Owner.
// @Tags User
// @Param id path string true "update user by id"
// @Param user body models.UserUpdate true "update user"
// @Security cookieAuth
// @Produce application/json
// @Success 200 {object} response.Response{}
// @Router /user/{id} [patch]
func (uc *UserController) UpdateUser(c *gin.Context) {
	var user models.User
	var userUpdate models.UserUpdate
	id := c.Param("id")
	result := initializers.DB.First(&user, "id = ?", id)
	if result.RowsAffected == 0 {
		var errors []response.ErrorDetail = make([]response.ErrorDetail, 0, 1)
		errors = append(errors, response.ErrorDetail{
			ErrorType:    response.ErrorTypeError,
			ErrorMessage: fmt.Sprintf("user not found"),
		})
		response.BadRequest(c, http.StatusNotFound, "user not found", errors)
		return
	}
	err := c.BindJSON(&userUpdate)
	if err != nil {
		var errors []response.ErrorDetail = make([]response.ErrorDetail, 0, 1)
		errors = append(errors, response.ErrorDetail{
			ErrorType:    response.ErrorTypeError,
			ErrorMessage: fmt.Sprintf("%v", err),
		})
		response.BadRequest(c, http.StatusBadRequest, "error to read body", errors)
		return
	}
	update := initializers.DB.Model(&user).Updates(userUpdate.ToUpdateUserModel())
	if update.Error != nil {
		var errors []response.ErrorDetail = make([]response.ErrorDetail, 0, 1)
		errors = append(errors, response.ErrorDetail{
			ErrorType:    response.ErrorTypeError,
			ErrorMessage: fmt.Sprintf("%v", update.Error),
		})
		response.BadRequest(c, http.StatusBadRequest, "failed to update user", errors)
		return
	}
	response.Ok(c, http.StatusOK, "success", nil)
}

// DeleteUser  godoc
// @Summary Delete users
// @Description Remove users data by id. Only Admin or Owner.
// @Tags User
// @Param id path string true "remove user by id"
// @Security cookieAuth
// @Produce application/json
// @Success 200 {object} response.Response{}
// @Router /user/{id} [delete]
func (uc *UserController) DeleteUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	result := initializers.DB.First(&user, "id = ?", id)
	if result.RowsAffected == 0 {
		var errors []response.ErrorDetail = make([]response.ErrorDetail, 0, 1)
		errors = append(errors, response.ErrorDetail{
			ErrorType:    response.ErrorTypeError,
			ErrorMessage: fmt.Sprintf("user not found"),
		})
		response.BadRequest(c, http.StatusNotFound, "user not found", errors)
		return
	}
	delete := initializers.DB.Delete(&user)
	if delete.Error != nil {
		var errors []response.ErrorDetail = make([]response.ErrorDetail, 0, 1)
		errors = append(errors, response.ErrorDetail{
			ErrorType:    response.ErrorTypeError,
			ErrorMessage: fmt.Sprintf("%v", delete.Error),
		})
		response.BadRequest(c, http.StatusBadRequest, "failed to delete user", errors)
		return
	}
	response.Ok(c, http.StatusOK, "success", nil)
}

// GetOneUser  godoc
// @Summary Get one user
// @Description Get user by id. Only Admin or Owner.
// @Tags User
// @Param id path string true "Get user by id"
// @Security cookieAuth
// @Produce application/json
// @Success 200 {object} response.Response{}
// @Router /user/{id} [get]
func (uc *UserController) GetUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	result := initializers.DB.First(&user, "id = ?", id)
	if result.RowsAffected == 0 {
		var errors []response.ErrorDetail = make([]response.ErrorDetail, 0, 1)
		errors = append(errors, response.ErrorDetail{
			ErrorType:    response.ErrorTypeError,
			ErrorMessage: fmt.Sprintf("user not found"),
		})
		response.BadRequest(c, http.StatusNotFound, "user not found", errors)
		return
	}
	response.Ok(c, http.StatusOK, "success", user.UserToUser())
}

// GetAllUsers  godoc
// @Summary Get all users
// @Description Get all users. Only Admin.
// @Tags User
// @Security cookieAuth
// @Produce application/json
// @Success 200 {object} response.Response{}
// @Router /user [get]
func (uc *UserController) GetAllUsers(c *gin.Context) {
	var users []models.User
	result := initializers.DB.Find(&users)
	if result.RowsAffected == 0 {
		var errors []response.ErrorDetail = make([]response.ErrorDetail, 0, 1)
		errors = append(errors, response.ErrorDetail{
			ErrorType:    response.ErrorTypeError,
			ErrorMessage: fmt.Sprintf("user not found"),
		})
		response.BadRequest(c, http.StatusNotFound, "user not found", errors)
		return
	}
	response.Ok(c, http.StatusOK, "success", &users)
}
