package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sudhir512kj/ecommerce_backend/internal/models"
)

// OrderRepository is an interface that defines the methods for interacting with orders
type OrderRepository interface {
	CreateOrder(ctx context.Context, userID int, cartItems []models.CartItem, paymentID int) (*models.Order, error)
	GetOrdersByUserID(ctx context.Context, userID int) ([]*models.Order, error)
	GetOrderByID(ctx context.Context, userID, orderID int) (*models.Order, error)
	CancelOrder(ctx context.Context, userID, orderID int) error
	getOrderItems(ctx context.Context, orderID int) ([]models.OrderItem, error)
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) CreateOrder(ctx context.Context, userID int, cartItems []models.CartItem, paymentID int) (*models.Order, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var order models.Order
	order.UserID = userID
	order.PaymentID = paymentID
	order.Status = models.OrderStatusPaymentComplete

	result, err := tx.ExecContext(ctx, "INSERT INTO orders (user_id, payment_id, status) VALUES (?, ?, ?)", order.UserID, order.PaymentID, order.Status)
	if err != nil {
		return nil, err
	}
	order.ID, _ = result.LastInsertId()

	for _, item := range cartItems {
		fmt.Print(item)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *orderRepository) GetOrdersByUserID(ctx context.Context, userID int) ([]*models.Order, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, user_id, payment_id, status FROM orders WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.PaymentID, &order.Status); err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *orderRepository) GetOrderByID(ctx context.Context, userID, orderID int) (*models.Order, error) {
	var order models.Order
	row := r.db.QueryRowContext(ctx, "SELECT id, user_id, payment_id, status FROM orders WHERE id = ? AND user_id = ?", orderID, userID)
	if err := row.Scan(&order.ID, &order.UserID, &order.PaymentID, &order.Status); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &order, nil
}

func (r *orderRepository) CancelOrder(ctx context.Context, userID, orderID int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "UPDATE orders SET status = ? WHERE id = ? AND user_id = ?", models.OrderStatusCancelled, orderID, userID)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *orderRepository) getOrderItems(ctx context.Context, orderID int) ([]models.OrderItem, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT product_id, quantity FROM order_items WHERE order_id = ?", orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderItems []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.ProductID, &item.Quantity); err != nil {
			return nil, err
		}
		orderItems = append(orderItems, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orderItems, nil
}
