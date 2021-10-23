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

func TestSort(t *testing.T) {
	type Entity struct {
		gorm.Model
		Int  int
		Date *time.Time
	}
	db, err := gorm.Open(sqlite.Open(":memory:"))
	assert.Nil(t, err)
	assert.Nil(t, db.AutoMigrate(&Entity{}))
	baseTime := time.Now()
	totalCount := 100
	for i := 0; i < totalCount; i++ {
		date := baseTime.Add(time.Second * time.Duration(i))
		assert.Nil(t, db.Create(&Entity{Int: i, Date: &date}).Error)
	}
	t.Run("number asc", func(t *testing.T) {
		var ents []Entity
		assert.Nil(t, db.Scopes(gormquery.Sort(core.Sort{"int": core.ASC})).Find(&ents).Error)
		var prev int
		for _, ent := range ents {
			assert.GreaterOrEqual(t, ent.Int, prev)
			prev = ent.Int
		}
	})
	t.Run("number desc", func(t *testing.T) {
		var ents []Entity
		assert.Nil(t, db.Scopes(gormquery.Sort(core.Sort{"int": core.DESC})).Find(&ents).Error)
		prev := 1000
		for _, ent := range ents {
			assert.LessOrEqual(t, ent.Int, prev)
			prev = ent.Int
		}
	})
	t.Run("date asc", func(t *testing.T) {
		var ents []Entity
		assert.Nil(t, db.Scopes(gormquery.Sort(core.Sort{"date": core.ASC})).Find(&ents).Error)
		prev := baseTime.Add(-time.Hour)
		for _, ent := range ents {
			assert.True(t, ent.Date.After(prev))
			prev = *ent.Date
		}
	})
	t.Run("date desc", func(t *testing.T) {
		var ents []Entity
		assert.Nil(t, db.Scopes(gormquery.Sort(core.Sort{"date": core.DESC})).Find(&ents).Error)
		prev := time.Now().Add(time.Hour * 100)
		for _, ent := range ents {
			assert.True(t, ent.Date.Before(prev))
			prev = *ent.Date
		}
	})
}
