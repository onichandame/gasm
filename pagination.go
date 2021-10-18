package gormquery

import (
	"github.com/onichandame/gormquery/dto"
	"gorm.io/gorm"
)

func Pagination(pagination *dto.PaginationInputDTO) func(*gorm.DB) *gorm.DB {
	offset := (pagination.Page - 1) * pagination.PerPage
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(int(offset)).Limit(int(pagination.PerPage))
	}
}
