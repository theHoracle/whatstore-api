package controllers

import (
	"math"

	"github.com/gofiber/fiber/v2"
)

type PaginationResponse struct {
	Data        interface{} `json:"data"`
	Total       int64       `json:"total"`
	Page        int         `json:"page"`
	PerPage     int         `json:"per_page"`
	TotalPages  int         `json:"total_pages"`
	HasNext     bool        `json:"has_next"`
	HasPrevious bool        `json:"has_previous"`
}

func NewPaginationResponse(data interface{}, total int64, page, perPage int) PaginationResponse {
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return PaginationResponse{
		Data:        data,
		Total:       total,
		Page:        page,
		PerPage:     perPage,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}
}

func paginate(c *fiber.Ctx) (page int, perPage int) {
	page = c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}

	perPage = c.QueryInt("per_page", 10)
	if perPage < 1 {
		perPage = 10
	}

	return page, perPage
}
