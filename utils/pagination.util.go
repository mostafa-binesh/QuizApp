package utils

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"math"
)

type Pagination struct {
	Limit      int   `json:"limit" query:"limit"` // per page
	Page       int   `json:"page" query:"page"`   // current page
	TotalRows  int64 `json:"totalRows"`           // total records
	TotalPages int   `json:"totalPages"`          // total pages
	// Rows       interface{} `json:"rows"`
}

func ParsedPagination(c *fiber.Ctx) *Pagination {
	pagination := new(Pagination)
	if err := c.QueryParser(pagination); err != nil {
		// ResErr(c, err.Error())
		panic("FAILED TO PARSE PAGINATION QUERY")
	}
	return pagination
}
func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}
func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}
func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}
func Paginate(value interface{}, pagination *Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var totalRows int64
		// D.DB().Model(value).Count(&totalRows)
		countDBSession := db.Session(&gorm.Session{Initialized: true})
		countDBSession.Model(value).Count(&totalRows)
		pagination.GetLimit()
		pagination.TotalRows = totalRows
		totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
		pagination.TotalPages = totalPages
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}
