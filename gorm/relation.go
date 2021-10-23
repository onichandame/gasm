package gormquery

import (
	"gorm.io/gorm"
)

func Relation(relations []string) func(*gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		for _, r := range relations {
			d = d.Preload(r)
		}
		return d
	}
}
