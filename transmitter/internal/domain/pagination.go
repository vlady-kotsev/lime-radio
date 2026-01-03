package domain

type PaginationParams struct {
	Page     int `query:"page" validate:"min=1" json:"page"`
	PageSize int `query:"page_size" validate:"min=1,max=30" json:"page_size"`
}

func (p *PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *PaginationParams) GetLimit() int {
	return p.PageSize
}

func NewPaginationParams(page, pageSize int) *PaginationParams {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 5
	}
	if pageSize > 30 {
		pageSize = 30
	}

	return &PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}
}

type PaginatedResult[T any] struct {
	Data    []T
	Page    int
	HasNext bool
}

func NewPaginatedResult[T any](data []T, total int, params *PaginationParams) *PaginatedResult[T] {
	totalPages := (total + params.PageSize - 1) / params.PageSize

	return &PaginatedResult[T]{
		Data:    data,
		Page:    params.Page,
		HasNext: params.Page < totalPages,
	}
}
