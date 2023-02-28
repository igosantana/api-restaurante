package utils

import (
	"api-restaurante/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func GetErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "max":
		return "Should be max " + fe.Param()
	case "email":
		return "This field should be email format"
	}
	return "Unknown error"
}

func Ok(c *gin.Context, status int, message string, data interface{}) {
	c.AbortWithStatusJSON(http.StatusOK, models.Response{
		Data:    data,
		Status:  status,
		Message: message,
	})
}

func BadRequest(c *gin.Context, status int, message string, errors []models.ErrorDetail) {
	c.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
		Error:   errors,
		Status:  status,
		Message: message,
	})
}
