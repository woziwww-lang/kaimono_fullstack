package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/price-comparison/server/internal/repository"
)

type StoreHandler struct {
	storeRepo *repository.StoreRepository
}

func NewStoreHandler(storeRepo *repository.StoreRepository) *StoreHandler {
	return &StoreHandler{storeRepo: storeRepo}
}

// GetNearbyStores handles GET /api/stores/nearby
// Query params: lat (latitude), lon (longitude), radius (meters, default: 5000)
func (h *StoreHandler) GetNearbyStores(c *gin.Context) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")
	radiusStr := c.DefaultQuery("radius", "5000")

	if latStr == "" || lonStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lat and lon are required"})
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid latitude"})
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid longitude"})
		return
	}

	radius, err := strconv.Atoi(radiusStr)
	if err != nil || radius <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid radius"})
		return
	}

	stores, err := h.storeRepo.FindNearby(lat, lon, radius)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stores": stores,
		"count":  len(stores),
	})
}

// GetAllStores handles GET /api/stores
func (h *StoreHandler) GetAllStores(c *gin.Context) {
	stores, err := h.storeRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stores": stores,
		"count":  len(stores),
	})
}

// GetStoreByID handles GET /api/stores/:id
func (h *StoreHandler) GetStoreByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid store id"})
		return
	}

	store, err := h.storeRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if store == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "store not found"})
		return
	}

	c.JSON(http.StatusOK, store)
}
