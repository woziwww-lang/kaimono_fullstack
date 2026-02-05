package repository

import (
	"database/sql"
	"fmt"

	"github.com/price-comparison/server/internal/domain"
)

type StoreRepository struct {
	db *sql.DB
}

func NewStoreRepository(db *sql.DB) *StoreRepository {
	return &StoreRepository{db: db}
}

// FindNearby finds stores within a specified radius (in meters) from a given point
func (r *StoreRepository) FindNearby(lat, lon float64, radiusMeters int) ([]domain.Store, error) {
	query := `
		SELECT
			id,
			name,
			address,
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
	`

	rows, err := r.db.Query(query, lat, lon, radiusMeters)
	if err != nil {
		return nil, fmt.Errorf("failed to query nearby stores: %w", err)
	}
	defer rows.Close()

	var stores []domain.Store
	for rows.Next() {
		var store domain.Store
		var distance float64
		err := rows.Scan(
			&store.ID,
			&store.Name,
			&store.Address,
			&store.Latitude,
			&store.Longitude,
			&distance,
			&store.CreatedAt,
			&store.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan store: %w", err)
		}
		store.Distance = &distance
		stores = append(stores, store)
	}

	return stores, nil
}

// FindAll returns all stores
func (r *StoreRepository) FindAll() ([]domain.Store, error) {
	query := `
		SELECT
			id,
			name,
			address,
			ST_Y(location::geometry) as latitude,
			ST_X(location::geometry) as longitude,
			created_at,
			updated_at
		FROM stores
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query stores: %w", err)
	}
	defer rows.Close()

	var stores []domain.Store
	for rows.Next() {
		var store domain.Store
		err := rows.Scan(
			&store.ID,
			&store.Name,
			&store.Address,
			&store.Latitude,
			&store.Longitude,
			&store.CreatedAt,
			&store.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan store: %w", err)
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
			ST_Y(location::geometry) as latitude,
			ST_X(location::geometry) as longitude,
			created_at,
			updated_at
		FROM stores
		WHERE id = $1
	`

	var store domain.Store
	err := r.db.QueryRow(query, id).Scan(
		&store.ID,
		&store.Name,
		&store.Address,
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

	return &store, nil
}
