package gimquery

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/onichandame/gim"
	"github.com/onichandame/go-crud/core"
	goutils "github.com/onichandame/go-utils"
)

type OperationConfig struct {
	Disabled bool
	gim.Module
}
type OperationModuleConfig struct {
	Disabled bool
	gim.Module
	Many OperationConfig
	One  OperationConfig
}
type CreateModuleConfig struct {
	Input interface{}
	OperationConfig
}
type UpdateModuleConfig struct {
	Input interface{}
	OperationModuleConfig
}
type DeleteModuleConfig struct {
	OperationModuleConfig
}
type GimModuleConfig struct {
	gim.Module
	Assembler    core.Assembler
	QueryService core.QueryService
	Endpoint     string
	Entity       interface{}
	Output       interface{}
	Read         OperationModuleConfig
	Create       CreateModuleConfig
	Update       UpdateModuleConfig
	Delete       DeleteModuleConfig
}

func CreateGimModule(config GimModuleConfig) *gim.Module {
	if config.QueryService == nil {
		panic(fmt.Errorf("query service must be defined"))
	}
	if config.Create.Input == nil {
		config.Create.Input = config.Output
	}
	if config.Update.Input == nil {
		config.Update.Input = config.Output
	}
	basepath := strings.ToLower(goutils.UnwrapType(reflect.TypeOf(config.Output)).Name())
	if !reflect.ValueOf(config.Endpoint).IsZero() {
		basepath = config.Endpoint
	}
	mod := &config.Module
	initModule(mod)
	if config.Assembler == nil {
		var assm core.DefaultAssembler
		assm.Entity = config.Entity
		assm.DTO = config.Output
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
			mod.Controllers = append(mod.Controllers, &gim.Controller{
				Path: basepath,
				Routes: []*gim.Route{
					{
						Endpoint: ":id",
						Get: func(args gim.RouteArgs) interface{} {
							return config.Assembler.ConvertToDTO(config.QueryService.FindByID(args.Ctx.Param("id")))
						},
					},
				},
			})
		}
		if !config.Read.Many.Disabled {
			importModule(mod, &config.Read.Many.Module)
			mod := &config.Read.Many.Module
			initModule(mod)
			mod.Controllers = append(mod.Controllers, &gim.Controller{
				Path: basepath,
				Routes: []*gim.Route{
					{
						Get: func(args gim.RouteArgs) interface{} {
							input := parseQuery(args.Ctx)
							return config.Assembler.ConvertToDTOs(config.QueryService.Find(*input))
						},
					},
				},
			})
		}
	}
	// create
	if !config.Create.Disabled {
		importModule(mod, &config.Create.Module)
		mod := &config.Create.Module
		initModule(mod)
		mod.Controllers = append(mod.Controllers, &gim.Controller{
			Path: basepath,
			Routes: []*gim.Route{
				{
					Post: func(args gim.RouteArgs) interface{} {
						input := parseCreateInput(args.Ctx, config.Create.Input)
						return config.Assembler.ConvertToDTO(config.QueryService.Create(config.Assembler.ConvertToCreateEntity(input)))
					},
				},
			},
		})
	}
	// update
	if !config.Update.Disabled {
		importModule(mod, &config.Update.Module)
		mod := &config.Update.Module
		initModule(mod)
		if !config.Update.One.Disabled {
			importModule(mod, &config.Update.One.Module)
			mod := &config.Update.One.Module
			initModule(mod)
			mod.Controllers = append(mod.Controllers, &gim.Controller{
				Path: basepath,
				Routes: []*gim.Route{
					{
						Endpoint: ":id",
						Put: func(args gim.RouteArgs) interface{} {
							input := parseUpdateInput(args.Ctx, config.Update.Input)
							return config.Assembler.ConvertToDTO(config.QueryService.UpdateOne(args.Ctx.Param("id"), config.Assembler.ConvertToUpdateEntity(input)))
						},
					},
				},
			})
		}
		if !config.Update.Many.Disabled {
			importModule(mod, &config.Update.Many.Module)
			mod := &config.Update.Many.Module
			initModule(mod)
			mod.Controllers = append(mod.Controllers, &gim.Controller{
				Path: basepath,
				Routes: []*gim.Route{
					{
						Put: func(args gim.RouteArgs) interface{} {
							filter := parseFilter(args.Ctx)
							input := parseUpdateInput(args.Ctx, config.Update.Input)
							return config.Assembler.ConvertToDTO(config.QueryService.UpdateMany(*filter, config.Assembler.ConvertToUpdateEntity(input)))
						},
					},
				},
			})
		}
	}
	// delete
	if !config.Delete.Disabled {
		importModule(mod, &config.Delete.Module)
		mod := &config.Delete.Module
		initModule(mod)
		if !config.Delete.One.Disabled {
			importModule(mod, &config.Delete.One.Module)
			mod := &config.Delete.One.Module
			initModule(mod)
			mod.Controllers = append(mod.Controllers, &gim.Controller{
				Path: basepath,
				Routes: []*gim.Route{
					{
						Endpoint: ":id",
						Delete: func(args gim.RouteArgs) interface{} {
							return config.Assembler.ConvertToDTO(config.QueryService.DeleteOne(args.Ctx.Param("id")))
						},
					},
				},
			})
		}
		if !config.Delete.Many.Disabled {
			importModule(mod, &config.Delete.Many.Module)
			mod := &config.Delete.Many.Module
			initModule(mod)
			mod.Controllers = append(mod.Controllers, &gim.Controller{
				Path: basepath,
				Routes: []*gim.Route{
					{
						Endpoint: "",
						Delete: func(args gim.RouteArgs) interface{} {
							filter := parseFilter(args.Ctx)
							return config.Assembler.ConvertToDTO(config.QueryService.DeleteMany(*filter))
						},
					},
				},
			})
		}
	}
	return mod
}
