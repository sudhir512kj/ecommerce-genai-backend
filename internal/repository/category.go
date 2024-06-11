package repository

import (
	"context"
	"database/sql"

	"github.com/sudhir512kj/ecommerce_backend/internal/models"
)

type CategoryRepository interface {
	GetAllCategories(ctx context.Context) ([]*models.Category, error)
	GetCategoryByID(ctx context.Context, id int) (*models.Category, error)
	UpdateCategory(ctx context.Context, category *models.Category) error
	DeleteCategory(ctx context.Context, id int) error

	GetAllSubcategoriesByCategoryID(ctx context.Context, categoryID int) ([]*models.Subcategory, error)
	GetSubcategoryByID(ctx context.Context, id int) (*models.Subcategory, error)
	CreateSubcategory(ctx context.Context, subcategory *models.Subcategory) (*models.Subcategory, error)
	UpdateSubcategory(ctx context.Context, subcategory *models.Subcategory) error
	DeleteSubcategory(ctx context.Context, id int) error
}

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

// Implement the methods for the CategoryRepository interface
func (r *categoryRepository) GetAllCategories(ctx context.Context) ([]*models.Category, error) {
	query := `SELECT id, name FROM categories`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]*models.Category, 0)
	for rows.Next() {
		category := &models.Category{}
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepository) GetCategoryByID(ctx context.Context, id int) (*models.Category, error) {
	query := `SELECT id, name FROM categories WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, id)

	category := &models.Category{}
	if err := row.Scan(&category.ID, &category.Name); err != nil {
		return nil, err
	}

	return category, nil
}

func (r *categoryRepository) UpdateCategory(ctx context.Context, category *models.Category) error {
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

	query := `UPDATE categories SET name = $1 WHERE id = $2`
	_, err = tx.ExecContext(ctx, query, category.Name, category.ID)
	if err != nil {
		return err
	}

	// Update subcategories associated with the category
	_, err = tx.ExecContext(ctx, "UPDATE subcategories SET category_id = $1 WHERE category_id = $2", category.ID, category.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *categoryRepository) DeleteCategory(ctx context.Context, id int) error {
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

	// Delete subcategories associated with the category
	if _, err := tx.ExecContext(ctx, "DELETE FROM subcategories WHERE category_id = ?", id); err != nil {
		return err
	}

	// Delete the category
	if _, err := tx.ExecContext(ctx, "DELETE FROM categories WHERE id = ?", id); err != nil {
		return err
	}

	return nil
}

func (r *categoryRepository) GetAllSubcategoriesByCategoryID(ctx context.Context, categoryID int) ([]*models.Subcategory, error) {
	query := `SELECT id, name, category_id FROM subcategories WHERE category_id = $1`

	rows, err := r.db.QueryContext(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subcategories := make([]*models.Subcategory, 0)
	for rows.Next() {
		subcategory := &models.Subcategory{}
		if err := rows.Scan(&subcategory.ID, &subcategory.Name, &subcategory.CategoryID); err != nil {
			return nil, err
		}
		subcategories = append(subcategories, subcategory)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subcategories, nil
}

func (r *categoryRepository) GetSubcategoryByID(ctx context.Context, id int) (*models.Subcategory, error) {
	query := `SELECT id, name, category_id FROM subcategories WHERE id = $1`

	subcategory := &models.Subcategory{}
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&subcategory.ID, &subcategory.Name, &subcategory.CategoryID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return subcategory, nil
}

func (r *categoryRepository) CreateSubcategory(ctx context.Context, subcategory *models.Subcategory) (*models.Subcategory, error) {
	query := `INSERT INTO subcategories (name, category_id) VALUES ($1, $2) RETURNING id`

	var createdID int
	if err := r.db.QueryRowContext(ctx, query, subcategory.Name, subcategory.CategoryID).Scan(&createdID); err != nil {
		return nil, err
	}

	subcategory.ID = createdID
	return subcategory, nil
}

func (r *categoryRepository) UpdateSubcategory(ctx context.Context, subcategory *models.Subcategory) error {
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

	query := `UPDATE subcategories SET name = $1, category_id = $2 WHERE id = $3`
	_, err = tx.ExecContext(ctx, query, subcategory.Name, subcategory.CategoryID, subcategory.ID)
	if err != nil {
		return err
	}

	// Update any products associated with the subcategory
	_, err = tx.ExecContext(ctx, "UPDATE products SET sub_category_id = $1 WHERE sub_category_id = $2", subcategory.ID, subcategory.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *categoryRepository) DeleteSubcategory(ctx context.Context, id int) error {
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

	// Delete the subcategory
	if _, err := tx.ExecContext(ctx, "DELETE FROM subcategories WHERE id = ?", id); err != nil {
		return err
	}

	return nil
}
