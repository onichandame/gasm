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
type Order string

const (
	ASC  Order = `ASC`
	DESC Order = `DESC`
)

type Sort map[string]Order

type Query struct {
	Filter     `json:"filter"`
	Pagination `json:"pagination"`
	Sort       `json:"sort"`
	Select     []string `json:"select"`
	Relations  []string `json:"relations"` // selection of relations not supported
}

type SingleQuery struct {
	ID interface{} `json:"id"`
}
