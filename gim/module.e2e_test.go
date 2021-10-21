package gimquery_test

import (
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

type Entity struct {
	gorm.Model
	Name string
}
type EntityDTO struct {
	Name string
}

func TestModule(t *testing.T) {
	createDB := func() *gorm.DB {
		if db, err := gorm.Open(sqlite.Open(":memory:")); err != nil {
			panic(err)
		} else {
			assert.Nil(t, db.AutoMigrate(&Entity{}))
			return db
		}
	}
	createApp := func(db *gorm.DB) *gim.Module {
		return gimquery.CreateGimModule(gimquery.GimModuleConfig{
			Endpoint: "entity",
			Entity:   &Entity{},
			Output:   &EntityDTO{},
			DB:       db,
		})
	}
	t.Run("read one", func(t *testing.T) {
		db := createDB()
		app := createApp(db)
		eng := app.Bootstrap()
		if err := db.Create(&Entity{}).Error; err != nil {
			panic(err)
		}
		w := httptest.NewRecorder()
		req := newRequest("/entity/get", `{"id":1}`)
		eng.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
	})
}

func newRequest(url string, body string) *http.Request {
	reader := strings.NewReader(body)
	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}
