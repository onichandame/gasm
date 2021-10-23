package gormquery

import (
	"fmt"
	"strings"

	"github.com/onichandame/go-crud/core"
	"gorm.io/gorm"
)

func getQuery(f *core.Filter) (query string, values []interface{}) {
	getFieldQuery := func(field string, condition interface{}) (query string, values []interface{}) {
		if condition == nil {
			query = fmt.Sprintf("%s AND %s IS ?", query, field)
			values = append(values, nil)
		} else if f, ok := condition.(map[string]interface{}); ok {
			if is, ok := f["is"]; ok {
				query = fmt.Sprintf("%s AND %s IS ?", query, field)
				values = append(values, is)
			}
			if eq, ok := f["eq"]; ok {
				query = fmt.Sprintf("%s AND %s = ?", query, field)
				values = append(values, eq)
			}
			if gt, ok := f["gt"]; ok {
				query = fmt.Sprintf("%s AND %s > ?", query, field)
				values = append(values, gt)
			}
			if lt, ok := f["lt"]; ok {
				query = fmt.Sprintf("%s AND %s < ?", query, field)
				values = append(values, lt)
			}
			if gte, ok := f["gte"]; ok {
				query = fmt.Sprintf("%s AND %s >= ?", query, field)
				values = append(values, gte)
			}
			if lte, ok := f["lte"]; ok {
				query = fmt.Sprintf("%s AND %s <= ?", query, field)
				values = append(values, lte)
			}
			if in, ok := f["in"]; ok {
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

func Filter(filter *core.Filter) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter == nil {
			return db
		} else {
			q, v := getQuery(filter)
			return db.Where(q, v...)
		}
	}
}
