package core

import (
	"fmt"
	"reflect"

	"github.com/jinzhu/copier"
	goutils "github.com/onichandame/go-utils"
)

// convert between dto and entity
type Assembler interface {
	// entity to dto for response
	ConvertToDTO(interface{}) interface{}
	// create input dto to create entity for consumption by orm
	ConvertToCreateEntity(interface{}) interface{}
	// update input dto to update entity form consumption by orm
	ConvertToUpdateEntity(interface{}) interface{}
}

type DefaultAssembler struct {
	DTO    interface{}
	Entity interface{}
}

func (a *DefaultAssembler) ConvertToQuery(q *Query) *Query { return q }
func (a *DefaultAssembler) ConvertToDTO(entity interface{}) interface{} {
	dto := reflect.New(goutils.UnwrapType(reflect.TypeOf(a.Entity))).Interface()
	if err := copier.Copy(dto, entity); err != nil {
		panic(err)
	}
	return dto
}
func (a *DefaultAssembler) ConvertToCreateEntity(dto interface{}) interface{} {
	ent := reflect.New(goutils.UnwrapType(reflect.TypeOf(a.Entity))).Interface()
	if err := copier.Copy(ent, dto); err != nil {
		panic(err)
	}
	return ent
}
func (a *DefaultAssembler) ConvertToUpdateEntity(dto interface{}) interface{} {
	ent := reflect.New(goutils.UnwrapType(reflect.TypeOf(a.Entity))).Interface()
	if err := copier.Copy(ent, dto); err != nil {
		panic(err)
	}
	return ent
}

func GetAssemblerToken(dtoname string) string {
	return fmt.Sprintf("__%s_assembler__", dtoname)
}
