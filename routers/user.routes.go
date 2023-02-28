package routers

import (
	"api-restaurante/controllers"
	"api-restaurante/middlewares"

	"github.com/gin-gonic/gin"
)

type UserRouterController struct {
	userController controllers.UserController
}

func NewUserRouteController(userController controllers.UserController) UserRouterController {
	return UserRouterController{userController}
}

func (rc *UserRouterController) UserRoute(rg *gin.RouterGroup) {
	router := rg.Group("/user")

	router.POST("/", rc.userController.CreateUser)
	router.Use(middlewares.ValidateToken(), middlewares.IsOwnerOrAdmin())
	router.PATCH("/:id", rc.userController.UpdateUser)
	router.DELETE("/:id", rc.userController.DeleteUser)
	router.GET("/:id", rc.userController.GetUser)
	router.GET("/", rc.userController.GetAllUsers)
}
