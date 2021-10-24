package core

import "fmt"

type QueryService interface {
	Find(Query) interface{}
	FindByID(id interface{}) interface{}
	Create(interface{}) interface{}
	UpdateOne(interface{}, interface{}) interface{}
	UpdateMany(Filter, interface{}) int
	DeleteOne(interface{}) interface{}
	DeleteMany(Filter) int
}

func GetQueryServiceToken(entityname string) string {
	return fmt.Sprintf("__%s_queryservice__", entityname)
}
