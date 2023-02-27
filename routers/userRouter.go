package routers

import (
	"api-restaurante/controllers"
	"api-restaurante/middlewares"

	"github.com/gin-gonic/gin"
)

func HandleRouters() {
	r := gin.Default()
	r.POST("/login", controllers.Login)
	// User routers
	user := r.Group("/user")
	user.POST("/", controllers.CreateUser)
	user.Use(middlewares.ValidateToken(), middlewares.IsOwnerOrAdmin())
	user.PATCH("/:id", controllers.UpdateUser)
	user.DELETE("/:id", controllers.DeleteUser)
	user.GET("/:id", controllers.GetUser)
	// Product routers
	product := r.Group("/product")
	product.POST("/", controllers.CreateProduct)
	product.PATCH("/:id", controllers.UpdateProduct)
	product.GET("/", controllers.GetAllProducts)
	product.DELETE("/:id", controllers.DeleteProduct)
	r.Run()
}
