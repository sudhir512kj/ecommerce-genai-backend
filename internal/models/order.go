package models

// Order represents a user's order
type Order struct {
	ID        int64       `json:"id"`
	UserID    int         `json:"user_id"`
	PaymentID int         `json:"payment_id"`
	Status    OrderStatus `json:"status"`
	Items     []OrderItem `json:"items"`
}

// OrderItem represents an item in the user's order
type OrderItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPaymentComplete OrderStatus = "Payment Complete"
	OrderStatusPending         OrderStatus = "Pending"
	OrderStatusPlanning        OrderStatus = "Planning"
	OrderStatusShipping        OrderStatus = "Shipping"
	OrderStatusComplete        OrderStatus = "Complete"
	OrderStatusUnfulfillable   OrderStatus = "Unfulfillable"
	OrderStatusCancelled       OrderStatus = "Cancelled"
)
