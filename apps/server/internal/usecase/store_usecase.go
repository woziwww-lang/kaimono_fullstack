package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/price-comparison/server/internal/domain"
	"github.com/price-comparison/server/internal/query"
)

type StoreRepository interface {
	FindNearby(lat, lon float64, radiusMeters int, limit, offset int) ([]domain.Store, error)
	FindAll(filters query.StoreFilters, limit, offset int, sortField, sortOrder string) ([]domain.Store, error)
	FindByID(id int) (*domain.Store, error)
}

type StoreUsecase struct {
	repo     StoreRepository
	cache    Cache
	cacheTTL time.Duration
}

func NewStoreUsecase(repo StoreRepository, cache Cache, cacheTTL time.Duration) *StoreUsecase {
	return &StoreUsecase{repo: repo, cache: cache, cacheTTL: cacheTTL}
}

func (u *StoreUsecase) List(opts StoreListOptions) ([]domain.Store, error) {
	limit := normalizeLimit(opts.Limit)
	offset := normalizeOffset(opts.Offset)
	sortField, sortOrder := normalizeStoreSort(opts.Sort, opts.UserLocation != nil)

	filters := query.StoreFilters{
		Query:        opts.Query,
		Category:     opts.Category,
		Bounds:       opts.Bounds,
		UserLocation: opts.UserLocation,
	}

	cacheKey := buildStoreCacheKey(filters, limit, offset, sortField, sortOrder)
	if u.cache != nil {
		if cached, err := u.cache.Get(context.Background(), cacheKey); err == nil {
			var stores []domain.Store
			if err := json.Unmarshal([]byte(cached), &stores); err == nil {
				return stores, nil
			}
		}
	}

	stores, err := u.repo.FindAll(filters, limit, offset, sortField, sortOrder)
	if err != nil {
		return nil, err
	}

	if u.cache != nil {
		if payload, err := json.Marshal(stores); err == nil {
			_ = u.cache.Set(context.Background(), cacheKey, string(payload), u.cacheTTL)
		}
	}

	return stores, nil
}

func (u *StoreUsecase) Nearby(opts StoreNearbyOptions) ([]domain.Store, error) {
	if opts.Radius <= 0 {
		return nil, fmt.Errorf("radius must be positive")
	}
	limit := normalizeLimit(opts.Limit)
	offset := normalizeOffset(opts.Offset)
	cacheKey := fmt.Sprintf("stores:nearby:%.5f:%.5f:%d:%d:%d", opts.Latitude, opts.Longitude, opts.Radius, limit, offset)
	if u.cache != nil {
		if cached, err := u.cache.Get(context.Background(), cacheKey); err == nil {
			var stores []domain.Store
			if err := json.Unmarshal([]byte(cached), &stores); err == nil {
				return stores, nil
			}
		}
	}

	stores, err := u.repo.FindNearby(opts.Latitude, opts.Longitude, opts.Radius, limit, offset)
	if err != nil {
		return nil, err
	}

	if u.cache != nil {
		if payload, err := json.Marshal(stores); err == nil {
			_ = u.cache.Set(context.Background(), cacheKey, string(payload), u.cacheTTL)
		}
	}

	return stores, nil
}

func (u *StoreUsecase) GetByID(id int) (*domain.Store, error) {
	if id <= 0 {
		return nil, fmt.Errorf("id must be positive")
	}
	return u.repo.FindByID(id)
}

func normalizeStoreSort(sort Sort, hasLocation bool) (string, string) {
	field := sort.Field
	order := normalizeOrder(sort.Order)

	switch field {
	case "distance":
		if hasLocation {
			return "distance", order
		}
		return "name", order
	case "price":
		return "price", order
	case "created_at":
		return "created_at", order
	case "name":
		fallthrough
	default:
		return "name", order
	}
}

func buildStoreCacheKey(filters query.StoreFilters, limit, offset int, sortField, sortOrder string) string {
	boundsKey := "none"
	if filters.Bounds != nil {
		boundsKey = fmt.Sprintf("%.4f:%.4f:%.4f:%.4f",
			filters.Bounds.MinLat,
			filters.Bounds.MinLon,
			filters.Bounds.MaxLat,
			filters.Bounds.MaxLon,
		)
	}

	locationKey := "none"
	if filters.UserLocation != nil {
		locationKey = fmt.Sprintf("%.4f:%.4f", filters.UserLocation.Lat, filters.UserLocation.Lon)
	}

	return fmt.Sprintf("stores:list:%s:%s:%s:%s:%d:%d:%s:%s",
		filters.Query,
		filters.Category,
		boundsKey,
		locationKey,
		limit,
		offset,
		sortField,
		sortOrder,
	)
}
