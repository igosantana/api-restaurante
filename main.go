package main

import (
	"api-restaurante/controllers"
	_ "api-restaurante/docs"
	"api-restaurante/initializers"
	"api-restaurante/routers"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

// @title Restaurant service API
// @version 1.0
// @description A Restaurant service API in Go using Gin framework
// @host 	localhost:3000
// @BasePath /api
func main() {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowCredentials = true
	server.Use(cors.New(corsConfig))
	router := server.Group("/api")
	// add swagger
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Welcome to restaurant api"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})
	AuthRouteController.AuthRoute(router)
	UserRouteController.UserRoute(router)
	ProductRouteController.ProductRoute(router)
	server.Run()
}
