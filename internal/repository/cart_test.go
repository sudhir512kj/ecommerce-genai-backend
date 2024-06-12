package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/sudhir512kj/ecommerce_backend/internal/models"
)

func TestCartRepository_CreateCart(t *testing.T) {
	// Set up the test environment
	db := dbTest.GetDb()

	cartRepo := NewCartRepository(db)
	productRepo := NewProductRepository(db)
	userRepo := NewUserRepository(db)

	t.Run("should create a new cart and add product to it", func(t *testing.T) {
		newEmail := generateRandomEmail()

		// Arrange
		user := &models.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     newEmail,
			Password:  "password123",
		}
		err := userRepo.CreateUser(context.Background(), user)
		require.NoError(t, err)

		EUser, err := userRepo.GetUserByEmail(context.Background(), user.Email)
		require.NoError(t, err)

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

		// Act
		cart, err := cartRepo.AddToCart(context.Background(), EUser.ID, product.ID, 2)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, cart.ID)
		assert.Len(t, cart.Items, 1)
		assert.Equal(t, product.ID, cart.Items[0].ProductID)
		assert.Equal(t, int64(2), cart.Items[0].Quantity)
	})

	t.Run("should return error when user does not exist", func(t *testing.T) {

		// Arrange
		nonExistentUserID := int(999)
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

		// Act
		_, err = cartRepo.AddToCart(context.Background(), nonExistentUserID, product.ID, 2)

		// Assert
		assert.Error(t, err)
		// assert.Equal(t, ErrUserNotFound, err)
	})

	t.Run("should return error when product does not exist", func(t *testing.T) {
		newEmail := generateRandomEmail()

		// Arrange
		user := &models.User{
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     newEmail,
			Password:  "password123",
		}
		err := userRepo.CreateUser(context.Background(), user)
		require.NoError(t, err)

		EUser, err := userRepo.GetUserByEmail(context.Background(), user.Email)
		require.NoError(t, err)

		nonExistentProductID := int(999)

		// Act
		_, err = cartRepo.AddToCart(context.Background(), EUser.ID, nonExistentProductID, 2)

		// Assert
		assert.Error(t, err)
		// assert.Equal(t, ErrProductNotFound, err)
	})
}

func TestCartRepository_GetCartByUserID(t *testing.T) {
	// Set up the test environment
	db := dbTest.GetDb()

	cartRepo := NewCartRepository(db)
	userRepo := NewUserRepository(db)

	t.Run("should get cart by user ID", func(t *testing.T) {
		newEmail := generateRandomEmail()

		// Arrange
		user := &models.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     newEmail,
			Password:  "password123",
		}
		err := userRepo.CreateUser(context.Background(), user)
		require.NoError(t, err)

		EUser, err := userRepo.GetUserByEmail(context.Background(), user.Email)
		require.NoError(t, err)

		cart, err := cartRepo.AddToCart(context.Background(), EUser.ID, 1, 2)
		require.NoError(t, err)

		// Act
		foundCart, err := cartRepo.GetCart(context.Background(), EUser.ID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, cart.ID, foundCart.ID)
		assert.Equal(t, cart.UserID, foundCart.UserID)
		assert.Equal(t, cart.Items, foundCart.Items)
	})

	t.Run("should return error when user does not have a cart", func(t *testing.T) {
		newEmail := generateRandomEmail()

		// Arrange
		user := &models.User{
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     newEmail,
			Password:  "password123",
		}
		err := userRepo.CreateUser(context.Background(), user)
		require.NoError(t, err)

		Euser, err := userRepo.GetUserByEmail(context.Background(), user.Email)
		require.NoError(t, err)

		// Act
		_, err = cartRepo.GetCart(context.Background(), Euser.ID)

		// Assert
		assert.Error(t, err)
		// assert.Equal(t, ErrCartNotFound, err)
	})
}

