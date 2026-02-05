package domain

import "time"

// Store represents a retail store with geographic location
type Store struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Distance  *float64  `json:"distance,omitempty"` // Distance in meters (only for nearby queries)
	MinPrice  *float64  `json:"min_price,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Product represents a product item
type Product struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	Barcode   string    `json:"barcode"`
	CreatedAt time.Time `json:"created_at"`
}

// Price represents a price record for a product at a store
type Price struct {
	ID         int       `json:"id"`
	StoreID    int       `json:"store_id"`
	ProductID  int       `json:"product_id"`
	Price      float64   `json:"price"`
	Currency   string    `json:"currency"`
	RecordedAt time.Time `json:"recorded_at"`
	CreatedAt  time.Time `json:"created_at"`

	// Joined data
	Store   *Store   `json:"store,omitempty"`
	Product *Product `json:"product,omitempty"`
}

// PriceComparison represents a product with prices from multiple stores
type PriceComparison struct {
	Product      Product `json:"product"`
	Prices       []Price `json:"prices"`
	LowestPrice  float64 `json:"lowest_price"`
	HighestPrice float64 `json:"highest_price"`
	AveragePrice float64 `json:"average_price"`
}

type PriceSummary struct {
	MinPrice *float64 `json:"min_price,omitempty"`
	MaxPrice *float64 `json:"max_price,omitempty"`
	AvgPrice *float64 `json:"avg_price,omitempty"`
	Currency string   `json:"currency,omitempty"`
}

type DailyPriceStats struct {
	Date     time.Time `json:"date"`
	AvgPrice float64   `json:"avg_price"`
	MinPrice float64   `json:"min_price"`
	MaxPrice float64   `json:"max_price"`
	Count    int       `json:"count"`
}

type StorePriceStats struct {
	StoreID  int               `json:"store_id"`
	Category string            `json:"category,omitempty"`
	Query    string            `json:"query,omitempty"`
	Days     int               `json:"days"`
	Summary  PriceSummary      `json:"summary"`
	Daily    []DailyPriceStats `json:"daily"`
}
