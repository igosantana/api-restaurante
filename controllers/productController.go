package controllers

import (
	"api-restaurante/awsBucket"
	"api-restaurante/initializers"
	"api-restaurante/models"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateProduct(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	fmt.Println(form)
	file := form.File["file"][0]
	name := form.Value["name"][0]
	description := form.Value["description"][0]
	price := form.Value["price"][0]
	category := form.Value["category"][0]
	imageUrl, err := awsBucket.UploadImage(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	var product models.Product
	p, err := strconv.ParseFloat(price, 64)
	if err != nil {
		log.Printf("error: %v", err)
	}
	product.Name = name
	product.Description = description
	product.Price = p
	product.Category = category
	product.Image = imageUrl
	initializers.DB.Create(&product)
	c.JSON(http.StatusOK, gin.H{
		"message": "product created",
	})
}
