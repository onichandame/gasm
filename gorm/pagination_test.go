package gormquery_test

import (
	"testing"

	"github.com/onichandame/go-crud/core"
	gormquery "github.com/onichandame/go-crud/gorm"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestPagination(t *testing.T) {
	type Entity struct {
		gorm.Model
	}
	db, err := gorm.Open(sqlite.Open(":memory:"))
	assert.Nil(t, err)
	assert.Nil(t, db.AutoMigrate(&Entity{}))
	totalCount := 100
	for i := 0; i < totalCount; i++ {
		assert.Nil(t, db.Create(&Entity{}).Error)
	}
	t.Run("unpaged", func(t *testing.T) {
		var ents []Entity
		assert.Nil(t, db.Scopes(gormquery.Pagination(nil)).Find(&ents).Error)
		assert.Len(t, ents, totalCount)
	})
	t.Run("first page", func(t *testing.T) {
		var ents []Entity
		assert.Nil(t, db.Scopes(gormquery.Pagination(&core.Pagination{Page: 1, PerPage: 10})).Find(&ents).Error)
		assert.Len(t, ents, 10)
		for _, ent := range ents {
			assert.GreaterOrEqual(t, int(ent.ID), 0)
			assert.LessOrEqual(t, int(ent.ID), 10)
		}
	})
	t.Run("fifth page", func(t *testing.T) {
		var ents []Entity
		assert.Nil(t, db.Scopes(gormquery.Pagination(&core.Pagination{Page: 5, PerPage: 10})).Find(&ents).Error)
		assert.Len(t, ents, 10)
		for _, ent := range ents {
			assert.GreaterOrEqual(t, int(ent.ID), 40)
			assert.LessOrEqual(t, int(ent.ID), 50)
		}
	})
}
