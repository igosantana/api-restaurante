package controllers

import (
	"api-restaurante/data/response"
	"api-restaurante/initializers"
	"api-restaurante/models"
	"api-restaurante/token"
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

// UserLogin	godoc
// @Summary Users Login
// @Description User Login. Set cookie with name `Authorization`. You need to include this cookie in subsequent requests.
// @Tags Auth
// @Produce application/json
// @Success 200 {object} response.Response{}
// @Router /auth/login [post]
func (ac *AuthController) Login(c *gin.Context) {
	var userLogin models.UserLogin
	if err := c.BindJSON(&userLogin); err != nil {
		var errors []response.ErrorDetail = make([]response.ErrorDetail, 0, 1)
		errors = append(errors, response.ErrorDetail{
			ErrorType:    response.ErrorTypeValidation,
			ErrorMessage: fmt.Sprintf("%v", err),
		})
		response.BadRequest(c, http.StatusBadRequest, "invalid request", errors)
		return
	}
	var user models.User
	result := initializers.DB.First(&user, "email = ?", userLogin.Email)
	if result.RowsAffected == 0 {
		var errors []response.ErrorDetail = make([]response.ErrorDetail, 0, 1)
		errors = append(errors, response.ErrorDetail{
			ErrorType:    response.ErrorTypeValidation,
			ErrorMessage: fmt.Sprintf("%v", "invalid email or password"),
		})
		response.BadRequest(c, http.StatusBadRequest, "credential error", errors)
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password))
	if err != nil {
		var errors []response.ErrorDetail = make([]response.ErrorDetail, 0, 1)
		errors = append(errors, response.ErrorDetail{
			ErrorType:    response.ErrorTypeValidation,
			ErrorMessage: fmt.Sprintf("%v", err),
		})
		response.BadRequest(c, http.StatusBadRequest, "invalid email or password", errors)
		return
	}
	var claims = &models.JwtClaims{}
	claims.UserId = user.ID.String()
	claims.Roles = user.Roles
	var tokenCreationTime = time.Now().UTC()
	var expirationTime = tokenCreationTime.Add(time.Duration(2) * time.Hour)
	tokenString, err := token.GenerateToken(claims, expirationTime)
	if err != nil {
		var errors []response.ErrorDetail = make([]response.ErrorDetail, 0, 1)
		errors = append(errors, response.ErrorDetail{
			ErrorType:    response.ErrorTypeError,
			ErrorMessage: fmt.Sprintf("%v", err),
		})
		response.BadRequest(c, http.StatusBadRequest, "failed to generate token", errors)
	}
	c.SetCookie("Authorization", tokenString, 3600*2, "", "", false, true)
	response.Ok(c, http.StatusOK, "success", nil)
}

// UserLogout	godoc
// @Summary Users Logout
// @Description Logout user.
// @Tags Auth
// @Security cookieAuth
// @Produce application/json
// @Success 200 {object} response.Response{}
// @Router /auth/logout [get]
func (ac *AuthController) LogoutUser(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "", "", false, true)
	response.Ok(c, http.StatusOK, "success", nil)
}
