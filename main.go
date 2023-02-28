package main

import (
	"api-restaurante/controllers"
	"api-restaurante/initializers"
	"api-restaurante/routers"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	server              *gin.Engine
	AuthController      controllers.AuthController
	AuthRouteController routers.AuthRouteController

	UserController      controllers.UserController
	UserRouteController routers.UserRouterController

	ProductController      controllers.ProductController
	ProductRouteController routers.ProductRouterController
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()

	AuthController = controllers.NewAuthController(initializers.DB)
	AuthRouteController = routers.NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(initializers.DB)
	UserRouteController = routers.NewUserRouteController(UserController)

	ProductController = controllers.NewProductController(initializers.DB)
	ProductRouteController = routers.NewProductRouterController(ProductController)

	server = gin.Default()
}

func main() {
	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Welcome to restaurant api"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})
	AuthRouteController.AuthRoute(router)
	UserRouteController.UserRoute(router)
	ProductRouteController.ProductRoute(router)
	server.Run()
}
