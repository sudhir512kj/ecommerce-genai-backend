package models

// Payment represents a payment for an order
type Payment struct {
	ID     int64         `json:"id"`
	UserID int           `json:"user_id"`
	Amount float64       `json:"amount"`
	Status PaymentStatus `json:"status"`
}

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "Pending"
	PaymentStatusSuccessful PaymentStatus = "Successful"
	PaymentStatusFailed     PaymentStatus = "Failed"
)

// CreatePaymentRequest represents the request to create a new payment
type CreatePaymentRequest struct {
	UserID int     `json:"user_id"`
	Amount float64 `json:"amount"`
}

// UpdatePaymentRequest represents the request to update a payment
type UpdatePaymentRequest struct {
	ID     int           `json:"id"`
	Status PaymentStatus `json:"status"`
}
