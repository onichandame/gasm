package gormquery

import (
	"github.com/onichandame/go-crud/core"
	"gorm.io/gorm"
)

func Pagination(pagination *core.Pagination) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pagination == nil {
			return db
		} else {
			offset := (pagination.Page - 1) * pagination.PerPage
			return db.Offset(int(offset)).Limit(int(pagination.PerPage))
		}

	}
}
