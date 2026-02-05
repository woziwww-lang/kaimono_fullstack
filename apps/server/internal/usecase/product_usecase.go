package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/price-comparison/server/internal/domain"
)

type ProductRepository interface {
	FindAll(limit, offset int, sortField, sortOrder string) ([]domain.Product, error)
	FindByID(id int) (*domain.Product, error)
	Search(keyword string, limit, offset int, sortField, sortOrder string) ([]domain.Product, error)
	ListCategories() ([]string, error)
}

type ProductUsecase struct {
	repo     ProductRepository
	cache    Cache
	cacheTTL time.Duration
}

func NewProductUsecase(repo ProductRepository, cache Cache, cacheTTL time.Duration) *ProductUsecase {
	return &ProductUsecase{repo: repo, cache: cache, cacheTTL: cacheTTL}
}

func (u *ProductUsecase) List(opts ProductListOptions) ([]domain.Product, error) {
	limit := normalizeLimit(opts.Limit)
	offset := normalizeOffset(opts.Offset)
	sortField, sortOrder := normalizeProductSort(opts.Sort)

	cacheKey := fmt.Sprintf("products:list:%d:%d:%s:%s", limit, offset, sortField, sortOrder)
	if u.cache != nil {
		if cached, err := u.cache.Get(context.Background(), cacheKey); err == nil {
			var products []domain.Product
			if err := json.Unmarshal([]byte(cached), &products); err == nil {
				return products, nil
			}
		}
	}

	products, err := u.repo.FindAll(limit, offset, sortField, sortOrder)
	if err != nil {
		return nil, err
	}

	if u.cache != nil {
		if payload, err := json.Marshal(products); err == nil {
			_ = u.cache.Set(context.Background(), cacheKey, string(payload), u.cacheTTL)
		}
	}

	return products, nil
}

func (u *ProductUsecase) GetByID(id int) (*domain.Product, error) {
	if id <= 0 {
		return nil, fmt.Errorf("id must be positive")
	}
	return u.repo.FindByID(id)
}

func (u *ProductUsecase) Search(opts ProductSearchOptions) ([]domain.Product, error) {
	if opts.Keyword == "" {
		return nil, fmt.Errorf("keyword is required")
	}
	limit := normalizeLimit(opts.Limit)
	offset := normalizeOffset(opts.Offset)
	sortField, sortOrder := normalizeProductSort(opts.Sort)

	cacheKey := fmt.Sprintf("products:search:%s:%d:%d:%s:%s", opts.Keyword, limit, offset, sortField, sortOrder)
	if u.cache != nil {
		if cached, err := u.cache.Get(context.Background(), cacheKey); err == nil {
			var products []domain.Product
			if err := json.Unmarshal([]byte(cached), &products); err == nil {
				return products, nil
			}
		}
	}

	products, err := u.repo.Search(opts.Keyword, limit, offset, sortField, sortOrder)
	if err != nil {
		return nil, err
	}

	if u.cache != nil {
		if payload, err := json.Marshal(products); err == nil {
			_ = u.cache.Set(context.Background(), cacheKey, string(payload), u.cacheTTL)
		}
	}

	return products, nil
}

func (u *ProductUsecase) ListCategories() ([]string, error) {
	cacheKey := "products:categories"
	if u.cache != nil {
		if cached, err := u.cache.Get(context.Background(), cacheKey); err == nil {
			var categories []string
			if err := json.Unmarshal([]byte(cached), &categories); err == nil {
				return categories, nil
			}
		}
	}

	categories, err := u.repo.ListCategories()
	if err != nil {
		return nil, err
	}

	if u.cache != nil {
		if payload, err := json.Marshal(categories); err == nil {
			_ = u.cache.Set(context.Background(), cacheKey, string(payload), u.cacheTTL)
		}
	}

	return categories, nil
}

func normalizeProductSort(sort Sort) (string, string) {
	field := sort.Field
	order := normalizeOrder(sort.Order)

	switch field {
	case "created_at":
		return "created_at", order
	case "name":
		fallthrough
	default:
		return "name", order
	}
}
