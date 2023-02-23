package routers

import (
	"api-restaurante/controllers"
	"api-restaurante/middlewares"

	"github.com/gin-gonic/gin"
)

func HandleRouters() {
	r := gin.Default()
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.PATCH("/user/:id", middlewares.AuthIsUser, controllers.UpdateUser)
	r.DELETE("/user/:id", middlewares.AuthIsUser, controllers.DeleteUser)
	r.GET("/user/:id", middlewares.AuthIsUser, controllers.GetUser)
	r.Run()
}
