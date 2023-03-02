package controllers

import (
	"api-restaurante/data/response"
	"api-restaurante/helper"
	"api-restaurante/initializers"
	"api-restaurante/models"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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

// CreateProduct  godoc
// @Summary Create products
// @Description Save products data in DB. Only Admin or Owner.
// @accept mpfd
// @Param Product body models.CreateProductForm true "A JSON form object containing the products requirements"
// @Tags Product
// @Security ApiKeyAuth
// @Produce application/json
// @Success 200 {object} response.Response{}
// @Router /product [post]
func (pc *ProductController) CreateProduct(c *gin.Context) {
	var form models.CreateProductForm
	err := c.Bind(&form)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]response.ErrorDetail, len(ve))
			for i, fe := range ve {
				out[i] = response.ErrorDetail{Field: fe.Field(), ErrorMessage: response.GetErrorMsg(fe)}
			}
			response.BadRequest(c, http.StatusBadRequest, "validation error", out)
		}
		return
	}
	imageUrl, err := helper.ImageUploadHelper(form.File)
	if err != nil {
		var errors []response.ErrorDetail = make([]response.ErrorDetail, 0, 1)
		errors = append(errors, response.ErrorDetail{
			ErrorType:    response.ErrorTypeError,
			ErrorMessage: fmt.Sprintf("%v", err),
		})
		response.BadRequest(c, http.StatusBadRequest, "failed to upload image", errors)
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
	create := initializers.DB.Create(&product)
	if create.Error != nil {
		var errors []response.ErrorDetail = make([]response.ErrorDetail, 0, 1)
		errors = append(errors, response.ErrorDetail{
			ErrorType:    response.ErrorTypeError,
			ErrorMessage: fmt.Sprintf("%v", create.Error),
		})
		response.BadRequest(c, http.StatusBadRequest, "failed to create product", errors)
		return
	}
	response.Ok(c, http.StatusOK, "success", nil)
}

// UpdateProduct  godoc
// @Summary Update products
// @Description Update and Save products data in DB. Only Admin or Owner.
// @Param Product body models.UpdateProduct true "A JSON form object containing the products requirements"
// @Param id path string true "Update product by id."
// @Tags Product
// @Security ApiKeyAuth
// @Produce application/json
// @Success 200 {object} response.Response{}
// @Router /product/{id} [patch]
func (pc *ProductController) UpdateProduct(c *gin.Context) {
	var updateProduct models.UpdateProduct
	err := c.BindJSON(&updateProduct)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]response.ErrorDetail, len(ve))
			for i, fe := range ve {
				out[i] = response.ErrorDetail{Field: fe.Field(), ErrorMessage: response.GetErrorMsg(fe)}
			}
			response.BadRequest(c, http.StatusBadRequest, "validation error", out)
		}
		return
	}
	var product models.Product
	id := c.Param("id")
	result := initializers.DB.First(&product, "id = ?", id)
	if result.RowsAffected == 0 {
		var errors []response.ErrorDetail = make([]response.ErrorDetail, 0, 1)
		errors = append(errors, response.ErrorDetail{
			ErrorType:    response.ErrorTypeError,
			ErrorMessage: fmt.Sprintf("product not found"),
		})
		response.BadRequest(c, http.StatusNotFound, "product not found", errors)
		return
	}
	update := initializers.DB.Model(&product).Updates(updateProduct.ToUpdateProductModel())
	if update.Error != nil {
		var errors []response.ErrorDetail = make([]response.ErrorDetail, 0, 1)
		errors = append(errors, response.ErrorDetail{
			ErrorType:    response.ErrorTypeError,
			ErrorMessage: fmt.Sprintf("%v", update.Error),
		})
		response.BadRequest(c, http.StatusBadRequest, "failed to update product", errors)
		return
	}
	response.Ok(c, http.StatusOK, "success", nil)
}

// GetAllProducts  godoc
// @Summary Get all products
// @Description Get all products data in DB.
// @Param name query string false "Get products by name"
// @Param category query string false "Get products by category"
// @Tags Product
// @Produce application/json
// @Success 200 {object} response.Response{} "A JSON with products"
// @Router /product [get]
func (pc *ProductController) GetAllProducts(c *gin.Context) {
	var productsDB []models.Product
	name := c.Query("name")
	category := c.Query("category")
	result := initializers.DB.Find(&productsDB)
	if result.RowsAffected == 0 {
		var errors []response.ErrorDetail = make([]response.ErrorDetail, 0, 1)
		errors = append(errors, response.ErrorDetail{
			ErrorType:    response.ErrorTypeError,
			ErrorMessage: fmt.Sprintf("products not found"),
		})
		response.BadRequest(c, http.StatusNotFound, "products not found", errors)
		return
	}
	if name != "" {
		result := initializers.DB.Find(&productsDB, "name LIKE ?", fmt.Sprintf("%%%s%%", strings.ToLower(name)))
		if result.RowsAffected == 0 {
			var errors []response.ErrorDetail = make([]response.ErrorDetail, 0, 1)
			errors = append(errors, response.ErrorDetail{
				ErrorType:    response.ErrorTypeError,
				ErrorMessage: fmt.Sprintf("products not found"),
			})
			response.BadRequest(c, http.StatusNotFound, "products not found", errors)
			return
		}
	}
	if category != "" {
		result := initializers.DB.Find(&productsDB, "category = ?", strings.ToLower(category))
		if result.RowsAffected == 0 {
			var errors []response.ErrorDetail = make([]response.ErrorDetail, 0, 1)
			errors = append(errors, response.ErrorDetail{
				ErrorType:    response.ErrorTypeError,
				ErrorMessage: fmt.Sprintf("products not found"),
			})
			response.BadRequest(c, http.StatusNotFound, "products not found", errors)
			return
		}
	}
	var allProducts []models.ToGetAllProducts
	for _, product := range productsDB {
		allProducts = append(allProducts, product.ToGetAll())
	}
	response.Ok(c, http.StatusOK, "success", &allProducts)
}

// DeleteProduct  godoc
// @Summary Delete product
// @Description Remove products data in DB. Only Admin or Owner.
// @Param id path string true "Remove product by id."
// @Tags Product
// @Security ApiKeyAuth
// @Produce application/json
// @Success 200 {object} response.Response{}
// @Router /product/{id} [delete]
func (pc *ProductController) DeleteProduct(c *gin.Context) {
	var product models.Product
	id := c.Param("id")
	result := initializers.DB.First(&product, "id = ?", id)
	if result.RowsAffected == 0 {
		var errors []response.ErrorDetail = make([]response.ErrorDetail, 0, 1)
		errors = append(errors, response.ErrorDetail{
			ErrorType:    response.ErrorTypeError,
			ErrorMessage: fmt.Sprintf("product not found"),
		})
		response.BadRequest(c, http.StatusNotFound, "product not found", errors)
		return
	}
	delete := initializers.DB.Delete(&product)
	if delete.Error != nil {
		var errors []response.ErrorDetail = make([]response.ErrorDetail, 0, 1)
		errors = append(errors, response.ErrorDetail{
			ErrorType:    response.ErrorTypeError,
			ErrorMessage: fmt.Sprintf("%v", delete.Error),
		})
		response.BadRequest(c, http.StatusBadRequest, "failed to delete product", errors)
		return
	}
	response.Ok(c, http.StatusOK, "success", nil)
}
