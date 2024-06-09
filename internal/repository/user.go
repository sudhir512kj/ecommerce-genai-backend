package repository

import (
	"context"
	"database/sql"

	"github.com/sudhir512kj/ecommerce_backend/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	CreateOTP(ctx context.Context, otp *models.OTP) error
	GetOTPByUserID(ctx context.Context, userID int) (*models.OTP, error)
	GetOTPByOTP(otp string) (*models.OTP, error)
	DeleteOTP(id int) error
	CreateAddress(ctx context.Context, address *models.Address) error
	GetAddressesByUserID(ctx context.Context, userID int) ([]*models.Address, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	// Implement database operations to create a new user
	query := `
        INSERT INTO users (first_name, last_name, email, password, permissions)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `
	err := r.db.QueryRowContext(ctx, query, user.FirstName, user.LastName, user.Email, user.Password, "{seller}").Scan(&user.ID)
	return err
}

func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	// Implement database operations to update an existing user
	query := `
        UPDATE users
        SET first_name = $1, last_name = $2, email = $3
        WHERE id = $4
    `
	_, err := r.db.ExecContext(ctx, query, user.FirstName, user.LastName, user.Email, user.ID)
	return err
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	// Implement database operations to retrieve a user by email
	query := `
        SELECT id, first_name, last_name, email, password
        FROM users
        WHERE email = $1
    `
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	query := `
        SELECT id, first_name, last_name, email, password
        FROM users
        WHERE id = $1
    `
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) CreateOTP(ctx context.Context, otp *models.OTP) error {
	// Implement database operations to create a new OTP
	query := `
        INSERT INTO otps (user_id, otp, expires_at)
        VALUES ($1, $2, $3)
    `
	_, err := r.db.ExecContext(ctx, query, otp.UserID, otp.OTP, otp.ExpiresAt)
	return err
}

func (r *userRepository) GetOTPByUserID(ctx context.Context, userID int) (*models.OTP, error) {
	// Implement database operations to retrieve an OTP by user ID
	query := `
        SELECT id, user_id, otp, expires_at
        FROM otps
        WHERE user_id = $1
        ORDER BY id DESC
        LIMIT 1
    `
	otp := &models.OTP{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&otp.ID, &otp.UserID, &otp.OTP, &otp.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return otp, nil
}

func (r *userRepository) GetOTPByOTP(otp string) (*models.OTP, error) {
	query := "SELECT id, user_id, otp, expires_at, created_at, updated_at FROM otps WHERE otp = $1"
	var o models.OTP
	err := r.db.QueryRow(query, otp).Scan(&o.ID, &o.UserID, &o.OTP, &o.ExpiresAt, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &o, nil
}

func (r *userRepository) DeleteOTP(id int) error {
	query := "DELETE FROM otps WHERE id = $1"
	_, err := r.db.Exec(query, id)
	return err
}

func (r *userRepository) CreateAddress(ctx context.Context, address *models.Address) error {
	// Implement database operations to create a new address
	query := `
        INSERT INTO addresses (user_id, street, city, state, country, zipcode)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `
	err := r.db.QueryRowContext(ctx, query, address.UserID, address.Street, address.City, address.State, address.Country, address.Zipcode).Scan(&address.ID)
	return err
}

func (r *userRepository) GetAddressesByUserID(ctx context.Context, userID int) ([]*models.Address, error) {
	// Implement database operations to retrieve addresses by user ID
	query := `
        SELECT id, user_id, street, city, state, country, zipcode
        FROM addresses
        WHERE user_id = $1
    `
	var addresses []*models.Address
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		address := &models.Address{}
		if err := rows.Scan(&address.ID, &address.UserID, &address.Street, &address.City, &address.State, &address.Country, &address.Zipcode); err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return addresses, nil
}
