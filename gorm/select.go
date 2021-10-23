package gormquery

import "gorm.io/gorm"

func Select(s []string) func(*gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		return d.Select(s)
	}
}
