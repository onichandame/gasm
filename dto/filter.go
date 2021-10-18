package dto

type FieldFilterOperator string

const (
	Is  FieldFilterOperator = "is"
	Eq  FieldFilterOperator = "eq"
	GT  FieldFilterOperator = "gt"
	LT  FieldFilterOperator = "lt"
	GTE FieldFilterOperator = "gte"
	LTE FieldFilterOperator = "lte"
	In  FieldFilterOperator = "in"
)

type FieldFilter map[FieldFilterOperator]interface{}

type Filter struct {
	Fields map[string]interface{} `json:"fields"`
	And    []*Filter              `json:"and"`
	Or     []*Filter              `json:"or"`
	Not    []*Filter              `json:"not"`
}
