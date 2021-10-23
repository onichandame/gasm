package gormquery_test

import (
	"testing"

	gormquery "github.com/onichandame/go-crud/gorm"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSelect(t *testing.T) {
	type Entity struct {
		gorm.Model
		Name *string
		Int  *int
	}
	db, err := gorm.Open(sqlite.Open(":memory:"))
	assert.Nil(t, err)
	assert.Nil(t, db.AutoMigrate(&Entity{}))
	name := "asdf"
	i := 1
	assert.Nil(t, db.Create(&Entity{Name: &name, Int: &i}).Error)
	t.Run("no select(all)", func(t *testing.T) {
		var ent Entity
		assert.Nil(t, db.First(&ent).Error)
		assert.NotNil(t, ent.Name)
		assert.NotNil(t, ent.Int)
	})
	t.Run("select fields", func(t *testing.T) {
		var ent Entity
		assert.Nil(t, db.Scopes(gormquery.Select([]string{"name"})).First(&ent).Error)
		assert.NotNil(t, ent.Name)
		assert.Nil(t, ent.Int)
	})
}
