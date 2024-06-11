package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sudhir512kj/ecommerce_backend/internal/models"
	"github.com/sudhir512kj/ecommerce_backend/internal/repository"
)

type ProductHandler struct {
	productRepo      repository.ProductRepository
	categoryRepo     repository.CategoryRepository
	productImageRepo repository.ProductImageRepository
}

func NewProductHandler(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository, productImageRepo repository.ProductImageRepository) *ProductHandler {
	return &ProductHandler{
		productRepo:      productRepo,
		categoryRepo:     categoryRepo,
		productImageRepo: productImageRepo,
	}
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	// Implement the logic to fetch all products with pagination
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}

	products, err := h.productRepo.GetAllProducts(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	// Implement the logic to fetch a product by ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := h.productRepo.GetProductByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	// Implement the logic to create a new product
	var product *models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Check if the category ID exists
	category, err := h.categoryRepo.GetCategoryByID(c.Request.Context(), product.CategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if category == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Create the product
	newProduct, err := h.productRepo.CreateProduct(c.Request.Context(), product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, newProduct)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	// Implement the logic to update an existing product
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.productRepo.GetProductByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Check if the category and subcategory exist
	category, err := h.categoryRepo.GetCategoryByID(c.Request.Context(), req.CategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if category == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category not found"})
		return
	}

	subcategory, err := h.categoryRepo.GetSubcategoryByID(c.Request.Context(), req.SubCategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if subcategory == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Subcategory not found"})
		return
	}

	// Update the product
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.CategoryID = req.CategoryID
	product.SubCategoryID = req.SubCategoryID
	product.Images = req.Images
	product.UpdatedAt = time.Now()

	if err := h.productRepo.UpdateProduct(c.Request.Context(), product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update product images
	if err := h.productImageRepo.DeleteProductImagesByProductID(c.Request.Context(), product.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, imageURL := range req.Images {
		productImage := &models.ProductImage{
			ProductID: product.ID,
			ImageURL:  imageURL,
		}
		_, err := h.productImageRepo.CreateProductImage(c.Request.Context(), productImage)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, product)

}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	// Implement the logic to delete a product
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := h.productRepo.GetProductByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if err := h.productImageRepo.DeleteProductImagesByProductID(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.productRepo.DeleteProduct(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
