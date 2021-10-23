package core

type Filter struct {
	Fields map[string]interface{} `json:"fields"`
	And    []*Filter              `json:"and"`
	Or     []*Filter              `json:"or"`
	Not    []*Filter              `json:"not"`
}
type Pagination struct {
	Page    uint `json:"page" validate:"required,min=1"`
	PerPage uint `json:"per_page" validate:"required,min=1"`
}

const (
	ASC  = `ASC`
	DESC = `DESC`
)

type Sort map[string]string

type Query struct {
	Filter     *Filter     `json:"filter"`
	Pagination *Pagination `json:"pagination"`
	Sort       Sort        `json:"sort"`
	Select     []string    `json:"select"`
	Relations  []string    `json:"relations"`
}

type SingleQuery struct {
	ID interface{} `json:"id"`
}
