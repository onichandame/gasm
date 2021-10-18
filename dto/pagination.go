package dto

type PaginationInputDTO struct {
	Page    uint `json:"page" validate:"required,min=1"`
	PerPage uint `json:"per_page" validate:"required,min=1"`
}
