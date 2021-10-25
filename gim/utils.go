package gimquery

import (
	"reflect"

	"github.com/fatih/structtag"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/onichandame/gim"
	"github.com/onichandame/go-crud/core"
	goutils "github.com/onichandame/go-utils"
)

func parseFilter(c *gin.Context) *core.Filter {
	query := parseQuery(c)
	return query.Filter
}

func parseQuery(c *gin.Context) *core.Query {
	var queryInput core.Query
	goutils.Assert(c.ShouldBindQuery(&queryInput))
	goutils.Assert(validator.New().Struct(&queryInput))
	return &queryInput
}

func parseCreateInput(c *gin.Context, dto interface{}) interface{} {
	input := reflect.New(goutils.UnwrapType(reflect.TypeOf(dto))).Interface()
	goutils.Assert(c.ShouldBindJSON(&input))
	goutils.Assert(validator.New().Struct(input))
	return input
}

func parseUpdateInput(c *gin.Context, dto interface{}) interface{} {
	getTagFieldMap := func() (m map[string]string) {
		m = make(map[string]string)
		t := goutils.UnwrapType(reflect.TypeOf(dto))
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			tags, err := structtag.Parse(string(field.Tag))
			goutils.Assert(err)
			jsonTag, err := tags.Get("json")
			goutils.Assert(err)
			m[jsonTag.Name] = field.Name
		}
		return m
	}
	tagfieldmap := getTagFieldMap()
	update := make(map[string]interface{})
	goutils.Assert(c.ShouldBindJSON(&update))
	for field := range update {
		if _, ok := tagfieldmap[field]; !ok {
			delete(update, field)
		}
	}
	return update
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
