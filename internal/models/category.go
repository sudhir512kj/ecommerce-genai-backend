package models

type Category struct {
	ID   int
	Name string
}

type Subcategory struct {
	ID         int
	Name       string
	CategoryID int
}
