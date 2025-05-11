package utils

import (
	"github.com/gofiber/fiber/v2"
	"math"
	"net/http"
	"strconv"
)

func Paginate(total int64, page, perPage int) (map[string]interface{}, error) {
	if perPage <= 0 {
		perPage = 1
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	var nextPage, prevPage *int
	if page < totalPages {
		next := page + 1
		nextPage = &next
	}
	if page > 1 {
		prev := page - 1
		prevPage = &prev
	}

	return map[string]interface{}{
		"current_page":   page,
		"items_per_page": perPage,
		"next_page":      nextPage,
		"previous_page":  prevPage,
		"total_count":    total,
		"total_pages":    totalPages,
	}, nil
}

func ConvertStringToInt64(str string) int64 {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func ConvertInt64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func ConvertInt64ToUint(i int64) uint {
	return uint(i)
}

func ConvertStringToUint(str string) uint {
	i, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		panic(err)
	}
	return uint(i)
}

type ServiceError struct {
	StatusCode int
	Message    string
	Err        error
}

func (e *ServiceError) Error() string {
	return e.Message
}

type CalculateOffsetStruct struct {
	CurrentPage  int
	ItemsPerPage int
	OrderBy      string
	SortBy       string
	Offset       int
}

func CalculateOffset(currentPage, itemsPerPage int, sortBy, orderBy string) CalculateOffsetStruct {
	if orderBy == "" {
		orderBy = "created_at"
	}
	if sortBy != "asc" && sortBy != "desc" {
		sortBy = "desc"
	}

	offset := (currentPage - 1) * itemsPerPage
	if offset < 0 {
		offset = 0
	}

	return CalculateOffsetStruct{
		CurrentPage:  currentPage,
		ItemsPerPage: itemsPerPage,
		OrderBy:      orderBy,
		SortBy:       sortBy,
		Offset:       offset,
	}
}

// BindJson for Fiber
func BindJson(c *fiber.Ctx, request interface{}) *ServiceError {
	if err := c.BodyParser(request); err != nil {
		return &ServiceError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid input",
		}
	}
	return nil
}
