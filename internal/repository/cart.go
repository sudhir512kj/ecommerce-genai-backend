package repository

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/sudhir512kj/ecommerce_backend/internal/models"
)

// CartRepository is an interface that defines the methods for interacting with the cart
type CartRepository interface {
	GetCart(ctx context.Context, userID int) (*models.Cart, error)
	GetCartItems(ctx context.Context, userID int) ([]*models.CartItem, error)
	AddToCart(ctx context.Context, userID, productID, quantity int) (*models.Cart, error)
	UpdateCart(ctx context.Context, userID int, cartItems []models.CartItem) (*models.Cart, error)
	DeleteFromCart(ctx context.Context, userID, productID int) (*models.Cart, error)
	SaveForLater(ctx context.Context, userID, productID int) (*models.Cart, error)
	CreateDefaultCart(ctx context.Context, userID int) error
}

type cartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) CartRepository {
	return &cartRepository{
		db: db,
	}
}

func (r *cartRepository) GetCart(ctx context.Context, userID int) (*models.Cart, error) {
	var cart models.Cart
	row := r.db.QueryRowContext(ctx, "SELECT id, user_id, total_amount FROM carts WHERE user_id = "+strconv.Itoa(userID))
	if err := row.Scan(&cart.ID, &cart.UserID, &cart.TotalAmount); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	var cartItems []models.CartItem
	cartItems, err := r.getCartItems(ctx, (int)(cart.ID))
	if err != nil {
		return nil, err
	}
	cart.Items = cartItems

	return &cart, nil
}

func (r *cartRepository) GetCartItems(ctx context.Context, userID int) ([]*models.CartItem, error) {
	query := "SELECT product_id, quantity FROM carts WHERE user_id = " + strconv.Itoa(userID)
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cartItems []*models.CartItem
	for rows.Next() {
		var cartItem models.CartItem
		if err := rows.Scan(&cartItem.ProductID, &cartItem.Quantity); err != nil {
			return nil, err
		}
		cartItems = append(cartItems, &cartItem)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cartItems, nil
}

func (r *cartRepository) AddToCart(ctx context.Context, userID, productID, quantity int) (*models.Cart, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// var cart models.Cart
	// row := tx.QueryRowContext(ctx, "SELECT id, user_id, total_amount FROM carts WHERE user_id = "+strconv.Itoa(userID))
	// if err := row.Scan(&cart.ID, &cart.UserID, &cart.TotalAmount); err != nil {
	// 	if err == sql.ErrNoRows {
	// 		// Create a new cart
	// 		result, err := tx.ExecContext(ctx, "INSERT INTO carts (user_id, total_amount) VALUES ("+strconv.Itoa(userID)+", 0)")
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		cart.ID, _ = result.LastInsertId()
	// 		cart.UserID = userID
	// 	} else {
	// 		return nil, nil
	// 	}
	// }

	var cart models.Cart
	row := tx.QueryRowContext(ctx, "SELECT id, user_id, total_amount FROM carts WHERE user_id = "+strconv.Itoa(userID))
	if err := row.Scan(&cart.ID, &cart.UserID, &cart.TotalAmount); err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	if cart.ID == 0 {
		result, err := tx.ExecContext(ctx, "INSERT INTO carts (user_id, total_amount) VALUES ("+strconv.Itoa(userID)+", 0)")
		if err != nil {
			return nil, err
		}
		// lastInsertedID :=
		cart.ID, _ = result.LastInsertId()
		cart.ID += 1
		cart.UserID = userID
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO cart_items (cart_id, product_id, quantity) VALUES ("+strconv.Itoa((int)(cart.ID))+", "+strconv.Itoa(productID)+", "+strconv.Itoa(quantity)+")")
	if err != nil {
		return nil, err
	}

	cart.Items, err = r.getCartItems(ctx, (int)(cart.ID))
	if err != nil {
		return nil, err
	}

	return &cart, nil
}

func (r *cartRepository) CreateDefaultCart(ctx context.Context, userID int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	var cart models.Cart
	result, err := tx.ExecContext(ctx, "INSERT INTO carts (user_id, total_amount) VALUES ("+strconv.Itoa(userID)+", 0)")
	if err != nil {
		return err
	}

	cart.ID, _ = result.LastInsertId()
	cart.UserID = userID
	cart.TotalAmount = 0

	return nil
}

func (r *cartRepository) UpdateCart(ctx context.Context, userID int, cartItems []models.CartItem) (*models.Cart, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var cart models.Cart
	row := tx.QueryRowContext(ctx, "SELECT id, user_id, total_amount FROM carts WHERE user_id = ?", userID)
	if err := row.Scan(&cart.ID, &cart.UserID, &cart.TotalAmount); err != nil {
		return nil, err
	}

	for _, item := range cartItems {
		if item.Quantity == 0 {
			_, err = tx.ExecContext(ctx, "DELETE FROM cart_items WHERE cart_id = ? AND product_id = ?", cart.ID, item.ProductID)
			if err != nil {
				return nil, err
			}
		} else {
			_, err = tx.ExecContext(ctx, "INSERT INTO cart_items (cart_id, product_id, quantity) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE quantity = ?", cart.ID, item.ProductID, item.Quantity, item.Quantity)
			if err != nil {
				return nil, err
			}
		}
	}

	cart.Items, err = r.getCartItems(ctx, (int)(cart.ID))
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &cart, nil
}

func (r *cartRepository) DeleteFromCart(ctx context.Context, userID, productID int) (*models.Cart, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var cart models.Cart
	row := tx.QueryRowContext(ctx, "SELECT id, user_id, total_amount FROM carts WHERE user_id = ?", userID)
	if err := row.Scan(&cart.ID, &cart.UserID, &cart.TotalAmount); err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM cart_items WHERE cart_id = ? AND product_id = ?", cart.ID, productID)
	if err != nil {
		return nil, err
	}

	cart.Items, err = r.getCartItems(ctx, (int)(cart.ID))
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &cart, nil
}

func (r *cartRepository) SaveForLater(ctx context.Context, userID, productID int) (*models.Cart, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var cart models.Cart
	row := tx.QueryRowContext(ctx, "SELECT id, user_id, total_amount FROM carts WHERE user_id = ?", userID)
	if err := row.Scan(&cart.ID, &cart.UserID, &cart.TotalAmount); err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM cart_items WHERE cart_id = ? AND product_id = ?", cart.ID, productID)
	if err != nil {
		return nil, err
	}

	cart.Items, err = r.getCartItems(ctx, (int)(cart.ID))
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &cart, nil
}

func (r *cartRepository) getCartItems(ctx context.Context, cartID int) ([]models.CartItem, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT product_id, quantity FROM cart_items WHERE cart_id = "+strconv.Itoa(cartID)+"")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cartItems []models.CartItem
	for rows.Next() {
		var item models.CartItem
		if err := rows.Scan(&item.ProductID, &item.Quantity); err != nil {
			return nil, err
		}
		cartItems = append(cartItems, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cartItems, nil
}
