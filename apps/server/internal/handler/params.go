package handler

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/price-comparison/server/internal/query"
	"github.com/price-comparison/server/internal/usecase"
)

func parsePagination(c *gin.Context) (limit int, offset int, err error) {
	limitStr := c.DefaultQuery("limit", "")
	offsetStr := c.DefaultQuery("offset", "0")

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return 0, 0, err
		}
	}

	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			return 0, 0, err
		}
	}
	if limit <= 0 {
		limit = usecase.DefaultLimit
	}
	if limit > usecase.MaxLimit {
		limit = usecase.MaxLimit
	}
	if offset < 0 {
		offset = 0
	}

	return limit, offset, nil
}

func parseSort(c *gin.Context) (field, order string) {
	field = c.DefaultQuery("sort", "")
	order = c.DefaultQuery("order", "")
	return field, order
}

func parseBounds(c *gin.Context) (*query.Bounds, error) {
	bbox := c.Query("bbox")
	if bbox == "" {
		return nil, nil
	}

	parts := splitAndTrim(bbox)
	if len(parts) != 4 {
		return nil, strconv.ErrSyntax
	}

	minLon, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return nil, err
	}
	minLat, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return nil, err
	}
	maxLon, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return nil, err
	}
	maxLat, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return nil, err
	}

	return &query.Bounds{
		MinLat: minLat,
		MinLon: minLon,
		MaxLat: maxLat,
		MaxLon: maxLon,
	}, nil
}

func parseUserLocation(c *gin.Context) (*query.GeoPoint, error) {
	latStr := c.Query("user_lat")
	lonStr := c.Query("user_lon")
	if latStr == "" || lonStr == "" {
		return nil, nil
	}
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return nil, err
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		return nil, err
	}
	return &query.GeoPoint{Lat: lat, Lon: lon}, nil
}

func splitAndTrim(value string) []string {
	raw := strings.Split(value, ",")
	out := make([]string, 0, len(raw))
	for _, item := range raw {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
