package gormquery

import (
	"fmt"

	"github.com/onichandame/go-crud/core"
	"gorm.io/gorm"
)

func Sort(sort core.Sort) func(*gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		if sort == nil {
			return d
		} else {
			for field, order := range sort {
				d = d.Order(fmt.Sprintf("%s %v", field, order))
			}
			return d
		}
	}
}
