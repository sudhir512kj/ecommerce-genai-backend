package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sudhir512kj/ecommerce_backend/internal/models"
	"github.com/sudhir512kj/ecommerce_backend/internal/repository"
)

// CartHandler handles the cart-related operations
type CartHandler struct {
	cartRepo repository.CartRepository
}

func NewCartHandler(cartRepo repository.CartRepository) *CartHandler {
	return &CartHandler{
		cartRepo: cartRepo,
	}
}

// AddToCart adds a product to the user's cart
func (h *CartHandler) AddToCart(c *gin.Context) {
	var req models.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, _ := c.Get("user_id")
	cart, err := h.cartRepo.AddToCart(c.Request.Context(), userId.(int), req.ProductID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func (h *CartHandler) GetCartInfo(c *gin.Context) {
	// uid, _ := c.Get("user_id")
	userID, _ := c.Get("user_id")

	cartItems, err := h.cartRepo.GetCart(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cart_items": cartItems,
	})
}

// UpdateCart updates the items in the user's cart
func (h *CartHandler) UpdateCart(c *gin.Context) {
	var req models.UpdateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, _ := c.Get("user_id")
	cart, err := h.cartRepo.UpdateCart(c.Request.Context(), userId.(int), req.CartItems)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// DeleteFromCart removes an item from the user's cart
func (h *CartHandler) DeleteFromCart(c *gin.Context) {
	var req models.DeleteFromCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, _ := c.Get("user_id")
	cart, err := h.cartRepo.DeleteFromCart(c.Request.Context(), userId.(int), req.ProductID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// SaveForLater moves an item from the cart to the user's saved items
func (h *CartHandler) SaveForLater(c *gin.Context) {
	var req models.SaveForLaterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, _ := c.Get("user_id")
	cart, err := h.cartRepo.SaveForLater(c.Request.Context(), userId.(int), req.ProductID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}
