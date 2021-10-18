package gormquery_test

import (
	"testing"
	"time"

	"github.com/onichandame/gormquery"
	"github.com/onichandame/gormquery/dto"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Entity struct {
	Str  string
	Date *time.Time
	Int  int
}

func TestFilter(t *testing.T) {
	initDB := func(t *testing.T) *gorm.DB {
		db, err := gorm.Open(sqlite.Open(":memory:"))
		assert.Nil(t, err)
		assert.Nil(t, db.AutoMigrate(&Entity{}))
		return db
	}
	t.Run("value", func(t *testing.T) {
		db := initDB(t)
		str := "asdf"
		assert.Nil(t, db.Create(&Entity{Str: str}).Error)
		assert.Nil(t, db.Create(&Entity{Str: "zxcv", Date: &time.Time{}}).Error)
		var filter dto.Filter
		filter.Fields = make(map[string]interface{})
		filter.Fields["date"] = nil
		var count int64
		assert.Nil(t, db.Model((&Entity{})).Scopes(gormquery.Filter(&filter)).Count(&count).Error)
		assert.Equal(t, 1, int(count))
		filter.Fields = make(map[string]interface{})
		filter.Fields["Str"] = str
		assert.Nil(t, db.Model((&Entity{})).Scopes(gormquery.Filter(&filter)).Count(&count).Error)
		assert.Equal(t, 1, int(count))
	})
	t.Run("is", func(t *testing.T) {
		db := initDB(t)
		assert.Nil(t, db.Create(&Entity{Str: "asdf"}).Error)
		assert.Nil(t, db.Create(&Entity{Date: &time.Time{}}).Error)
		var filter dto.Filter
		filter.Fields = make(map[string]interface{})
		filter.Fields["date"] = dto.FieldFilter{dto.Is: nil}
		var count int64
		assert.Nil(t, db.Model(&Entity{}).Scopes(gormquery.Filter(&filter)).Count(&count).Error)
		assert.Equal(t, int64(1), count)
	})
	t.Run("eq", func(t *testing.T) {
		db := initDB(t)
		str := "asdf"
		assert.Nil(t, db.Create(&Entity{Str: str}).Error)
		assert.Nil(t, db.Create(&Entity{Str: "qwer"}).Error)
		var filter dto.Filter
		filter.Fields = make(map[string]interface{})
		filter.Fields["str"] = dto.FieldFilter{dto.Eq: str}
		var count int64
		assert.Nil(t, db.Model(&Entity{}).Scopes(gormquery.Filter(&filter)).Count(&count).Error)
		assert.Equal(t, int64(1), count)
	})
	t.Run("gt", func(t *testing.T) {
		db := initDB(t)
		i := 10
		assert.Nil(t, db.Create(&Entity{Int: i - 1}).Error)
		assert.Nil(t, db.Create(&Entity{Int: i + 1}).Error)
		assert.Nil(t, db.Create(&Entity{Date: &time.Time{}}).Error)
		var filter dto.Filter
		filter.Fields = make(map[string]interface{})
		filter.Fields["int"] = dto.FieldFilter{dto.GT: i}
		var count int64
		assert.Nil(t, db.Model(&Entity{}).Scopes(gormquery.Filter(&filter)).Count(&count).Error)
		assert.Equal(t, int64(1), count)
		filter.Fields = make(map[string]interface{})
		filter.Fields["date"] = dto.FieldFilter{dto.GT: time.Now()}
		assert.Nil(t, db.Model(&Entity{}).Scopes(gormquery.Filter(&filter)).Count(&count).Error)
		assert.Equal(t, int64(0), count)
	})
	t.Run("lt", func(t *testing.T) {
		db := initDB(t)
		i := 10
		assert.Nil(t, db.Create(&Entity{Int: i - 1}).Error)
		assert.Nil(t, db.Create(&Entity{Int: i + 1}).Error)
		assert.Nil(t, db.Create(&Entity{Date: &time.Time{}}).Error)
		var filter dto.Filter
		filter.Fields = make(map[string]interface{})
		filter.Fields["int"] = dto.FieldFilter{dto.LT: i}
		var count int64
		assert.Nil(t, db.Model(&Entity{}).Scopes(gormquery.Filter(&filter)).Count(&count).Error)
		assert.Equal(t, int64(1), count)
		filter.Fields = make(map[string]interface{})
		filter.Fields["date"] = dto.FieldFilter{dto.LT: time.Now()}
		assert.Nil(t, db.Model(&Entity{}).Scopes(gormquery.Filter(&filter)).Count(&count).Error)
		assert.Equal(t, int64(1), count)
	})
}
