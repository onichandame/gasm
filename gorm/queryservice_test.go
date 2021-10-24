package gormquery_test

import (
	"testing"

	"github.com/onichandame/go-crud/core"
	gormquery "github.com/onichandame/go-crud/gorm"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGORMQueryService(t *testing.T) {
	type Entity struct {
		gorm.Model
		Name string
	}
	createSvc := func() (*gormquery.GORMQueryService, *gorm.DB) {
		db, err := gorm.Open(sqlite.Open(":memory:"))
		assert.Nil(t, err)
		assert.Nil(t, db.AutoMigrate(&Entity{}))
		return gormquery.CreateGORMQueryService(db, &Entity{}), db
	}
	t.Run("create one", func(t *testing.T) {
		svc, db := createSvc()
		ent := svc.Create(&Entity{})
		assert.Nil(t, db.First(&Entity{}, ent.(*Entity).ID).Error)
	})
	t.Run("read one", func(t *testing.T) {
		svc, db := createSvc()
		var ent Entity
		assert.Nil(t, db.Create(&ent).Error)
		assert.Equal(t, ent.ID, svc.FindByID(ent.ID).(*Entity).ID)
	})
	t.Run("read many", func(t *testing.T) {
		svc, db := createSvc()
		assert.Nil(t, db.Create([]Entity{{}, {}}).Error)
		ents := svc.Find(core.Query{})
		assert.Len(t, ents, 2)
	})
	t.Run("update one", func(t *testing.T) {
		svc, db := createSvc()
		assert.Nil(t, db.Create(&Entity{Name: "asdf"}).Error)
		res := svc.UpdateOne(1, map[string]interface{}{"name": "zxcv"}).(*Entity)
		assert.Equal(t, "zxcv", res.Name)
		var ent Entity
		assert.Nil(t, db.First(&ent).Error)
		assert.Equal(t, res.Name, ent.Name)
	})
	t.Run("update many", func(t *testing.T) {
		svc, db := createSvc()
		assert.Nil(t, db.Create([]Entity{{Name: "asdf"}, {Name: "qwer"}}).Error)
		assert.Equal(t, 2, svc.UpdateMany(core.Filter{Fields: map[string]interface{}{"id": map[string]interface{}{"gt": 0}}}, map[string]interface{}{"name": "zxcv"}))
		var count int64
		assert.Nil(t, db.Model(&Entity{}).Where("name = ?", "zxcv").Count(&count).Error)
		assert.Equal(t, 2, int(count))
	})
	t.Run("delete one", func(t *testing.T) {
		svc, db := createSvc()
		assert.Nil(t, db.Create(&Entity{}).Error)
		assert.NotNil(t, svc.DeleteOne(1))
		var count int64
		assert.Nil(t, db.Model(&Entity{}).Count(&count).Error)
		assert.Equal(t, 0, int(count))
	})
	t.Run("delete many", func(t *testing.T) {
		svc, db := createSvc()
		assert.Nil(t, db.Create([]Entity{{}, {}}).Error)
		assert.Equal(t, 2, svc.DeleteMany(core.Filter{Fields: map[string]interface{}{"id": map[string]interface{}{"gt": 0}}}))
		var count int64
		assert.Nil(t, db.Model(&Entity{}).Count(&count).Error)
		assert.Equal(t, 0, int(count))
	})
}
