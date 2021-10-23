package core

import (
	"fmt"
	"reflect"

	"github.com/fatih/structtag"
	"github.com/jinzhu/copier"
	goutils "github.com/onichandame/go-utils"
)

// convert between dto and entity
type Assembler interface {
	// dto query to entity query
	ConvertToQuery(Query) Query
	// entity to dto for response
	ConvertToDTO(interface{}) interface{}
	// entities to dtos for response
	ConvertToDTOs(interface{}) interface{}
	// create input dto to create entity for consumption by orm
	ConvertToCreateEntity(interface{}) interface{}
	// create input dtos to create entities for consumption by orm
	ConvertToCreateEntities(interface{}) interface{}
	// update input dto to update entity form consumption by orm
	ConvertToUpdateEntity(interface{}) interface{}
}

type DefaultAssembler struct {
	DTO    interface{}
	Entity interface{}
}

func (a *DefaultAssembler) getDTOTagFieldMap() (m map[string]string) {
	m = make(map[string]string)
	t := goutils.UnwrapType(reflect.TypeOf(a.DTO))
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tags, err := structtag.Parse(string(field.Tag))
		if err != nil {
			panic(err)
		}
		jsonTag, err := tags.Get("json")
		if err != nil {
			panic(err)
		}
		m[jsonTag.Name] = field.Name
	}
	return m
}

func (a *DefaultAssembler) ConvertToQuery(q Query) Query {
	tagFieldMap := a.getDTOTagFieldMap()
	var parseFilter func(in *Filter)
	parseFilter = func(in *Filter) {
		for _, and := range in.And {
			parseFilter(and)
		}
		for _, or := range in.Or {
			parseFilter(or)
		}
		for _, not := range in.Not {
			parseFilter(not)
		}
		for tagname, _ := range in.Fields {
			if fieldname, ok := tagFieldMap[tagname]; ok {
				in.Fields[fieldname] = in.Fields[tagname]
				delete(in.Fields, tagname)
			}
		}
	}
	parseFilter(q.Filter)
	parseSort := func() {}
	parseSort()
	parseSelect := func() {}
	parseSelect()
	parseRelations := func() {
		for i, tagname := range q.Relations {
			if fieldname, ok := tagFieldMap[tagname]; ok {
				q.Relations[i] = fieldname
			}
		}
	}
	parseRelations()
	return q
}

func (a *DefaultAssembler) ConvertToDTO(entity interface{}) interface{} {
	dto := reflect.New(goutils.UnwrapType(reflect.TypeOf(a.DTO))).Interface()
	if err := copier.Copy(dto, entity); err != nil {
		panic(err)
	}
	return dto
}
func (a *DefaultAssembler) ConvertToDTOs(entities interface{}) interface{} {
	dtos := reflect.MakeSlice(reflect.SliceOf(goutils.UnwrapType(reflect.TypeOf(a.DTO))), 0, 0).Interface()
	if err := copier.Copy(&dtos, entities); err != nil {
		panic(err)
	}
	return dtos
}
func (a *DefaultAssembler) ConvertToCreateEntity(dto interface{}) interface{} {
	ent := reflect.New(goutils.UnwrapType(reflect.TypeOf(a.Entity))).Interface()
	if err := copier.Copy(ent, dto); err != nil {
		panic(err)
	}
	return ent
}
func (a *DefaultAssembler) ConvertToCreateEntities(dto interface{}) interface{} {
	ents := reflect.MakeSlice(reflect.SliceOf(goutils.UnwrapType(reflect.TypeOf(a.Entity))), 0, 0).Interface()
	if err := copier.Copy(&ents, dto); err != nil {
		panic(err)
	}
	return ents
}
func (a *DefaultAssembler) ConvertToUpdateEntity(dto interface{}) interface{} {
	return dto.(map[string]interface{})
}

func GetAssemblerToken(dtoname string) string {
	return fmt.Sprintf("__%s_assembler__", dtoname)
}