func TestCartRepository_AddItemToCart(t *testing.T) {
	// Set up the test environment
	db := dbTest.GetDb()

	cartRepo := NewCartRepository(db)
	productRepo := NewProductRepository(db)
	userRepo := NewUserRepository(db)

	t.Run("should add item to cart", func(t *testing.T) {
		// Arrange
		user := &models.User{
			FirstName: "Bob",
			LastName:  "Smith",
			Email:     "bob@example.com",
			Password:  "password123",
		}
		err := userRepo.CreateUser(context.Background(), user)
		require.NoError(t, err)

		EUser, err := userRepo.GetUserByEmail(context.Background(), user.Email)
		require.NoError(t, err)

		_, err = cartRepo.AddToCart(context.Background(), EUser.ID, 1, 2)
		require.NoError(t, err)

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

		// Act
		_, err = cartRepo.AddToCart(context.Background(), EUser.ID, addedProduct.ID, 3)

		// Assert
		assert.NoError(t, err)

		updatedCart, err := cartRepo.GetCart(context.Background(), EUser.ID)
		require.NoError(t, err)
		assert.Len(t, updatedCart.Items, 2)
		assert.Equal(t, product.ID, updatedCart.Items[1].ProductID)
		assert.Equal(t, int64(3), updatedCart.Items[1].Quantity)
	})

	t.Run("should return error when cart does not exist", func(t *testing.T) {
		// Arrange
		nonExistentUserID := int(999)
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

		// Act
		_, err = cartRepo.AddToCart(context.Background(), nonExistentUserID, addedProduct.ID, 2)

		// Assert
		assert.Error(t, err)
		// assert.Equal(t, ErrCartNotFound, err)
	})

	t.Run("should return error when product does not exist", func(t *testing.T) {
		// Arrange
		user := &models.User{
			FirstName: "Alice",
			LastName:  "Johnson",
			Email:     "alice@example.com",
			Password:  "password123",
		}
		err := userRepo.CreateUser(context.Background(), user)
		require.NoError(t, err)

		EUser, err := userRepo.GetUserByEmail(context.Background(), user.Email)
		require.NoError(t, err)

		_, err = cartRepo.AddToCart(context.Background(), user.ID, 1, 2)
		require.NoError(t, err)

		nonExistentProductID := int(999)

		// Act
		_, err = cartRepo.AddToCart(context.Background(), EUser.ID, nonExistentProductID, 2)

		// Assert
		assert.Error(t, err)
		// assert.Equal(t, ErrProductNotFound, err)
	})
}

func TestCartRepository_RemoveItemFromCart(t *testing.T) {
	db := dbTest.GetDb()

	cartRepo := NewCartRepository(db)
	// productRepo := NewProductRepository(db)
	userRepo := NewUserRepository(db)

	// t.Run("should remove item from cart", func(t *testing.T) {
	// 	// Arrange
	// 	user := &models.User{
	// 		FirstName: "Charlie",
	// 		LastName:  "Brown",
	// 		Email:     "charlie@example.com",
	// 		Password:  "password123",
	// 	}
	// 	err := userRepo.CreateUser(context.Background(), user)
	// 	require.NoError(t, err)

	// 	cart, err := cartRepo.AddToCart(context.Background(), user.ID, 1, 2)
	// 	require.NoError(t, err)

	// 	product := &models.Product{
	// 		Name:          "Test Product",
	// 		Description:   "This is a test product",
	// 		Price:         19.99,
	// 		CategoryID:    1,
	// 		SubCategoryID: 1,
	// 	}
	// 	addedProduct, err := productRepo.CreateProduct(context.Background(), product)

	// 	// Assert
	// 	assert.NoError(t, err)
	// 	assert.NotEmpty(t, addedProduct.ID)
	// 	assert.NotEmpty(t, addedProduct.CreatedAt)
	// 	assert.NotEmpty(t, addedProduct.UpdatedAt)

	// 	err = cartRepo.AddToCart(context.Background(), cart.ID, product.ID, 3)
	// 	require.NoError(t, err)

	// 	// Act
	// 	err = cartRepo.RemoveItemFromCart(context.Background(), cart.ID, product.ID)

	// 	// Assert
	// 	assert.NoError(t, err)

	// 	updatedCart, err := cartRepo.GetCartByUserID(context.Background(), user.ID)
	// 	require.NoError(t, err)
	// 	assert.Len(t, updatedCart.Items, 1)
	// })

	// t.Run("should return error when cart does not exist", func(t *testing.T) {
	// 	// Arrange
	// 	nonExistentCartID := int(999)
	// 	product := &models.Product{
	// 		Name:          "Test Product",
	// 		Description:   "This is a test product",
	// 		Price:         19.99,
	// 		CategoryID:    1,
	// 		SubCategoryID: 1,
	// 	}
	// 	err := productRepo.CreateProduct(context.Background(), product)
	// 	require.NoError(t, err)

	// 	// Act
	// 	err = cartRepo.DeleteFromCart(context.Background(), nonExistentCartID, product.ID)

	// 	// Assert
	// 	assert.Error(t, err)
	// 	// assert.Equal(t, ErrCartNotFound, err)
	// })

	t.Run("should return error when product does not exist in cart", func(t *testing.T) {
		// Arrange
		user := &models.User{
			FirstName: "David",
			LastName:  "Lee",
			Email:     "david@example.com",
			Password:  "password123",
		}
		err := userRepo.CreateUser(context.Background(), user)
		require.NoError(t, err)

		EUser, err := userRepo.GetUserByEmail(context.Background(), user.Email)
		require.NoError(t, err)

		_, err = cartRepo.AddToCart(context.Background(), EUser.ID, 1, 2)
		require.NoError(t, err)

		nonExistentProductID := int(999)

		// Act
		_, err = cartRepo.DeleteFromCart(context.Background(), EUser.ID, nonExistentProductID)

		// Assert
		assert.Error(t, err)
		// assert.Equal(t, ErrCartItemNotFound, err)
	})
}
