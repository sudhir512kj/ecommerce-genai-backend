package models

// Cart represents a user's shopping cart
type Cart struct {
	ID          int64      `json:"id"`
	UserID      int        `json:"user_id"`
	TotalAmount float64    `json:"total_amount"`
	Items       []CartItem `json:"items"`
}

// CartItem represents an item in the user's shopping cart
type CartItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

// AddToCartRequest represents the request to add a product to the cart
type AddToCartRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

// UpdateCartRequest represents the request to update the cart
type UpdateCartRequest struct {
	CartItems []CartItem `json:"cart_items"`
}

// DeleteFromCartRequest represents the request to delete a product from the cart
type DeleteFromCartRequest struct {
	ProductID int `json:"product_id"`
}

// SaveForLaterRequest represents the request to save a product for later
type SaveForLaterRequest struct {
	ProductID int `json:"product_id"`
}
