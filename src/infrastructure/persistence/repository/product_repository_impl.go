package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/ports/repository"
	"mini-market/src/core/domain/ports/unitofwork"
	"mini-market/src/core/domain/query"
	"mini-market/src/infrastructure/_errors"
	"mini-market/src/infrastructure/persistence/mapper"

	"gorm.io/gorm"
)

type ProductRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewProductRepositoryImpl(baseRepository *BaseRepository) repository.ProductRepository {
	return &ProductRepositoryImpl{BaseRepository: baseRepository}
}

func (this *ProductRepositoryImpl) Create(ctx context.Context, product *entity.ProductEntity) (*entity.ProductEntity, error) {
	var row mapper.ProductRow
	err := this.db(ctx).Raw(
		`INSERT INTO products (name, price, stock_quantity, created_at, updated_at)
		 VALUES (?, ?, ?, now(), now())
		 RETURNING id, name, price, stock_quantity`,
		product.Name, product.Price, product.StockQuantity,
	).Scan(&row).Error
	if err != nil {
		return nil, _errors.RawSQLErrorWrap(err)
	}
	return mapper.ProductRowToEntity(&row), nil
}

func (this *ProductRepositoryImpl) GetByID(ctx context.Context, id uint) (*entity.ProductEntity, error) {
	var row mapper.ProductRow
	err := this.db(ctx).Raw(
		`SELECT id, name, price, stock_quantity FROM products WHERE id = ? AND deleted_at IS NULL`,
		id,
	).Scan(&row).Error
	if err != nil {
		return nil, _errors.RawSQLErrorWrap(err)
	}
	if row.ID == 0 {
		return nil, _errors.RawSQLErrorWrap(gorm.ErrRecordNotFound)
	}
	return mapper.ProductRowToEntity(&row), nil
}

func (this *ProductRepositoryImpl) GetByIDs(ctx context.Context, ids []uint) ([]*entity.ProductEntity, error) {
	var rows []mapper.ProductRow
	err := this.db(ctx).Raw(
		`SELECT id, name, price, stock_quantity FROM products WHERE id = ANY(?::bigint[]) AND deleted_at IS NULL`,
		int64ArrayLiteral(toInt64Slice(ids)),
	).Scan(&rows).Error
	if err != nil {
		return nil, _errors.RawSQLErrorWrap(err)
	}

	result := make([]*entity.ProductEntity, len(rows))
	for i := range rows {
		result[i] = mapper.ProductRowToEntity(&rows[i])
	}
	return result, nil
}

func (this *ProductRepositoryImpl) List(ctx context.Context, filter query.ProductListFilter) (*entity.PagingEntity[*entity.ProductEntity], error) {
	var total int64
	if err := this.db(ctx).Raw(
		`SELECT count(*) FROM products WHERE deleted_at IS NULL`,
	).Scan(&total).Error; err != nil {
		return nil, _errors.RawSQLErrorWrap(err)
	}

	var rows []mapper.ProductRow
	if err := this.db(ctx).Raw(
		`SELECT id, name, price, stock_quantity FROM products
		 WHERE deleted_at IS NULL
		 ORDER BY id
		 LIMIT ? OFFSET ?`,
		filter.Limit, filter.Offset(),
	).Scan(&rows).Error; err != nil {
		return nil, _errors.RawSQLErrorWrap(err)
	}

	items := make([]*entity.ProductEntity, len(rows))
	for i := range rows {
		items[i] = mapper.ProductRowToEntity(&rows[i])
	}
	return entity.NewPagingEntity(filter.Page, filter.Limit, total, items), nil
}

func (this *ProductRepositoryImpl) ReserveStockBatchTx(ctx context.Context, tx unitofwork.Tx, reservations []repository.StockReservation) (map[uint]int64, []uint, error) {
	ids, quantities := toReservationSlices(reservations)

	var rows []mapper.ReservedProductRow
	err := this.tx(ctx, tx).Raw(
		`UPDATE products p
		 SET stock_quantity = p.stock_quantity - v.quantity, updated_at = now()
		 FROM (SELECT unnest(?::bigint[]) AS id, unnest(?::bigint[]) AS quantity) v
		 WHERE p.id = v.id AND p.stock_quantity >= v.quantity AND p.deleted_at IS NULL
		 RETURNING p.id, p.price`,
		int64ArrayLiteral(ids), int64ArrayLiteral(quantities),
	).Scan(&rows).Error
	if err != nil {
		return nil, nil, _errors.RawSQLErrorWrap(err)
	}

	prices := make(map[uint]int64, len(rows))
	for _, row := range rows {
		prices[row.ID] = row.Price
	}

	var failed []uint
	for _, r := range reservations {
		if _, ok := prices[r.ProductID]; !ok {
			failed = append(failed, r.ProductID)
		}
	}

	return prices, failed, nil
}

func (this *ProductRepositoryImpl) ReleaseStockBatchTx(ctx context.Context, tx unitofwork.Tx, reservations []repository.StockReservation) error {
	ids, quantities := toReservationSlices(reservations)

	result := this.tx(ctx, tx).Exec(
		`UPDATE products p
		 SET stock_quantity = p.stock_quantity + v.quantity, updated_at = now()
		 FROM (SELECT unnest(?::bigint[]) AS id, unnest(?::bigint[]) AS quantity) v
		 WHERE p.id = v.id`,
		int64ArrayLiteral(ids), int64ArrayLiteral(quantities),
	)
	return _errors.RawSQLErrorWrap(result.Error)
}

func toInt64Slice(ids []uint) []int64 {
	result := make([]int64, len(ids))
	for i, id := range ids {
		result[i] = int64(id)
	}
	return result
}

func toReservationSlices(reservations []repository.StockReservation) ([]int64, []int64) {
	ids := make([]int64, len(reservations))
	quantities := make([]int64, len(reservations))
	for i, r := range reservations {
		ids[i] = int64(r.ProductID)
		quantities[i] = r.Quantity
	}
	return ids, quantities
}

func int64ArrayLiteral(values []int64) string {
	parts := make([]string, len(values))
	for i, v := range values {
		parts[i] = strconv.FormatInt(v, 10)
	}
	return fmt.Sprintf("{%s}", strings.Join(parts, ","))
}
