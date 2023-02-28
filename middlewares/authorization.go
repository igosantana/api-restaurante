package middlewares

import (
	"api-restaurante/models"
	"api-restaurante/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReturnUnauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{
		Error: []models.ErrorDetail{
			{
				ErrorType:    models.ErrorTypeUnauthorized,
				ErrorMessage: "You are not authorized to access this path",
			},
		},
		Status:  http.StatusUnauthorized,
		Message: "Unauthorized access",
	})
}

func ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("Authorization")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "fail", "message": "You are not logged in",
			})
			return
		}
		valid, claims := token.VerifyToken(tokenString)
		if !valid {
			ReturnUnauthorized(c)
		}
		if len(c.Keys) == 0 {
			c.Keys = make(map[string]interface{})
		}
		c.Keys["userId"] = claims.UserId
		c.Keys["roles"] = claims.Roles
	}
}

func IsOwnerOrAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userId := c.Keys["userId"]
		userRoles := c.Keys["roles"]
		roles := userRoles.([]string)
		set := make(map[string]bool)
		for _, v := range roles {
			set[v] = true
		}
		_, ok := set["Admin"]
		if id != userId && !ok {
			ReturnUnauthorized(c)
		} else if id == userId && !ok {
			c.Next()
		} else if id != userId && ok {
			c.Next()
		}
	}
}

func Authorization(validRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(c.Keys) == 0 {
			ReturnUnauthorized(c)
		}
		rolesVal := c.Keys["roles"]
		if rolesVal == nil {
			ReturnUnauthorized(c)
		}
		roles := rolesVal.([]string)
		validation := make(map[string]bool)
		for _, val := range roles {
			validation[val] = true
		}
		for _, val := range validRoles {
			if _, ok := validation[val]; !ok {
				ReturnUnauthorized(c)
			}
		}
	}
}
