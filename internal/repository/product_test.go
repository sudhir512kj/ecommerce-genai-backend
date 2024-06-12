package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sudhir512kj/ecommerce_backend/internal/models"
)

func TestProductRepository(t *testing.T) {

	productRepo := NewProductRepository(dbTest.GetDb())

	t.Run("CreateProduct", func(t *testing.T) {
		// Arrange
		product := &models.Product{
			Name:          "Test Product",
			Description:   "This is a test product",
			Price:         19.99,
			CategoryID:    1,
			SubCategoryID: 1,
		}

		// Act
		addedProduct, err := productRepo.CreateProduct(context.Background(), product)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, addedProduct.ID)
		assert.NotEmpty(t, addedProduct.CreatedAt)
		assert.NotEmpty(t, addedProduct.UpdatedAt)
	})

	t.Run("GetProductByID", func(t *testing.T) {
		// Arrange
		product := &models.Product{
			Name:          "Another Test Product",
			Description:   "This is another test product",
			Price:         24.99,
			CategoryID:    1,
			SubCategoryID: 1,
		}
		addedProduct, err := productRepo.CreateProduct(context.Background(), product)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, addedProduct.ID)
		assert.NotEmpty(t, addedProduct.CreatedAt)
		assert.NotEmpty(t, addedProduct.UpdatedAt)

		// Act
		foundProduct, err := productRepo.GetProductByID(context.Background(), addedProduct.ID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, addedProduct.ID, foundProduct.ID)
		assert.Equal(t, addedProduct.Name, foundProduct.Name)
		assert.Equal(t, addedProduct.Description, foundProduct.Description)
		assert.Equal(t, addedProduct.Price, foundProduct.Price)
	})

	t.Run("UpdateProduct", func(t *testing.T) {
		// Arrange
		product := &models.Product{
			Name:          "Test Product",
			Description:   "This is a test product",
			Price:         19.99,
			CategoryID:    1,
			SubCategoryID: 1,
		}
		addedProduct, err := productRepo.CreateProduct(context.Background(), product)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, addedProduct.ID)
		assert.NotEmpty(t, addedProduct.CreatedAt)
		assert.NotEmpty(t, addedProduct.UpdatedAt)

		updatedProduct := &models.Product{
			ID:          product.ID,
			Name:        "Updated Test Product",
			Description: "This is an updated test product",
			Price:       24.99,
		}

		// Act
		err = productRepo.UpdateProduct(context.Background(), updatedProduct)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "Updated Test Product", updatedProduct.Name)
		assert.Equal(t, "This is an updated test product", updatedProduct.Description)
		assert.Equal(t, float64(24.99), updatedProduct.Price)
	})

	t.Run("DeleteProduct", func(t *testing.T) {
		// Arrange
		product := &models.Product{
			Name:          "Test Product to Delete",
			Description:   "This is a test product to delete",
			Price:         19.99,
			CategoryID:    1,
			SubCategoryID: 1,
		}
		addedProduct, err := productRepo.CreateProduct(context.Background(), product)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, addedProduct.ID)
		assert.NotEmpty(t, addedProduct.CreatedAt)
		assert.NotEmpty(t, addedProduct.UpdatedAt)

		// Act
		err = productRepo.DeleteProduct(context.Background(), addedProduct.ID)

		// Assert
		assert.NoError(t, err)

		// Verify the product is deleted
		_, err = productRepo.GetProductByID(context.Background(), addedProduct.ID)
		assert.NoError(t, err)
	})

	t.Run("ListProducts", func(t *testing.T) {
		// Arrange
		product1 := &models.Product{
			Name:          "Test Product 1",
			Description:   "This is test product 1",
			Price:         9.99,
			CategoryID:    1,
			SubCategoryID: 1,
		}
		addedProduct1, err := productRepo.CreateProduct(context.Background(), product1)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, addedProduct1.ID)
		assert.NotEmpty(t, addedProduct1.CreatedAt)
		assert.NotEmpty(t, addedProduct1.UpdatedAt)

		product2 := &models.Product{
			Name:          "Test Product 2",
			Description:   "This is test product 2",
			Price:         14.99,
			CategoryID:    1,
			SubCategoryID: 1,
		}
		addedProduct2, err := productRepo.CreateProduct(context.Background(), product2)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, addedProduct2.ID)
		assert.NotEmpty(t, addedProduct2.CreatedAt)
		assert.NotEmpty(t, addedProduct2.UpdatedAt)

		// Act
		products, err := productRepo.GetAllProducts(context.Background(), 2, 2)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, products, 2)
	})
}
