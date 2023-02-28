package routers

import (
	"api-restaurante/controllers"
	"api-restaurante/middlewares"

	"github.com/gin-gonic/gin"
)

type AuthRouteController struct {
	authController controllers.AuthController
}

func NewAuthRouteController(authController controllers.AuthController) AuthRouteController {
	return AuthRouteController{authController}
}

func (rc *AuthRouteController) AuthRoute(rg *gin.RouterGroup) {
	router := rg.Group("/auth")

	router.POST("/login", rc.authController.Login)
	router.GET("/logout", middlewares.ValidateToken(), rc.authController.LogoutUser)
}
