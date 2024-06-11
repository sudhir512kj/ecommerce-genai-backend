package models

import "time"

type Product struct {
	ID            int
	Name          string
	Description   string
	Price         float64
	CategoryID    int
	SubCategoryID int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Images        []string
}

type CreateProductRequest struct {
	Name          string   `json:"name" binding:"required"`
	Description   string   `json:"description" binding:"required"`
	Price         float64  `json:"price" binding:"required,gte=0"`
	CategoryID    int      `json:"category_id" binding:"required"`
	SubCategoryID int      `json:"sub_category_id" binding:"required"`
	Images        []string `json:"images" binding:"dive,url,max=5"`
}

type UpdateProductRequest struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Price         float64  `json:"price" binding:"gte=0"`
	CategoryID    int      `json:"category_id"`
	SubCategoryID int      `json:"sub_category_id"`
	Images        []string `json:"images" binding:"dive,url,max=5"`
}
