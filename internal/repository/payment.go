package repository

import (
	"context"
	"database/sql"

	"github.com/sudhir512kj/ecommerce_backend/internal/models"
)

// PaymentRepository is an interface that defines the methods for interacting with payments
type PaymentRepository interface {
	ProcessPayment(ctx context.Context, userID int, amount float64) (*models.Payment, error)
}

type paymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) PaymentRepository {
	return &paymentRepository{
		db: db,
	}
}

func (r *paymentRepository) ProcessPayment(ctx context.Context, userID int, amount float64) (*models.Payment, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var payment models.Payment
	result, err := tx.ExecContext(ctx, "INSERT INTO payments (user_id, amount) VALUES (?, ?)", userID, amount)
	if err != nil {
		return nil, err
	}
	payment.ID, _ = result.LastInsertId()
	payment.Amount = amount

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &payment, nil
}
