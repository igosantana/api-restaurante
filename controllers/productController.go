package controllers

import (
	"api-restaurante/awsBucket"
	"api-restaurante/initializers"
	"api-restaurante/models"
	"api-restaurante/utils"
	"errors"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type form struct {
	Name        string                `form:"name" binding:"required,max=110"`
	Description string                `form:"description" binding:"required,max=255"`
	Price       string                `form:"price" binding:"required"`
	Category    string                `form:"category" binding:"required,max=110"`
	Quantity    string                `form:"quantity" binding:"required"`
	File        *multipart.FileHeader `form:"file" binding:"required"`
}

func CreateProduct(c *gin.Context) {
	var form form
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

type updateProduct struct {
	Name        string  `json:"name,omitempty" binding:"max=110"`
	Description string  `json:"description,omitempty" binding:"max=250"`
	Price       float64 `json:"price,omitempty"`
	Category    string  `json:"category,omitempty" binding:"max=110"`
	Quantity    int     `json:"quantity,omitempty"`
}

func (p *updateProduct) toModel() *models.Product {
	return &models.Product{
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Category:    p.Category,
		Quantity:    p.Quantity,
	}
}

func UpdateProduct(c *gin.Context) {
	var updateProduct updateProduct
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
	initializers.DB.Model(&product).Updates(updateProduct.toModel())
	c.JSON(http.StatusOK, gin.H{
		"message": "product updated",
	})
}

func GetAllProducts(c *gin.Context) {
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

func DeleteProduct(c *gin.Context) {
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
