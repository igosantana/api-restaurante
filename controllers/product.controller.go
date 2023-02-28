package controllers

import (
	"api-restaurante/awsBucket"
	"api-restaurante/initializers"
	"api-restaurante/models"
	"api-restaurante/utils"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type ProductController struct {
	DB *gorm.DB
}

func NewProductController(DB *gorm.DB) ProductController {
	return ProductController{DB}
}

func (pc *ProductController) CreateProduct(c *gin.Context) {
	var form models.CreateProductForm
	err := c.ShouldBind(&form)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]utils.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = utils.ErrorMsg{Field: fe.Field(), Message: utils.GetErrorMsg(fe)}
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
		}
		return
	}
	imageUrl, error := awsBucket.UploadImage(form.File)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": error,
		})
		return
	}
	p, err := strconv.ParseFloat(form.Price, 64)
	if err != nil {
		log.Printf("error: %v", err)
	}
	q, err := strconv.Atoi(form.Quantity)
	if err != nil {
		log.Printf("error: %v", err)
	}
	var product models.Product
	product.Name = form.Name
	product.Description = form.Description
	product.Price = p
	product.Category = form.Category
	product.Image = imageUrl
	product.Quantity = q
	initializers.DB.Create(&product)
	c.JSON(http.StatusOK, gin.H{
		"message": "product created",
	})
}

func (pc *ProductController) UpdateProduct(c *gin.Context) {
	var updateProduct models.UpdateProduct
	err := c.ShouldBind(&updateProduct)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]utils.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = utils.ErrorMsg{Field: fe.Field(), Message: utils.GetErrorMsg(fe)}
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
		}
		return
	}
	var product models.Product
	id := c.Param("id")
	result := initializers.DB.First(&product, "id = ?", id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Product not found",
		})
		return
	}
	initializers.DB.Model(&product).Updates(updateProduct.ToUpdateProductModel())
	c.JSON(http.StatusOK, gin.H{
		"message": "product updated",
	})
}

func (pc *ProductController) GetAllProducts(c *gin.Context) {
	var products []models.Product

	result := initializers.DB.Find(&products)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Product not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": &products,
	})
}

func (pc *ProductController) DeleteProduct(c *gin.Context) {
	var product models.Product
	id := c.Param("id")
	result := initializers.DB.First(&product, "id = ?", id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Product not found",
		})
		return
	}
	initializers.DB.Delete(&product)
}
