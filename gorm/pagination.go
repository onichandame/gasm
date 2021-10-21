package gormquery

import (
	"github.com/onichandame/go-crud/core"
	"gorm.io/gorm"
)

func Pagination(pagination *core.Pagination) func(*gorm.DB) *gorm.DB {
	offset := (pagination.Page - 1) * pagination.PerPage
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(int(offset)).Limit(int(pagination.PerPage))
	}
}
