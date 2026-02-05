package usecase

import "github.com/price-comparison/server/internal/query"

const (
	DefaultLimit = 20
	MaxLimit     = 100
)

type Pagination struct {
	Limit  int
	Offset int
}

type Sort struct {
	Field string
	Order string
}

type StoreListOptions struct {
	Pagination
	Sort
	Query        string
	Category     string
	Bounds       *query.Bounds
	UserLocation *query.GeoPoint
}

type StoreNearbyOptions struct {
	Latitude  float64
	Longitude float64
	Radius    int
	Pagination
}

type ProductListOptions struct {
	Pagination
	Sort
}

type ProductSearchOptions struct {
	Keyword string
	Pagination
	Sort
}

type PriceListOptions struct {
	ProductID int
	Pagination
	Sort
}

type StorePriceListOptions struct {
	StoreID  int
	Category string
	Pagination
	Sort
}

type StorePriceStatsOptions struct {
	StoreID  int
	Category string
	Query    string
	Days     int
}
