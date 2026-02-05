package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/price-comparison/server/internal/domain"
	"github.com/price-comparison/server/internal/query"
)

type StoreRepository struct {
	db *sql.DB
}

func NewStoreRepository(db *sql.DB) *StoreRepository {
	return &StoreRepository{db: db}
}

// FindNearby finds stores within a specified radius (in meters) from a given point
func (r *StoreRepository) FindNearby(lat, lon float64, radiusMeters int, limit, offset int) ([]domain.Store, error) {
	query := `
		SELECT
			id,
			name,
			address,
			phone,
			ST_Y(location::geometry) as latitude,
			ST_X(location::geometry) as longitude,
			ST_Distance(location, ST_GeographyFromText('POINT(' || $2 || ' ' || $1 || ')')) as distance,
			created_at,
			updated_at
		FROM stores
		WHERE ST_DWithin(
			location,
			ST_GeographyFromText('POINT(' || $2 || ' ' || $1 || ')'),
			$3
		)
		ORDER BY distance
		LIMIT $4 OFFSET $5
	`

	rows, err := r.db.Query(query, lat, lon, radiusMeters, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query nearby stores: %w", err)
	}
	defer rows.Close()

	var stores []domain.Store
	for rows.Next() {
		var store domain.Store
		var distance float64
		var phone sql.NullString
		err := rows.Scan(
			&store.ID,
			&store.Name,
			&store.Address,
			&phone,
			&store.Latitude,
			&store.Longitude,
			&distance,
			&store.CreatedAt,
			&store.UpdatedAt,
		)
		if phone.Valid {
			store.Phone = phone.String
		}
		if err != nil {
			return nil, fmt.Errorf("failed to scan store: %w", err)
		}
		store.Distance = &distance
		stores = append(stores, store)
	}

	return stores, nil
}

// FindAll returns all stores with filters
func (r *StoreRepository) FindAll(filters query.StoreFilters, limit, offset int, sortField, sortOrder string) ([]domain.Store, error) {
	var args []interface{}
	addArg := func(value interface{}) string {
		args = append(args, value)
		return fmt.Sprintf("$%d", len(args))
	}

	distanceExpr := "NULL"
	if filters.UserLocation != nil {
		lonArg := addArg(filters.UserLocation.Lon)
		latArg := addArg(filters.UserLocation.Lat)
		pointExpr := fmt.Sprintf("ST_SetSRID(ST_MakePoint(%s, %s), 4326)::geography", lonArg, latArg)
		distanceExpr = fmt.Sprintf("ST_Distance(s.location, %s)", pointExpr)
	}

	minPriceExpr := "NULL"
	priceJoin := ""
	{
		categoryClause := ""
		queryClause := ""
		joinProducts := ""
		if filters.Category != "" || filters.Query != "" {
			joinProducts = "JOIN products pr ON pr.id = p.product_id"
		}
		if filters.Category != "" {
			categoryArg := addArg(filters.Category)
			categoryClause = fmt.Sprintf("AND pr.category = %s", categoryArg)
		}
		if filters.Query != "" {
			queryArg := addArg("%" + filters.Query + "%")
			queryClause = fmt.Sprintf("AND pr.name ILIKE %s", queryArg)
		}

		priceJoin = fmt.Sprintf(`
		LEFT JOIN LATERAL (
			SELECT MIN(p.price) AS min_price
			FROM prices p
			%s
			WHERE p.store_id = s.id %s %s
		) price_summary ON true
		`, joinProducts, categoryClause, queryClause)
		minPriceExpr = "price_summary.min_price"
	}

	var conditions []string
	if filters.Bounds != nil {
		minLonArg := addArg(filters.Bounds.MinLon)
		minLatArg := addArg(filters.Bounds.MinLat)
		maxLonArg := addArg(filters.Bounds.MaxLon)
		maxLatArg := addArg(filters.Bounds.MaxLat)
		conditions = append(conditions, fmt.Sprintf(
			"ST_Intersects(s.location::geometry, ST_MakeEnvelope(%s, %s, %s, %s, 4326))",
			minLonArg, minLatArg, maxLonArg, maxLatArg,
		))
	}
	if filters.Category != "" || filters.Query != "" {
		conditions = append(conditions, "price_summary.min_price IS NOT NULL")
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	orderBy := "s.name"
	nulls := ""
	switch sortField {
	case "created_at":
		orderBy = "s.created_at"
	case "distance":
		orderBy = "distance"
		nulls = "NULLS LAST"
	case "price":
		orderBy = "min_price"
		nulls = "NULLS LAST"
	default:
		orderBy = "s.name"
	}

	orderClause := fmt.Sprintf("%s %s", orderBy, sortOrder)
	if nulls != "" {
		orderClause = fmt.Sprintf("%s %s", orderClause, nulls)
	}

	limitArg := addArg(limit)
	offsetArg := addArg(offset)

	query := fmt.Sprintf(`
		SELECT
			s.id,
			s.name,
			s.address,
			s.phone,
			ST_Y(s.location::geometry) as latitude,
			ST_X(s.location::geometry) as longitude,
			%s as distance,
			%s as min_price,
			s.created_at,
			s.updated_at
		FROM stores s
		%s
		%s
		ORDER BY %s
		LIMIT %s OFFSET %s
	`, distanceExpr, minPriceExpr, priceJoin, whereClause, orderClause, limitArg, offsetArg)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query stores: %w", err)
	}
	defer rows.Close()

	var stores []domain.Store
	for rows.Next() {
		var store domain.Store
		var distance sql.NullFloat64
		var minPrice sql.NullFloat64
		var phone sql.NullString
		err := rows.Scan(
			&store.ID,
			&store.Name,
			&store.Address,
			&phone,
			&store.Latitude,
			&store.Longitude,
			&distance,
			&minPrice,
			&store.CreatedAt,
			&store.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan store: %w", err)
		}
		if phone.Valid {
			store.Phone = phone.String
		}
		if distance.Valid {
			store.Distance = &distance.Float64
		}
		if minPrice.Valid {
			store.MinPrice = &minPrice.Float64
		}
		stores = append(stores, store)
	}

	return stores, nil
}

// FindByID finds a store by its ID
func (r *StoreRepository) FindByID(id int) (*domain.Store, error) {
	query := `
		SELECT
			id,
			name,
			address,
			phone,
			ST_Y(location::geometry) as latitude,
			ST_X(location::geometry) as longitude,
			created_at,
			updated_at
		FROM stores
		WHERE id = $1
	`

	var store domain.Store
	var phone sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&store.ID,
		&store.Name,
		&store.Address,
		&phone,
		&store.Latitude,
		&store.Longitude,
		&store.CreatedAt,
		&store.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find store: %w", err)
	}
	if phone.Valid {
		store.Phone = phone.String
	}

	return &store, nil
}
