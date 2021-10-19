package gormquery

import (
	"gorm.io/gorm"
)

type PaginationInput struct {
	Page    uint `json:"page" validate:"required,min=1"`
	PerPage uint `json:"per_page" validate:"required,min=1"`
}

func Pagination(pagination *PaginationInput) func(*gorm.DB) *gorm.DB {
	offset := (pagination.Page - 1) * pagination.PerPage
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(int(offset)).Limit(int(pagination.PerPage))
	}
}
