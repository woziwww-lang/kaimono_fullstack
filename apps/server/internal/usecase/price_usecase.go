package usecase

import (
	"fmt"

	"github.com/price-comparison/server/internal/domain"
)

type PriceRepository interface {
	FindByProductID(productID int, limit, offset int, sortField, sortOrder string) ([]domain.Price, error)
	FindByStoreID(storeID int, category string, limit, offset int, sortField, sortOrder string) ([]domain.Price, error)
	FindStorePriceStats(storeID int, category string, query string, days int) (domain.StorePriceStats, error)
}

type PriceUsecase struct {
	repo PriceRepository
}

func NewPriceUsecase(repo PriceRepository) *PriceUsecase {
	return &PriceUsecase{repo: repo}
}

func (u *PriceUsecase) ListByProduct(opts PriceListOptions) ([]domain.Price, error) {
	if opts.ProductID <= 0 {
		return nil, fmt.Errorf("product id must be positive")
	}
	limit := normalizeLimit(opts.Limit)
	offset := normalizeOffset(opts.Offset)
	sortField, sortOrder := normalizePriceSort(opts.Sort)
	return u.repo.FindByProductID(opts.ProductID, limit, offset, sortField, sortOrder)
}

func (u *PriceUsecase) ListByStore(opts StorePriceListOptions) ([]domain.Price, error) {
	if opts.StoreID <= 0 {
		return nil, fmt.Errorf("store id must be positive")
	}
	limit := normalizeLimit(opts.Limit)
	offset := normalizeOffset(opts.Offset)
	sortField, sortOrder := normalizePriceSort(opts.Sort)
	return u.repo.FindByStoreID(opts.StoreID, opts.Category, limit, offset, sortField, sortOrder)
}

func (u *PriceUsecase) GetStorePriceStats(opts StorePriceStatsOptions) (domain.StorePriceStats, error) {
	if opts.StoreID <= 0 {
		return domain.StorePriceStats{}, fmt.Errorf("store id must be positive")
	}
	days := opts.Days
	if days <= 0 {
		days = 14
	}
	if days > 60 {
		days = 60
	}
	return u.repo.FindStorePriceStats(opts.StoreID, opts.Category, opts.Query, days)
}

func normalizePriceSort(sort Sort) (string, string) {
	field := sort.Field
	order := normalizeOrder(sort.Order)

	switch field {
	case "recorded_at":
		return "recorded_at", order
	case "price":
		fallthrough
	default:
		return "price", order
	}
}
