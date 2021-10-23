package gimquery_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/onichandame/gim"
	gimquery "github.com/onichandame/go-crud/gim"
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
			return gimquery.CreateGimModule(gimquery.GimModuleConfig{
				Endpoint: "entity",
				Entity:   &Entity{},
				Output:   &EntityDTO{},
				DB:       db,
			}), db
		}
	}
	newRequest := func(url string, body string) *http.Request {
		reader := strings.NewReader(body)
		req, err := http.NewRequest("POST", url, reader)
		if err != nil {
			panic(err)
		}
		req.Header.Set("Content-Type", "application/json")
		return req
	}
	t.Run("read one", func(t *testing.T) {
		app, db := createApp()
		eng := app.Bootstrap()
		if err := db.Create(&Entity{}).Error; err != nil {
			panic(err)
		}
		w := httptest.NewRecorder()
		req := newRequest("/entity/get", `{"id":1}`)
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
			req := newRequest("/entity/list", `{}`)
			eng.ServeHTTP(w, req)
			var res []EntityDTO
			assert.Equal(t, 200, w.Code)
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &res))
			assert.Len(t, res, 2)
		})
		t.Run("filter", func(t *testing.T) {
			w := httptest.NewRecorder()
			req := newRequest("/entity/list", `{"filter":{"fields":{"name":"asdf"}}}`)
			eng.ServeHTTP(w, req)
			var res []EntityDTO
			assert.Equal(t, 200, w.Code)
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &res))
			assert.Len(t, res, 1)
		})
		t.Run("pagination", func(t *testing.T) {
			w := httptest.NewRecorder()
			req := newRequest("/entity/list", `{"pagination":{"per_page":1,"page":1}}`)
			eng.ServeHTTP(w, req)
			var res []EntityDTO
			assert.Equal(t, 200, w.Code)
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &res))
			assert.Len(t, res, 1)
			first := res[0]
			w = httptest.NewRecorder()
			req = newRequest("/entity/list", `{"pagination":{"per_page":1,"page":2}}`)
			eng.ServeHTTP(w, req)
			assert.Equal(t, 200, w.Code)
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &res))
			assert.Len(t, res, 1)
			second := res[0]
			assert.NotEqual(t, first.ID, second.ID)
		})
		t.Run("sorting", func(t *testing.T) {
			w := httptest.NewRecorder()
			req := newRequest("/entity/list", `{"sort":{"int":"asc"}}`)
			eng.ServeHTTP(w, req)
			assert.Equal(t, 200, w.Code)
			var res []EntityDTO
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &res))
			assert.Len(t, res, 2)
			assert.Less(t, res[0].Int, res[1].Int)
			w = httptest.NewRecorder()
			req = newRequest("/entity/list", `{"sort":{"int":"desc"}}`)
			eng.ServeHTTP(w, req)
			assert.Equal(t, 200, w.Code)
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &res))
			assert.Len(t, res, 2)
			assert.Greater(t, res[0].Int, res[1].Int)
		})
	})
	t.Run("create one", func(t *testing.T) {
		app, db := createApp()
		w := httptest.NewRecorder()
		req := newRequest("/entity/createOne", `{"name":"asdf","int":1}`)
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
	t.Run("create many", func(t *testing.T) {
		app, db := createApp()
		w := httptest.NewRecorder()
		req := newRequest("/entity/createMany", `[{"name":"asdf","int":1},{"name":"zxcv","int":2}]`)
		eng := app.Bootstrap()
		eng.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		var res []EntityDTO
		assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &res))
		assert.Len(t, res, 2)
		for _, v := range res {
			assert.Nil(t, db.First(&Entity{}, v.ID).Error)
		}
	})
}
