package controllers

import (
	"api-restaurante/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ok(c *gin.Context, status int, message string, data interface{}) {
	c.AbortWithStatusJSON(http.StatusOK, models.Response{
		Data:    data,
		Status:  status,
		Message: message,
	})
}

func badRequest(c *gin.Context, status int, message string, errors []models.ErrorDetail) {
	c.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
		Error:   errors,
		Status:  status,
		Message: message,
	})
}
