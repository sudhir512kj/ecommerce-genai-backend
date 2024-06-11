package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/sudhir512kj/ecommerce_backend/internal/models"
)

type ProductRepository interface {
	GetAllProducts(ctx context.Context, page, limit int) ([]*models.Product, error)
	GetProductByID(ctx context.Context, id int) (*models.Product, error)
	CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) error
	DeleteProduct(ctx context.Context, id int) error
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetAllProducts(ctx context.Context, page, limit int) ([]*models.Product, error) {
	// Implement the logic to fetch all products with pagination
	offset := (page - 1) * limit
	query := `SELECT id, name, description, price, category_id, subcategory_id, created_at, updated_at 
              FROM products 
              ORDER BY created_at DESC
              LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]*models.Product, 0, limit)
	for rows.Next() {
		product := &models.Product{}
		if err := rows.Scan(
			&product.ID, &product.Name, &product.Description, &product.Price,
			&product.CategoryID, &product.SubCategoryID, &product.CreatedAt, &product.UpdatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) GetProductByID(ctx context.Context, id int) (*models.Product, error) {
	// Implement the logic to fetch a product by ID
	query := `SELECT id, name, description, price, category_id, subcategory_id, created_at, updated_at 
              FROM products 
              WHERE id = $1`

	product := &models.Product{}
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID, &product.Name, &product.Description, &product.Price,
		&product.CategoryID, &product.SubCategoryID, &product.CreatedAt, &product.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return product, nil
}

func (r *productRepository) CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	// Implement the logic to create a new product
	query := `INSERT INTO products (name, description, price, category_id, subcategory_id, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7) 
              RETURNING id, created_at, updated_at`

	var createdProduct models.Product
	if err := r.db.QueryRowContext(ctx, query,
		product.Name, product.Description, product.Price, product.CategoryID, product.SubCategoryID,
		product.CreatedAt, product.UpdatedAt,
	).Scan(&createdProduct.ID, &createdProduct.CreatedAt, &createdProduct.UpdatedAt); err != nil {
		return nil, err
	}

	createdProduct.Name = product.Name
	createdProduct.Description = product.Description
	createdProduct.Price = product.Price
	createdProduct.CategoryID = product.CategoryID
	createdProduct.SubCategoryID = product.SubCategoryID

	return &createdProduct, nil
}

func (r *productRepository) UpdateProduct(ctx context.Context, product *models.Product) error {
	// Implement the logic to update an existing product
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

	query := `UPDATE products 
              SET name = $1, description = $2, price = $3, category_id = $4, subcategory_id = $5, updated_at = $6
              WHERE id = $7`
	_, err = tx.ExecContext(ctx, query,
		product.Name, product.Description, product.Price, product.CategoryID, product.SubCategoryID, time.Now(), product.ID)
	if err != nil {
		return err
	}

	// Delete existing product images
	_, err = tx.ExecContext(ctx, "DELETE FROM product_images WHERE product_id = $1", product.ID)
	if err != nil {
		return err
	}

	// Insert new product images
	for _, imageURL := range product.Images {
		_, err = tx.ExecContext(ctx, "INSERT INTO product_images (product_id, image_url) VALUES ($1, $2)", product.ID, imageURL)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *productRepository) DeleteProduct(ctx context.Context, id int) error {
	// Implement the logic to delete a product
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

	// Delete product images
	_, err = tx.ExecContext(ctx, "DELETE FROM product_images WHERE product_id = $1", id)
	if err != nil {
		return err
	}

	// Delete the product
	_, err = tx.ExecContext(ctx, "DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
