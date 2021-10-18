package gormquery

import (
	"fmt"
	"strings"

	"github.com/onichandame/gormquery/dto"
	"gorm.io/gorm"
)

func getQuery(f *dto.Filter) (query string, values []interface{}) {
	getFieldQuery := func(field string, condition interface{}) (query string, values []interface{}) {
		if condition == nil {
			query = fmt.Sprintf("%s AND %s IS ?", query, field)
			values = append(values, nil)
		} else if f, ok := condition.(dto.FieldFilter); ok {
			if is, ok := f[dto.Is]; ok {
				query = fmt.Sprintf("%s AND %s IS ?", query, field)
				values = append(values, is)
			}
			if eq, ok := f[dto.Eq]; ok {
				query = fmt.Sprintf("%s AND %s = ?", query, field)
				values = append(values, eq)
			}
			if gt, ok := f[dto.GT]; ok {
				query = fmt.Sprintf("%s AND %s > ?", query, field)
				values = append(values, gt)
			}
			if lt, ok := f[dto.LT]; ok {
				query = fmt.Sprintf("%s AND %s < ?", query, field)
				values = append(values, lt)
			}
			if gte, ok := f[dto.GTE]; ok {
				query = fmt.Sprintf("%s AND %s >= ?", query, field)
				values = append(values, gte)
			}
			if lte, ok := f[dto.LTE]; ok {
				query = fmt.Sprintf("%s AND %s <= ?", query, field)
				values = append(values, lte)
			}
			if in, ok := f[dto.In]; ok {
				query = fmt.Sprintf("%s AND %s In ?", query, field)
				values = append(values, in)
			}
		} else {
			query = fmt.Sprintf("%s AND %s = ?", query, field)
			values = append(values, condition)
		}
		query = strings.TrimSpace(query)
		query = strings.Trim(query, `AND`)
		return query, values
	}
	for _, and := range f.And {
		subQuery, subValues := getQuery(and)
		query = fmt.Sprintf("%s AND (%s)", query, subQuery)
		values = append(values, subValues...)
	}
	for _, or := range f.Or {
		subQuery, subValues := getQuery(or)
		query = fmt.Sprintf("%s OR (%s)", query, subQuery)
		values = append(values, subValues...)
	}
	for _, not := range f.Not {
		subQuery, subValues := getQuery(not)
		query = fmt.Sprintf("%s NOT (%s)", query, subQuery)
		values = append(values, subValues...)
	}
	for field, condition := range f.Fields {
		subQuery, subValues := getFieldQuery(field, condition)
		query = fmt.Sprintf("%s AND (%s)", query, subQuery)
		values = append(values, subValues...)
	}
	query = strings.TrimSpace(query)
	query = strings.Trim(query, `AND`)
	return query, values
}

func Filter(filter *dto.Filter) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q, v := getQuery(filter)
		return db.Where(q, v...)
	}
}
