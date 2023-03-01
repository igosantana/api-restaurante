package routers

import (
	"api-restaurante/controllers"
	"api-restaurante/middlewares"

	"github.com/gin-gonic/gin"
)

type ProductRouterController struct {
	productController controllers.ProductController
}

func NewProductRouterController(productController controllers.ProductController) ProductRouterController {
	return ProductRouterController{productController}
}

func (rc *ProductRouterController) ProductRoute(gr *gin.RouterGroup) {
	router := gr.Group("/product")
	router.GET("", rc.productController.GetAllProducts)
	router.Use(middlewares.ValidateToken(), middlewares.Authorization([]string{controllers.Admin, controllers.Owner}))
	router.POST("/", rc.productController.CreateProduct)
	router.PATCH("/:id", rc.productController.UpdateProduct)
	router.DELETE("/:id", rc.productController.DeleteProduct)
}
