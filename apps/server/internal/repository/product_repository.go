package repository

import (
	"database/sql"
	"fmt"

	"github.com/price-comparison/server/internal/domain"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// FindAll returns all products
func (r *ProductRepository) FindAll(limit, offset int, sortField, sortOrder string) ([]domain.Product, error) {
	query := `
		SELECT id, name, category, barcode, created_at
		FROM products
		ORDER BY %s %s
		LIMIT $1 OFFSET $2
	`

	orderedQuery := fmt.Sprintf(query, sortField, sortOrder)
	rows, err := r.db.Query(orderedQuery, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var product domain.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Category,
			&product.Barcode,
			&product.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	return products, nil
}

// FindByID finds a product by its ID
func (r *ProductRepository) FindByID(id int) (*domain.Product, error) {
	query := `
		SELECT id, name, category, barcode, created_at
		FROM products
		WHERE id = $1
	`

	var product domain.Product
	err := r.db.QueryRow(query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Category,
		&product.Barcode,
		&product.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find product: %w", err)
	}

	return &product, nil
}

// Search searches products by name
func (r *ProductRepository) Search(keyword string, limit, offset int, sortField, sortOrder string) ([]domain.Product, error) {
	query := `
		SELECT id, name, category, barcode, created_at
		FROM products
		WHERE name ILIKE '%' || $1 || '%'
		ORDER BY %s %s
		LIMIT $2 OFFSET $3
	`

	orderedQuery := fmt.Sprintf(query, sortField, sortOrder)
	rows, err := r.db.Query(orderedQuery, keyword, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var product domain.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Category,
			&product.Barcode,
			&product.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	return products, nil
}

// ListCategories returns distinct product categories
func (r *ProductRepository) ListCategories() ([]string, error) {
	query := `
		SELECT DISTINCT category
		FROM products
		WHERE category IS NOT NULL AND category <> ''
		ORDER BY category
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}

	return categories, nil
}
