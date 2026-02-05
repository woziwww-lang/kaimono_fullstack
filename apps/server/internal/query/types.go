package query

type Bounds struct {
	MinLat float64
	MinLon float64
	MaxLat float64
	MaxLon float64
}

type GeoPoint struct {
	Lat float64
	Lon float64
}

type StoreFilters struct {
	Query        string
	Category     string
	Bounds       *Bounds
	UserLocation *GeoPoint
}
