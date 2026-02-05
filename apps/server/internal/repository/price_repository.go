package repository

import (
	"database/sql"
	"fmt"

	"github.com/price-comparison/server/internal/domain"
)

type PriceRepository struct {
	db *sql.DB
}

func NewPriceRepository(db *sql.DB) *PriceRepository {
	return &PriceRepository{db: db}
}

// FindByProductID finds all prices for a specific product
func (r *PriceRepository) FindByProductID(productID int) ([]domain.Price, error) {
	query := `
		SELECT
			p.id,
			p.store_id,
			p.product_id,
			p.price,
			p.currency,
			p.recorded_at,
			p.created_at,
			s.id,
			s.name,
			s.address,
			ST_Y(s.location::geometry) as latitude,
			ST_X(s.location::geometry) as longitude,
			s.created_at,
			s.updated_at
		FROM prices p
		INNER JOIN stores s ON p.store_id = s.id
		WHERE p.product_id = $1
		ORDER BY p.price ASC
	`

	rows, err := r.db.Query(query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to query prices: %w", err)
	}
	defer rows.Close()

	var prices []domain.Price
	for rows.Next() {
		var price domain.Price
		var store domain.Store

		err := rows.Scan(
			&price.ID,
			&price.StoreID,
			&price.ProductID,
			&price.Price,
			&price.Currency,
			&price.RecordedAt,
			&price.CreatedAt,
			&store.ID,
			&store.Name,
			&store.Address,
			&store.Latitude,
			&store.Longitude,
			&store.CreatedAt,
			&store.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan price: %w", err)
		}

		price.Store = &store
		prices = append(prices, price)
	}

	return prices, nil
}

// FindRecentByStoreIDs finds recent prices for multiple stores
func (r *PriceRepository) FindRecentByStoreIDs(storeIDs []int, limit int) ([]domain.Price, error) {
	if len(storeIDs) == 0 {
		return []domain.Price{}, nil
	}

	query := `
		SELECT
			p.id,
			p.store_id,
			p.product_id,
			p.price,
			p.currency,
			p.recorded_at,
			p.created_at,
			pr.id,
			pr.name,
			pr.category,
			pr.barcode
		FROM prices p
		INNER JOIN products pr ON p.product_id = pr.id
		WHERE p.store_id = ANY($1)
		ORDER BY p.recorded_at DESC
		LIMIT $2
	`

	rows, err := r.db.Query(query, storeIDs, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query prices: %w", err)
	}
	defer rows.Close()

	var prices []domain.Price
	for rows.Next() {
		var price domain.Price
		var product domain.Product

		err := rows.Scan(
			&price.ID,
			&price.StoreID,
			&price.ProductID,
			&price.Price,
			&price.Currency,
			&price.RecordedAt,
			&price.CreatedAt,
			&product.ID,
			&product.Name,
			&product.Category,
			&product.Barcode,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan price: %w", err)
		}

		price.Product = &product
		prices = append(prices, price)
	}

	return prices, nil
}
