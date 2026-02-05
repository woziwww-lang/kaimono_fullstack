package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/price-comparison/server/internal/domain"
)

type PriceRepository struct {
	db *sql.DB
}

func NewPriceRepository(db *sql.DB) *PriceRepository {
	return &PriceRepository{db: db}
}

// FindByProductID finds all prices for a specific product
func (r *PriceRepository) FindByProductID(productID int, limit, offset int, sortField, sortOrder string) ([]domain.Price, error) {
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
		ORDER BY %s %s
		LIMIT $2 OFFSET $3
	`

	orderedQuery := fmt.Sprintf(query, sortField, sortOrder)
	rows, err := r.db.Query(orderedQuery, productID, limit, offset)
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

// FindByStoreID finds all prices for a specific store (optionally filtered by category)
func (r *PriceRepository) FindByStoreID(storeID int, category string, limit, offset int, sortField, sortOrder string) ([]domain.Price, error) {
	orderBy := "p.price"
	if sortField == "recorded_at" {
		orderBy = "p.recorded_at"
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
		WHERE p.store_id = $1
	`

	args := []interface{}{storeID}
	if category != "" {
		query += " AND pr.category = $2"
		args = append(args, category)
		query += fmt.Sprintf(" ORDER BY %s %s LIMIT $3 OFFSET $4", orderBy, sortOrder)
		args = append(args, limit, offset)
	} else {
		query += fmt.Sprintf(" ORDER BY %s %s LIMIT $2 OFFSET $3", orderBy, sortOrder)
		args = append(args, limit, offset)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query prices by store: %w", err)
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

func (r *PriceRepository) FindStorePriceStats(storeID int, category string, query string, days int) (domain.StorePriceStats, error) {
	if days <= 0 {
		days = 14
	}

	args := []interface{}{storeID, days}
	where := "WHERE p.store_id = $1 AND p.recorded_at >= NOW() - ($2 * INTERVAL '1 day')"
	if category != "" {
		args = append(args, category)
		where += fmt.Sprintf(" AND pr.category = $%d", len(args))
	}
	if query != "" {
		args = append(args, "%"+query+"%")
		where += fmt.Sprintf(" AND pr.name ILIKE $%d", len(args))
	}

	withClause := fmt.Sprintf(`
		WITH filtered AS (
			SELECT p.price, p.currency, p.recorded_at::date AS day
			FROM prices p
			JOIN products pr ON pr.id = p.product_id
			%s
		)
	`, where)

	summaryQuery := withClause + `
		SELECT MIN(price), MAX(price), AVG(price), MIN(currency)
		FROM filtered
	`

	var minPrice sql.NullFloat64
	var maxPrice sql.NullFloat64
	var avgPrice sql.NullFloat64
	var currency sql.NullString
	if err := r.db.QueryRow(summaryQuery, args...).Scan(&minPrice, &maxPrice, &avgPrice, &currency); err != nil {
		return domain.StorePriceStats{}, fmt.Errorf("failed to query price summary: %w", err)
	}

	summary := domain.PriceSummary{Currency: currency.String}
	if minPrice.Valid {
		summary.MinPrice = &minPrice.Float64
	}
	if maxPrice.Valid {
		summary.MaxPrice = &maxPrice.Float64
	}
	if avgPrice.Valid {
		summary.AvgPrice = &avgPrice.Float64
	}

	dailyQuery := withClause + `
		SELECT day, AVG(price), MIN(price), MAX(price), COUNT(*)
		FROM filtered
		GROUP BY day
		ORDER BY day
	`

	rows, err := r.db.Query(dailyQuery, args...)
	if err != nil {
		return domain.StorePriceStats{}, fmt.Errorf("failed to query daily stats: %w", err)
	}
	defer rows.Close()

	var daily []domain.DailyPriceStats
	for rows.Next() {
		var day time.Time
		var avg float64
		var min float64
		var max float64
		var count int
		if err := rows.Scan(&day, &avg, &min, &max, &count); err != nil {
			return domain.StorePriceStats{}, fmt.Errorf("failed to scan daily stats: %w", err)
		}
		daily = append(daily, domain.DailyPriceStats{
			Date:     day,
			AvgPrice: avg,
			MinPrice: min,
			MaxPrice: max,
			Count:    count,
		})
	}

	return domain.StorePriceStats{
		StoreID:  storeID,
		Category: category,
		Query:    query,
		Days:     days,
		Summary:  summary,
		Daily:    daily,
	}, nil
}
