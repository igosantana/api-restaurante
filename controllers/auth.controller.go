package controllers

import (
	"api-restaurante/initializers"
	"api-restaurante/models"
	"api-restaurante/token"
	"api-restaurante/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(DB *gorm.DB) AuthController {
	return AuthController{DB}
}

func (ac *AuthController) Login(c *gin.Context) {
	var userLogin models.UserLogin
	if err := c.BindJSON(&userLogin); err != nil {
		var errors []models.ErrorDetail = make([]models.ErrorDetail, 0, 1)
		errors = append(errors, models.ErrorDetail{
			ErrorType:    models.ErrorTypeValidation,
			ErrorMessage: fmt.Sprintf("%v", err),
		})
		utils.BadRequest(c, http.StatusBadRequest, "Invalid request", errors)
	}
	var user models.User
	result := initializers.DB.First(&user, "email = ?", userLogin.Email)
	if result.RowsAffected == 0 {
		var errors []models.ErrorDetail = make([]models.ErrorDetail, 0, 1)
		errors = append(errors, models.ErrorDetail{
			ErrorType:    models.ErrorTypeValidation,
			ErrorMessage: fmt.Sprintf("%v", "Invalid email or password"),
		})
		utils.BadRequest(c, http.StatusBadRequest, "Invalid email or password", errors)
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password))
	if err != nil {
		var errors []models.ErrorDetail = make([]models.ErrorDetail, 0, 1)
		errors = append(errors, models.ErrorDetail{
			ErrorType:    models.ErrorTypeValidation,
			ErrorMessage: fmt.Sprintf("%v", err),
		})
		utils.BadRequest(c, http.StatusBadRequest, "Invalid email or password", errors)
	}
	var claims = &models.JwtClaims{}
	claims.UserId = user.ID.String()
	claims.Roles = user.Roles
	var tokenCreationTime = time.Now().UTC()
	var expirationTime = tokenCreationTime.Add(time.Duration(2) * time.Hour)
	tokenString, err := token.GenerateToken(claims, expirationTime)
	if err != nil {
		utils.BadRequest(c, http.StatusBadRequest, "Error in generating token", []models.ErrorDetail{
			{
				ErrorType:    models.ErrorTypeError,
				ErrorMessage: err.Error(),
			},
		})
	}
	c.SetCookie("Authorization", tokenString, 3600*2, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (ac *AuthController) LogoutUser(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
