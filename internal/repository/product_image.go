package repository

import (
	"context"
	"database/sql"

	"github.com/sudhir512kj/ecommerce_backend/internal/models"
)

type ProductImageRepository interface {
	GetImagesByProductID(ctx context.Context, productID int) ([]*models.ProductImage, error)
	CreateProductImage(ctx context.Context, image *models.ProductImage) (*models.ProductImage, error)
	DeleteProductImage(ctx context.Context, id int) error
	DeleteProductImagesByProductID(ctx context.Context, productID int) error
}

type productImageRepository struct {
	db *sql.DB
}

func NewProductImageRepository(db *sql.DB) ProductImageRepository {
	return &productImageRepository{db: db}
}

// Implement the methods for the ProductImageRepository interface
func (r *productImageRepository) GetImagesByProductID(ctx context.Context, productID int) ([]*models.ProductImage, error) {
	query := `SELECT id, product_id, image_url FROM product_images WHERE product_id = $1`

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	images := make([]*models.ProductImage, 0)
	for rows.Next() {
		image := &models.ProductImage{}
		if err := rows.Scan(&image.ID, &image.ProductID, &image.ImageURL); err != nil {
			return nil, err
		}
		images = append(images, image)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return images, nil
}

func (r *productImageRepository) CreateProductImage(ctx context.Context, image *models.ProductImage) (*models.ProductImage, error) {
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

	query := `INSERT INTO product_images (product_id, image_url) VALUES ($1, $2) RETURNING id`
	var createdID int
	if err := tx.QueryRowContext(ctx, query, image.ProductID, image.ImageURL).Scan(&createdID); err != nil {
		return nil, err
	}

	image.ID = createdID
	return image, nil
}

func (r *productImageRepository) UpdateProductImage(ctx context.Context, image *models.ProductImage) error {
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

	query := `UPDATE product_images SET image_url = $1 WHERE id = $2`
	_, err = tx.ExecContext(ctx, query, image.ImageURL, image.ID)
	return err
}

func (r *productImageRepository) DeleteProductImage(ctx context.Context, id int) error {
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

	query := `DELETE FROM product_images WHERE id = $1`
	_, err = tx.ExecContext(ctx, query, id)
	return err
}

func (r *productImageRepository) DeleteProductImagesByProductID(ctx context.Context, productID int) error {
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

	query := `DELETE FROM product_images WHERE product_id = $1`
	_, err = tx.ExecContext(ctx, query, productID)
	return err
}
