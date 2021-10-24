package gimquery_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/onichandame/gim"
	gimquery "github.com/onichandame/go-crud/gim"
	gormquery "github.com/onichandame/go-crud/gorm"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestModule(t *testing.T) {
	type Entity struct {
		gorm.Model
		Int  int
		Name string
	}
	type EntityDTO struct {
		ID   uint   `json:"id"`
		Int  int    `json:"int"`
		Name string `json:"name"`
	}
	createApp := func() (*gim.Module, *gorm.DB) {
		if db, err := gorm.Open(sqlite.Open(":memory:")); err != nil {
			panic(err)
		} else {
			assert.Nil(t, db.AutoMigrate(&Entity{}))
			queryService := gormquery.CreateGORMQueryService(db, &Entity{})
			return gimquery.CreateGimModule(gimquery.GimModuleConfig{
				Endpoint:     "entities",
				QueryService: queryService,
				Entity:       &Entity{},
				Output:       &EntityDTO{},
			}), db
		}
	}
	newRequest := func(method string, url string, body string) *http.Request {
		reader := strings.NewReader(body)
		req, err := http.NewRequest(method, url, reader)
		if err != nil {
			panic(err)
		}
		if method != "GET" {
			req.Header.Set("Content-Type", "application/json")
		}
		return req
	}
	t.Run("read one", func(t *testing.T) {
		app, db := createApp()
		eng := app.Bootstrap()
		if err := db.Create(&Entity{}).Error; err != nil {
			panic(err)
		}
		w := httptest.NewRecorder()
		req := newRequest("GET", "/entities/1", ``)
		eng.ServeHTTP(w, req)
		var res EntityDTO
		assert.Equal(t, 200, w.Code)
		assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &res))
		assert.NotNil(t, res)
	})
	t.Run("read many", func(t *testing.T) {
		app, db := createApp()
		if err := db.Create(&[]Entity{
			{Int: 1},
			{
				Int:  2,
				Name: "asdf",
			},
		}).Error; err != nil {
			panic(err)
		}
		eng := app.Bootstrap()
		t.Run("all", func(t *testing.T) {
			w := httptest.NewRecorder()
			req := newRequest("GET", "/entities", ``)
			eng.ServeHTTP(w, req)
			var res []EntityDTO
			assert.Equal(t, 200, w.Code)
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &res))
			assert.Len(t, res, 2)
		})
		t.Run("filter", func(t *testing.T) {
			w := httptest.NewRecorder()
			req := newRequest("GET", `/entities?filter={"fields":{"name":"asdf"}}`, ``)
			eng.ServeHTTP(w, req)
			var res []EntityDTO
			assert.Equal(t, 200, w.Code)
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &res))
			assert.Len(t, res, 1)
		})
		t.Run("pagination", func(t *testing.T) {
			w := httptest.NewRecorder()
			req := newRequest("GET", `/entities?pagination={"perPage":1,"page":1}`, ``)
			eng.ServeHTTP(w, req)
			var res []EntityDTO
			assert.Equal(t, 200, w.Code)
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &res))
			assert.Len(t, res, 1)
			first := res[0]
			w = httptest.NewRecorder()
			req = newRequest("GET", `/entities?pagination={"perPage":1,"page":2}`, ``)
			eng.ServeHTTP(w, req)
			assert.Equal(t, 200, w.Code)
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &res))
			assert.Len(t, res, 1)
			second := res[0]
			assert.NotEqual(t, first.ID, second.ID)
		})
		t.Run("sorting", func(t *testing.T) {
			w := httptest.NewRecorder()
			req := newRequest("GET", `/entities?sort={"int":"asc"}`, ``)
			eng.ServeHTTP(w, req)
			assert.Equal(t, 200, w.Code)
			var res []EntityDTO
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &res))
			assert.Len(t, res, 2)
			assert.Less(t, res[0].Int, res[1].Int)
			w = httptest.NewRecorder()
			req = newRequest("GET", `/entities?sort={"int":"desc"}`, ``)
			eng.ServeHTTP(w, req)
			assert.Equal(t, 200, w.Code)
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &res))
			assert.Len(t, res, 2)
			assert.Greater(t, res[0].Int, res[1].Int)
		})
	})
	t.Run("create", func(t *testing.T) {
		app, db := createApp()
		w := httptest.NewRecorder()
		req := newRequest("POST", "/entities", `{"name":"asdf","int":1}`)
		eng := app.Bootstrap()
		eng.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		var res EntityDTO
		assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &res))
		var ent Entity
		assert.Nil(t, db.First(&ent, res.ID).Error)
		assert.Equal(t, "asdf", ent.Name)
		assert.Equal(t, 1, ent.Int)
	})
	t.Run("update one", func(t *testing.T) {
		app, db := createApp()
		assert.Nil(t, db.Create(&Entity{Name: "asdf"}).Error)
		w := httptest.NewRecorder()
		req := newRequest(http.MethodPut, "/entities/1", `{"name":"zxcv"}`)
		eng := app.Bootstrap()
		eng.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		var ent Entity
		assert.Nil(t, db.First(&ent, 1).Error)
		assert.Equal(t, "zxcv", ent.Name)
	})
	t.Run("update many", func(t *testing.T) {
		app, db := createApp()
		assert.Nil(t, db.Create([]Entity{{Name: "asdf"}, {Name: "asdf"}, {Name: "qwer"}}).Error)
		w := httptest.NewRecorder()
		req := newRequest(http.MethodPut, `/entities?filter={"fields":{"name":"asdf"}}`, `{"name":"zxcv"}`)
		eng := app.Bootstrap()
		eng.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		var ents []Entity
		assert.Nil(t, db.Where("name = ?", "zxcv").Find(&ents).Error)
		assert.Equal(t, 2, len(ents))
	})
	t.Run("delete one", func(t *testing.T) {
		app, db := createApp()
		assert.Nil(t, db.Create(&Entity{Name: "asdf"}).Error)
		w := httptest.NewRecorder()
		req := newRequest(http.MethodDelete, `/entities/1`, ``)
		eng := app.Bootstrap()
		eng.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		var count int64
		assert.Nil(t, db.Model(&Entity{}).Count(&count).Error)
		assert.Equal(t, 0, int(count))
	})
	t.Run("delete many", func(t *testing.T) {
		app, db := createApp()
		assert.Nil(t, db.Create([]Entity{{Name: "asdf"}, {Name: "asdf"}, {}}).Error)
		w := httptest.NewRecorder()
		req := newRequest(http.MethodDelete, `/entities?filter={"fields":{"name":"asdf"}}`, ``)
		eng := app.Bootstrap()
		eng.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.NotNil(t, db.Where("name = ?", "asdf").First(&Entity{}).Error)
		var count int64
		assert.Nil(t, db.Model(&Entity{}).Count(&count).Error)
		assert.Equal(t, 1, int(count))
	})
}
