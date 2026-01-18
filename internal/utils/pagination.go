package utils

import (
	"fmt"
	"komiku-scraper/internal/models"
)

// PaginationParams holds pagination parameters
type PaginationParams struct {
	Page  int
	Limit int
}

// DefaultPaginationParams returns default pagination (page 1, limit 20)
func DefaultPaginationParams() PaginationParams {
	return PaginationParams{
		Page:  1,
		Limit: 20,
	}
}

// Paginate applies pagination to a slice and returns paginated data with meta
func Paginate(items interface{}, params PaginationParams) (interface{}, models.Meta) {
	// Type assertion to get slice length
	var total int
	var result interface{}

	switch v := items.(type) {
	case []interface{}:
		total = len(v)
		result = paginateSlice(v, params)
	default:
		// If unknown type, return as-is
		return items, models.Meta{}
	}

	totalPages := (total + params.Limit - 1) / params.Limit
	if totalPages == 0 {
		totalPages = 1
	}

	return result, models.Meta{
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}
}

func paginateSlice(items []interface{}, params PaginationParams) []interface{} {
	total := len(items)

	if total == 0 {
		return items
	}

	start := (params.Page - 1) * params.Limit
	if start >= total {
		return []interface{}{}
	}

	end := start + params.Limit
	if end > total {
		end = total
	}

	return items[start:end]
}

// ParsePaginationParams parses page and limit from query parameters
func ParsePaginationParams(pageStr, limitStr string) PaginationParams {
	params := DefaultPaginationParams()

	if pageStr != "" {
		var page int
		if _, err := fmt.Sscanf(pageStr, "%d", &page); err == nil && page > 0 {
			params.Page = page
		}
	}

	if limitStr != "" {
		var limit int
		if _, err := fmt.Sscanf(limitStr, "%d", &limit); err == nil && limit > 0 && limit <= 100 {
			params.Limit = limit
		}
	}

	return params
}
