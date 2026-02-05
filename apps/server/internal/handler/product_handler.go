package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/price-comparison/server/internal/response"
	"github.com/price-comparison/server/internal/usecase"
)

type ProductHandler struct {
	productUsecase *usecase.ProductUsecase
	priceUsecase   *usecase.PriceUsecase
}

func NewProductHandler(productUsecase *usecase.ProductUsecase, priceUsecase *usecase.PriceUsecase) *ProductHandler {
	return &ProductHandler{
		productUsecase: productUsecase,
		priceUsecase:   priceUsecase,
	}
}

// GetAllProducts handles GET /api/products
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	limit, offset, err := parsePagination(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidArgument, "invalid pagination")
		return
	}
	sortField, sortOrder := parseSort(c)
	products, err := h.productUsecase.List(usecase.ProductListOptions{
		Pagination: usecase.Pagination{Limit: limit, Offset: offset},
		Sort:       usecase.Sort{Field: sortField, Order: sortOrder},
	})
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.ErrInternal, err.Error())
		return
	}

	response.OK(c, products, &response.Meta{
		Count:  len(products),
		Limit:  limit,
		Offset: offset,
	})
}

// GetProductByID handles GET /api/products/:id
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidArgument, "invalid product id")
		return
	}

	product, err := h.productUsecase.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.ErrInternal, err.Error())
		return
	}

	if product == nil {
		response.Error(c, http.StatusNotFound, response.ErrNotFound, "product not found")
		return
	}

	response.OK(c, product, nil)
}

// SearchProducts handles GET /api/products/search?q=keyword
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidArgument, "search keyword is required")
		return
	}

	limit, offset, err := parsePagination(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidArgument, "invalid pagination")
		return
	}
	sortField, sortOrder := parseSort(c)
	products, err := h.productUsecase.Search(usecase.ProductSearchOptions{
		Keyword:    keyword,
		Pagination: usecase.Pagination{Limit: limit, Offset: offset},
		Sort:       usecase.Sort{Field: sortField, Order: sortOrder},
	})
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.ErrInternal, err.Error())
		return
	}

	response.OK(c, products, &response.Meta{
		Count:  len(products),
		Limit:  limit,
		Offset: offset,
	})
}

// GetCategories handles GET /api/products/categories
func (h *ProductHandler) GetCategories(c *gin.Context) {
	categories, err := h.productUsecase.ListCategories()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.ErrInternal, err.Error())
		return
	}

	response.OK(c, categories, &response.Meta{
		Count: len(categories),
	})
}

// GetProductPrices handles GET /api/products/:id/prices
func (h *ProductHandler) GetProductPrices(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidArgument, "invalid product id")
		return
	}

	limit, offset, err := parsePagination(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.ErrInvalidArgument, "invalid pagination")
		return
	}
	sortField, sortOrder := parseSort(c)
	prices, err := h.priceUsecase.ListByProduct(usecase.PriceListOptions{
		ProductID:  id,
		Pagination: usecase.Pagination{Limit: limit, Offset: offset},
		Sort:       usecase.Sort{Field: sortField, Order: sortOrder},
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
