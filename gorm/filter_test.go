package gormquery_test

import (
	"testing"
	"time"

	"github.com/onichandame/go-crud/core"
	gormquery "github.com/onichandame/go-crud/gorm"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Entity struct {
	Str  *string
	Date *time.Time
	Int  *int
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
		assert.Nil(t, db.Create(&Entity{Str: &str}).Error)
		assert.Nil(t, db.Create(&Entity{Str: &[]string{"zxcv"}[0], Date: &time.Time{}}).Error)
		var filter core.Filter
		filter.Fields = make(map[string]interface{})
		filter.Fields["date"] = nil
		ents := []Entity{}
		assert.Nil(t, db.Model((&Entity{})).Scopes(gormquery.Filter(&filter)).Find(&ents).Error)
		assert.Equal(t, 1, len(ents))
		filter.Fields = make(map[string]interface{})
		filter.Fields["str"] = str
		assert.Nil(t, db.Model((&Entity{})).Scopes(gormquery.Filter(&filter)).Find(&ents).Error)
		assert.Equal(t, 1, len(ents))
		filter.Fields = make(map[string]interface{})
		filter.Fields["int"] = nil
		assert.Nil(t, db.Model((&Entity{})).Scopes(gormquery.Filter(&filter)).Find(&ents).Error)
		assert.Equal(t, 2, len(ents))
	})
	t.Run("is", func(t *testing.T) {
		db := initDB(t)
		assert.Nil(t, db.Create(&Entity{Date: &time.Time{}}).Error)
		var filter core.Filter
		ents := []Entity{}
		filter.Fields = make(map[string]interface{})
		filter.Fields["date"] = map[string]interface{}{"is": nil}
		assert.Nil(t, db.Model(&Entity{}).Scopes(gormquery.Filter(&filter)).Find(&ents).Error)
		assert.Equal(t, 0, len(ents))
		filter.Fields = make(map[string]interface{})
		filter.Fields["str"] = map[string]interface{}{"is": nil}
		assert.Nil(t, db.Model(&Entity{}).Scopes(gormquery.Filter(&filter)).Find(&ents).Error)
		assert.Equal(t, 1, len(ents))
	})
	t.Run("eq", func(t *testing.T) {
		db := initDB(t)
		str := "asdf"
		assert.Nil(t, db.Create(&Entity{Str: &str}).Error)
		var filter core.Filter
		ents := []Entity{}
		filter.Fields = make(map[string]interface{})
		filter.Fields["str"] = map[string]interface{}{"eq": str}
		assert.Nil(t, db.Model(&Entity{}).Scopes(gormquery.Filter(&filter)).Find(&ents).Error)
		assert.Equal(t, 1, len(ents))
		filter.Fields = make(map[string]interface{})
		filter.Fields["int"] = map[string]interface{}{"eq": 1}
		assert.Nil(t, db.Model(&Entity{}).Scopes(gormquery.Filter(&filter)).Find(&ents).Error)
		assert.Equal(t, 0, len(ents))
	})
	t.Run("gt", func(t *testing.T) {
		db := initDB(t)
		i := 10
		assert.Nil(t, db.Create(&Entity{Int: &[]int{i - 1}[0]}).Error)
		assert.Nil(t, db.Create(&Entity{Int: &[]int{i + 1}[0]}).Error)
		assert.Nil(t, db.Create(&Entity{Date: &[]time.Time{time.Now().Add(time.Hour)}[0]}).Error)
		assert.Nil(t, db.Create(&Entity{Date: &[]time.Time{time.Now().Add(-time.Hour)}[0]}).Error)
		var filter core.Filter
		ents := []Entity{}
		filter.Fields = make(map[string]interface{})
		filter.Fields["int"] = map[string]interface{}{"gt": i}
		assert.Nil(t, db.Model(&Entity{}).Scopes(gormquery.Filter(&filter)).Find(&ents).Error)
		assert.Equal(t, 1, len(ents))
		filter.Fields = make(map[string]interface{})
		filter.Fields["date"] = map[string]interface{}{"gt": time.Now()}
		assert.Nil(t, db.Model(&Entity{}).Scopes(gormquery.Filter(&filter)).Find(&ents).Error)
		assert.Equal(t, 1, len(ents))
	})
	t.Run("lt", func(t *testing.T) {
		db := initDB(t)
		i := 10
		assert.Nil(t, db.Create(&Entity{Int: &[]int{i - 1}[0]}).Error)
		assert.Nil(t, db.Create(&Entity{Int: &[]int{i + 1}[0]}).Error)
		assert.Nil(t, db.Create(&Entity{Date: &[]time.Time{time.Now().Add(-time.Hour)}[0]}).Error)
		assert.Nil(t, db.Create(&Entity{Date: &[]time.Time{time.Now().Add(time.Hour)}[0]}).Error)
		var filter core.Filter
		ents := make([]Entity, 0)
		filter.Fields = make(map[string]interface{})
		filter.Fields["int"] = map[string]interface{}{"lt": i}
		assert.Nil(t, db.Model(&Entity{}).Scopes(gormquery.Filter(&filter)).Find(&ents).Error)
		assert.Equal(t, 1, len(ents))
		filter.Fields = make(map[string]interface{})
		filter.Fields["date"] = map[string]interface{}{"lt": time.Now()}
		assert.Nil(t, db.Model(&Entity{}).Scopes(gormquery.Filter(&filter)).Find(&ents).Error)
		assert.Equal(t, 1, len(ents))
	})
	t.Run("gte", func(t *testing.T) {
		db := initDB(t)
		i := 10
		tm := time.Now()
		assert.Nil(t, db.Create(&Entity{Int: &[]int{i - 1}[0]}).Error)
		assert.Nil(t, db.Create(&Entity{Int: &i}).Error)
		assert.Nil(t, db.Create(&Entity{Int: &[]int{i + 1}[0]}).Error)
		assert.Nil(t, db.Create(&Entity{Date: &[]time.Time{tm.Add(time.Hour)}[0]}).Error)
		assert.Nil(t, db.Create(&Entity{Date: &tm}).Error)
		assert.Nil(t, db.Create(&Entity{Date: &[]time.Time{tm.Add(-time.Hour)}[0]}).Error)
		var filter core.Filter
		ents := make([]Entity, 0)
		filter.Fields = make(map[string]interface{})
		filter.Fields["int"] = map[string]interface{}{"gte": i}
		assert.Nil(t, db.Model(&Entity{}).Scopes(gormquery.Filter(&filter)).Find(&ents).Error)
		assert.Equal(t, 2, len(ents))
		filter.Fields = make(map[string]interface{})
		filter.Fields["date"] = map[string]interface{}{"gte": tm}
		assert.Nil(t, db.Model(&Entity{}).Scopes(gormquery.Filter(&filter)).Find(&ents).Error)
		assert.Equal(t, 2, len(ents))
	})
	t.Run("lte", func(t *testing.T) {
		db := initDB(t)
		i := 10
		tm := time.Now()
		assert.Nil(t, db.Create(&Entity{Int: &[]int{i - 1}[0]}).Error)
		assert.Nil(t, db.Create(&Entity{Int: &i}).Error)
		assert.Nil(t, db.Create(&Entity{Int: &[]int{i + 1}[0]}).Error)
		assert.Nil(t, db.Create(&Entity{Date: &[]time.Time{tm.Add(time.Hour)}[0]}).Error)
		assert.Nil(t, db.Create(&Entity{Date: &tm}).Error)
		assert.Nil(t, db.Create(&Entity{Date: &[]time.Time{tm.Add(-time.Hour)}[0]}).Error)
		var filter core.Filter
		ents := make([]Entity, 0)
		filter.Fields = make(map[string]interface{})
		filter.Fields["int"] = map[string]interface{}{"lte": i}
		assert.Nil(t, db.Model(&Entity{}).Scopes(gormquery.Filter(&filter)).Find(&ents).Error)
		assert.Equal(t, 2, len(ents))
		filter.Fields = make(map[string]interface{})
		filter.Fields["date"] = map[string]interface{}{"lte": tm}
		assert.Nil(t, db.Model(&Entity{}).Scopes(gormquery.Filter(&filter)).Find(&ents).Error)
		assert.Equal(t, 2, len(ents))
	})
	t.Run("in", func(t *testing.T) {
		db := initDB(t)
		strs := []string{
			`asdf`,
			`qwer`,
			`zxcv`,
		}
		for _, s := range strs {
			assert.Nil(t, db.Create(&Entity{Str: &s}).Error)
		}
		var filter core.Filter
		ents := make([]Entity, 0)
		filter.Fields = make(map[string]interface{})
		filter.Fields["str"] = map[string]interface{}{"in": strs}
		assert.Nil(t, db.Model(&Entity{}).Scopes(gormquery.Filter(&filter)).Find(&ents).Error)
		assert.Equal(t, len(strs), len(ents))
		filter.Fields = make(map[string]interface{})
		filter.Fields["str"] = map[string]interface{}{"in": strs[1:]}
		assert.Nil(t, db.Model(&Entity{}).Scopes(gormquery.Filter(&filter)).Find(&ents).Error)
		assert.Equal(t, len(strs)-1, len(ents))
	})
}
