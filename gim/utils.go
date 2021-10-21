package gimquery

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/onichandame/gim"
	"github.com/onichandame/go-crud/core"
	"gorm.io/gorm"
)

func parseQuery(c *gin.Context) *core.Query {
	var queryInput core.Query
	if err := c.ShouldBindJSON(&queryInput); err != nil {
		panic(err)
	}
	return &queryInput
}

func parseSingleQuery(c *gin.Context) *core.SingleQuery {
	var input core.SingleQuery
	if err := c.ShouldBindJSON(&input); err != nil {
		panic(err)
	}
	return &input
}

func getDB(app context.Context, raw interface{}) *gorm.DB {
	if db, ok := raw.(*gorm.DB); ok {
		return db
	} else {
		return app.Value(raw).(*gorm.DB)
	}
}

func initModule(mod *gim.Module) {
	if mod.Imports == nil {
		mod.Imports = make([]*gim.Module, 0)
	}
	if mod.Controllers == nil {
		mod.Controllers = make([]*gim.Controller, 0)
	}
	if mod.Middlewares == nil {
		mod.Middlewares = make([]*gim.Middleware, 0)
	}
}

func importModule(parent *gim.Module, child *gim.Module) {
	parent.Imports = append(parent.Imports, child)
}
