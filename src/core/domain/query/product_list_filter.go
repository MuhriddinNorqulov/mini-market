package query

type ProductListFilter struct {
	Page  int
	Limit int
}

func (this *ProductListFilter) Normalize() {
	if this.Page <= 0 {
		this.Page = 1
	}
	if this.Limit <= 0 {
		this.Limit = DefaultPageLimit
	}
}

func (this ProductListFilter) Offset() int {
	if this.Page <= 0 {
		return 0
	}
	return (this.Page - 1) * this.Limit
}
