package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sudhir512kj/ecommerce_backend/internal/repository"
)

// OrderHandler handles the order-related operations
type OrderHandler struct {
	orderRepo   repository.OrderRepository
	cartRepo    repository.CartRepository
	paymentRepo repository.PaymentRepository
}

func NewOrderHandler(orderRepo repository.OrderRepository, cartRepo repository.CartRepository, paymentRepo repository.PaymentRepository) *OrderHandler {
	return &OrderHandler{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		paymentRepo: paymentRepo,
	}
}

// Checkout creates a new order and processes the payment
func (h *OrderHandler) Checkout(c *gin.Context) {
	userId, _ := c.Get("user_id")
	cart, err := h.cartRepo.GetCart(c.Request.Context(), userId.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.paymentRepo.ProcessPayment(c.Request.Context(), userId.(int), cart.TotalAmount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderRepo.CreateOrder(c.Request.Context(), userId.(int), cart.Items, (int)(payment.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetOrders retrieves all orders associated with the user
func (h *OrderHandler) GetOrders(c *gin.Context) {
	userId, _ := c.Get("user_id")
	orders, err := h.orderRepo.GetOrdersByUserID(c.Request.Context(), userId.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetOrder retrieves a specific order associated with the user
func (h *OrderHandler) GetOrder(c *gin.Context) {
	userId, _ := c.Get("user_id")
	orderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := h.orderRepo.GetOrderByID(c.Request.Context(), userId.(int), orderId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// CancelOrder cancels a specific order associated with the user
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	userId, _ := c.Get("user_id")
	orderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	err = h.orderRepo.CancelOrder(c.Request.Context(), userId.(int), orderId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})
}
