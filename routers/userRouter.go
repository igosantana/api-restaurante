package routers

import (
	"api-restaurante/controllers"

	"github.com/gin-gonic/gin"
)

func HandleRouters() {
	r := gin.Default()
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.Run()
}
