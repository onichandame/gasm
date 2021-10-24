package core

type Filter struct {
	Fields map[string]interface{} `form:"fields"`
	And    []*Filter              `form:"and"`
	Or     []*Filter              `form:"or"`
	Not    []*Filter              `form:"not"`
}
type Pagination struct {
	Page    uint `form:"page" validate:"required,min=1"`
	PerPage uint `form:"perPage" validate:"required,min=1"`
}

const (
	ASC  = `ASC`
	DESC = `DESC`
)

type Sort map[string]string

type Query struct {
	Filter     *Filter     `form:"filter"`
	Pagination *Pagination `form:"pagination"`
	Sort       Sort        `form:"sort"`
	Select     []string    `form:"select"`
	Relations  []string    `form:"relations"`
}
