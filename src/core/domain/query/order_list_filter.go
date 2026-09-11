package query

type OrderListFilter struct {
	Page  int
	Limit int
}

func (this *OrderListFilter) Normalize() {
	if this.Page <= 0 {
		this.Page = 1
	}
	if this.Limit <= 0 {
		this.Limit = DefaultPageLimit
	}
}

func (this OrderListFilter) Offset() int {
	if this.Page <= 0 {
		return 0
	}
	return (this.Page - 1) * this.Limit
}
