package routers

import (
	"api-restaurante/controllers"

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
	router.POST("/", rc.productController.CreateProduct)
	router.PATCH("/:id", rc.productController.UpdateProduct)
	router.GET("/", rc.productController.GetAllProducts)
	router.DELETE("/:id", rc.productController.DeleteProduct)
}
