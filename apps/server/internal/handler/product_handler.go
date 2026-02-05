package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/price-comparison/server/internal/repository"
)

type ProductHandler struct {
	productRepo *repository.ProductRepository
	priceRepo   *repository.PriceRepository
}

func NewProductHandler(productRepo *repository.ProductRepository, priceRepo *repository.PriceRepository) *ProductHandler {
	return &ProductHandler{
		productRepo: productRepo,
		priceRepo:   priceRepo,
	}
}

// GetAllProducts handles GET /api/products
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	products, err := h.productRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"count":    len(products),
	})
}

// GetProductByID handles GET /api/products/:id
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	product, err := h.productRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// SearchProducts handles GET /api/products/search?q=keyword
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "search keyword is required"})
		return
	}

	products, err := h.productRepo.Search(keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"count":    len(products),
	})
}

// GetProductPrices handles GET /api/products/:id/prices
func (h *ProductHandler) GetProductPrices(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	prices, err := h.priceRepo.FindByProductID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"prices": prices,
		"count":  len(prices),
	})
}
