package entity

type PagingEntity[T any] struct {
	Page  int   `json:"page"`
	Size  int   `json:"size"`
	Total int64 `json:"total"`
	Items []T   `json:"items"`
}

func NewPagingEntity[T any](page, size int, total int64, items []T) *PagingEntity[T] {
	return &PagingEntity[T]{Page: page, Size: size, Total: total, Items: items}
}
