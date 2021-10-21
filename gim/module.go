package gimquery

import (
	"reflect"
	"strings"

	"github.com/onichandame/gim"
	"github.com/onichandame/go-crud/core"
	gormquery "github.com/onichandame/go-crud/gorm"
	goutils "github.com/onichandame/go-utils"
	"gorm.io/gorm"
)

type OperationConfig struct {
	Disabled bool
	Endpoint string
	gim.Module
}
type OperationModuleConfig struct {
	Disabled bool
	gim.Module
	Many OperationConfig
	One  OperationConfig
}
type QueryModuleConfig struct {
	OperationModuleConfig
}
type MutationModuleConfig struct {
	Input interface{}
	OperationModuleConfig
}
type GimModuleConfig struct {
	gim.Module
	Assembler core.Assembler
	DB        *gorm.DB
	Endpoint  string
	Entity    interface{}
	Output    interface{}
	Read      QueryModuleConfig
	Create    MutationModuleConfig
	Update    MutationModuleConfig
	Delete    MutationModuleConfig
}

func CreateGimModule(config GimModuleConfig) *gim.Module {
	basepath := strings.ToLower(goutils.UnwrapType(reflect.TypeOf(config.Output)).Name())
	if !reflect.ValueOf(config.Endpoint).IsZero() {
		basepath = config.Endpoint
	}
	mod := &config.Module
	initModule(mod)
	queryService := gormquery.CreateGORMQueryService(config.DB, config.Entity)
	if config.Assembler == nil {
		var assm core.DefaultAssembler
		assm.DTO = config.Output
		assm.Entity = config.Entity
		config.Assembler = &assm
	}
	// read
	if !config.Read.Disabled {
		importModule(mod, &config.Read.Module)
		mod := &config.Read.Module
		initModule(mod)
		if !config.Read.One.Disabled {
			importModule(mod, &config.Read.One.Module)
			mod := &config.Read.One.Module
			initModule(mod)
			endpoint := config.Read.One.Endpoint
			if endpoint == "" {
				endpoint = "get"
			}
			mod.Controllers = append(mod.Controllers, &gim.Controller{
				Path: basepath,
				Routes: []*gim.Route{
					{
						Endpoint: endpoint,
						Post: func(args gim.RouteArgs) interface{} {
							input := parseSingleQuery(args.Ctx)
							return config.Assembler.ConvertToDTO(queryService.FindByID(input.ID))
						},
					},
				},
			})
		}
	}
	return mod
}
