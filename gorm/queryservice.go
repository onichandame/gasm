package gormquery

import (
	"reflect"

	"github.com/jinzhu/copier"
	"github.com/onichandame/go-crud/core"
	goutils "github.com/onichandame/go-utils"
	"gorm.io/gorm"
)

type GORMQueryService struct {
	DB     *gorm.DB
	Entity interface{}
}

func (s *GORMQueryService) newEntity() interface{} {
	return reflect.New(goutils.UnwrapType(reflect.TypeOf(s.Entity))).Interface()
}
func (s *GORMQueryService) newEntitySlice() interface{} {
	return reflect.MakeSlice(goutils.UnwrapType(reflect.TypeOf(s.Entity)), 0, 0).Interface()
}

func (s *GORMQueryService) FindByID(id interface{}) interface{} {
	data := s.newEntity()
	if err := s.DB.First(data, id).Error; err != nil {
		panic(err)
	}
	return data
}
func (s *GORMQueryService) Find(q core.Query) interface{} {
	data := s.newEntitySlice()
	if err := s.DB.Scopes(Pagination(&q.Pagination), Filter(&q.Filter)).Find(data).Error; err != nil {
		panic(err)
	}
	return data
}
func (s *GORMQueryService) CreateOne(input interface{}) interface{} {
	data := s.newEntity()
	if err := copier.Copy(data, input); err != nil {
		panic(err)
	}
	if err := s.DB.Create(data).Error; err != nil {
		panic(err)
	}
	return input
}
func (s *GORMQueryService) CreateMany(input []interface{}) interface{} {
	data := s.newEntitySlice()
	if err := copier.Copy(data, input); err != nil {
		panic(err)
	}
	if err := s.DB.Create(data).Error; err != nil {
		panic(err)
	}
	return data
}
func (s *GORMQueryService) UpdateOne(id interface{}, update interface{}) interface{} {
	data := s.newEntity()
	if err := s.DB.First(data, id).Error; err != nil {
		panic(err)
	}
	if err := copier.Copy(data, update); err != nil {
		panic(err)
	}
	return data
}
func (s *GORMQueryService) UpdateMany(filter core.Filter, update interface{}) int {
	res := s.DB.Model(s.newEntity()).Scopes(Filter(&filter)).Updates(update)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}
func (s *GORMQueryService) DeleteOne(id interface{}) interface{} {
	data := s.newEntity()
	if err := s.DB.First(data, id).Error; err != nil {
		panic(err)
	}
	if err := s.DB.Delete(data).Error; err != nil {
		panic(err)
	}
	return data
}
func (s *GORMQueryService) DeleteMany(filter core.Filter) int {
	res := s.DB.Scopes(Filter(&filter)).Delete(s.newEntity())
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

func CreateGORMQueryService(db *gorm.DB, entity interface{}) *GORMQueryService {
	var svc GORMQueryService
	svc.DB = db
	svc.Entity = entity
	return &svc
}
