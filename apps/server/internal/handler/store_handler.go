package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/price-comparison/server/internal/response"
	"github.com/price-comparison/server/internal/usecase"
)

type StoreHandler struct {
	storeUsecase *usecase.StoreUsecase
	priceUsecase *usecase.PriceUsecase
}

func NewStoreHandler(storeUsecase *usecase.StoreUsecase, priceUsecase *usecase.PriceUsecase) *StoreHandler {
	return &StoreHandler{storeUsecase: storeUsecase, priceUsecase: priceUsecase}
}

// GetNearbyStores handles GET /api/stores/nearby
// Query params: lat (latitude), lon (longitude), radius (meters, default: 5000)
func (h *StoreHandler) GetNearbyStores(c *gin.Context) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")
	radiusStr := c.DefaultQuery("radius", "5000")
	limit, offset, err := parsePagination(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidArgument, "invalid pagination")
		return
	}

	if latStr == "" || lonStr == "" {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidArgument, "lat and lon are required")
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidArgument, "invalid latitude")
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidArgument, "invalid longitude")
		return
	}

	radius, err := strconv.Atoi(radiusStr)
	if err != nil || radius <= 0 {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidArgument, "invalid radius")
		return
	}

	stores, err := h.storeUsecase.Nearby(usecase.StoreNearbyOptions{
		Latitude:   lat,
		Longitude:  lon,
		Radius:     radius,
		Pagination: usecase.Pagination{Limit: limit, Offset: offset},
	})
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.ErrInternal, err.Error())
		return
	}

	response.OK(c, stores, &response.Meta{
		Count:  len(stores),
		Limit:  limit,
		Offset: offset,
	})
}

// GetAllStores handles GET /api/stores
func (h *StoreHandler) GetAllStores(c *gin.Context) {
	limit, offset, err := parsePagination(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidArgument, "invalid pagination")
		return
	}
	sortField, sortOrder := parseSort(c)

	bounds, err := parseBounds(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidArgument, "invalid bbox")
		return
	}

	userLocation, err := parseUserLocation(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidArgument, "invalid user location")
		return
	}

	stores, err := h.storeUsecase.List(usecase.StoreListOptions{
		Pagination: usecase.Pagination{Limit: limit, Offset: offset},
		Sort:       usecase.Sort{Field: sortField, Order: sortOrder},
		Query:      c.Query("q"),
		Category:   c.Query("category"),
		Bounds:     bounds,
		UserLocation: userLocation,
	})
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.ErrInternal, err.Error())
		return
	}

	response.OK(c, stores, &response.Meta{
		Count:  len(stores),
		Limit:  limit,
		Offset: offset,
	})
}

// GetStoreByID handles GET /api/stores/:id
func (h *StoreHandler) GetStoreByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidArgument, "invalid store id")
		return
	}

	store, err := h.storeUsecase.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.ErrInternal, err.Error())
		return
	}

	if store == nil {
		response.Error(c, http.StatusNotFound, response.ErrNotFound, "store not found")
		return
	}

	response.OK(c, store, nil)
}

// GetStorePrices handles GET /api/stores/:id/prices
func (h *StoreHandler) GetStorePrices(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidArgument, "invalid store id")
		return
	}

	limit, offset, err := parsePagination(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidArgument, "invalid pagination")
		return
	}
	sortField, sortOrder := parseSort(c)
	category := c.Query("category")

	prices, err := h.priceUsecase.ListByStore(usecase.StorePriceListOptions{
		StoreID:   id,
		Category:  category,
		Pagination: usecase.Pagination{Limit: limit, Offset: offset},
		Sort:      usecase.Sort{Field: sortField, Order: sortOrder},
	})
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.ErrInternal, err.Error())
		return
	}

	response.OK(c, prices, &response.Meta{
		Count:  len(prices),
		Limit:  limit,
		Offset: offset,
	})
}

// GetStorePriceStats handles GET /api/stores/:id/price-stats
func (h *StoreHandler) GetStorePriceStats(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidArgument, "invalid store id")
		return
	}

	category := c.Query("category")
	query := c.Query("q")
	days := 0
	if daysParam := c.Query("days"); daysParam != "" {
		parsed, err := strconv.Atoi(daysParam)
		if err != nil || parsed <= 0 {
			response.Error(c, http.StatusBadRequest, response.ErrInvalidArgument, "invalid days")
			return
		}
		days = parsed
	}

	stats, err := h.priceUsecase.GetStorePriceStats(usecase.StorePriceStatsOptions{
		StoreID:  id,
		Category: category,
		Query:    query,
		Days:     days,
	})
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.ErrInternal, err.Error())
		return
	}

	response.OK(c, stats, nil)
}
