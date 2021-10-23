package gormquery_test

import (
	"testing"

	gormquery "github.com/onichandame/go-crud/gorm"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Parent struct {
	gorm.Model
	Children []*Child
}
type Child struct {
	gorm.Model
	Parent   *Parent
	ParentID uint
}

func TestRelation(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"))
	assert.Nil(t, err)
	assert.Nil(t, db.AutoMigrate(&Parent{}, &Child{}))
	var childless, parent Parent
	var orphan, child Child
	assert.Nil(t, db.Create(&childless).Error)
	assert.Nil(t, db.Create(&orphan).Error)
	parent.Children = []*Child{{}}
	assert.Nil(t, db.Create(&parent).Error)
	assert.Len(t, parent.Children, 1)
	child = *parent.Children[0]
	t.Run("parent", func(t *testing.T) {
		var res Parent
		assert.Nil(t, db.Scopes(gormquery.Relation([]string{"Children"})).First(&res, parent.ID).Error)
		assert.Len(t, res.Children, 1)
		assert.Equal(t, child.ID, res.Children[0].ID)
	})
	t.Run("child", func(t *testing.T) {
		var res Child
		assert.Nil(t, db.Scopes(gormquery.Relation([]string{"Parent"})).First(&res, child.ID).Error)
		assert.NotNil(t, res.Parent)
	})
	t.Run("orphan", func(t *testing.T) {
		var res Child
		assert.Nil(t, db.Scopes(gormquery.Relation([]string{"Parent"})).First(&res, orphan.ID).Error)
		assert.Nil(t, res.Parent)
	})
}
