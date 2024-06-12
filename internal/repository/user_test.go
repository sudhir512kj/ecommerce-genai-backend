package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/sudhir512kj/ecommerce_backend/internal/models"
)

func TestUserRepository(t *testing.T) {
	// Set up the test environment

	userRepo := NewUserRepository(dbTest.GetDb())

	t.Run("CreateUser", func(t *testing.T) {
		newEmail := generateRandomEmail()

		// Arrange
		user := &models.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     newEmail,
			Password:  "password123",
		}

		// Act
		err := userRepo.CreateUser(context.Background(), user)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, user.ID)
	})

	t.Run("GetUserByEmail", func(t *testing.T) {
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

		// Act
		foundUser, err := userRepo.GetUserByEmail(context.Background(), newEmail)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, user.ID, foundUser.ID)
		assert.Equal(t, user.FirstName, foundUser.FirstName)
		assert.Equal(t, user.LastName, foundUser.LastName)
		assert.Equal(t, user.Email, foundUser.Email)
	})

	t.Run("GetUserByID", func(t *testing.T) {
		newEmail := generateRandomEmail()

		// Arrange
		user := &models.User{
			FirstName: "Bob",
			LastName:  "Smith",
			Email:     newEmail,
			Password:  "password123",
		}
		err := userRepo.CreateUser(context.Background(), user)
		require.NoError(t, err)

		// Act
		foundUser, err := userRepo.GetUserByID(context.Background(), user.ID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, user.ID, foundUser.ID)
		assert.Equal(t, user.FirstName, foundUser.FirstName)
		assert.Equal(t, user.LastName, foundUser.LastName)
		assert.Equal(t, user.Email, foundUser.Email)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		newEmail := generateRandomEmail()

		// Arrange
		user := &models.User{
			FirstName: "Alice",
			LastName:  "Johnson",
			Email:     newEmail,
			Password:  "password123",
		}
		err := userRepo.CreateUser(context.Background(), user)
		require.NoError(t, err)

		updatedUser := &models.User{
			ID:        user.ID,
			FirstName: "Alice",
			LastName:  "Williams",
			Email:     newEmail,
		}

		// Act
		err = userRepo.UpdateUser(context.Background(), updatedUser)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "Alice", updatedUser.FirstName)
		assert.Equal(t, "Williams", updatedUser.LastName)
		assert.Equal(t, newEmail, updatedUser.Email)
	})
}
