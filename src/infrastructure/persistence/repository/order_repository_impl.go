package repository

import (
	"context"

	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/entity/enum"
	"mini-market/src/core/domain/ports/repository"
	"mini-market/src/core/domain/ports/unitofwork"
	"mini-market/src/core/domain/query"
	"mini-market/src/infrastructure/_errors"
	"mini-market/src/infrastructure/persistence/mapper"

	"gorm.io/gorm"
)

type OrderRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewOrderRepositoryImpl(baseRepository *BaseRepository) repository.OrderRepository {
	return &OrderRepositoryImpl{BaseRepository: baseRepository}
}

func (this *OrderRepositoryImpl) CreateTx(ctx context.Context, tx unitofwork.Tx, order *entity.OrderEntity) (*entity.OrderEntity, error) {
	db := this.tx(ctx, tx)

	var orderRow mapper.OrderRow
	err := db.Raw(
		`INSERT INTO orders (user_id, idempotency_key, total_price, status, expires_at, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, now(), now())
		 RETURNING id, user_id, idempotency_key, total_price, status, expires_at, created_at`,
		order.UserID, order.IdempotencyKey, order.TotalPrice, order.Status, order.ExpiresAt,
	).Scan(&orderRow).Error
	if err != nil {
		return nil, _errors.RawSQLErrorWrap(err)
	}

	items := make([]*entity.OrderItemEntity, 0, len(order.Items))
	for _, item := range order.Items {
		var itemRow mapper.OrderItemRow
		err := db.Raw(
			`INSERT INTO order_items (order_id, product_id, quantity, unit_price, created_at, updated_at)
			 VALUES (?, ?, ?, ?, now(), now())
			 RETURNING id, product_id, quantity, unit_price`,
			orderRow.ID, item.ProductID, item.Quantity, item.UnitPrice,
		).Scan(&itemRow).Error
		if err != nil {
			return nil, _errors.RawSQLErrorWrap(err)
		}
		items = append(items, mapper.OrderItemRowToEntity(&itemRow))
	}

	result := mapper.OrderRowToEntity(&orderRow)
	result.Items = items
	return result, nil
}

func (this *OrderRepositoryImpl) GetByIdempotencyKey(ctx context.Context, userID uint, idempotencyKey string) (*entity.OrderEntity, error) {
	db := this.db(ctx)
	var row mapper.OrderRow
	err := db.Raw(
		`SELECT id, user_id, idempotency_key, total_price, status, expires_at, created_at
		 FROM orders WHERE user_id = ? AND idempotency_key = ? AND deleted_at IS NULL`,
		userID, idempotencyKey,
	).Scan(&row).Error
	if err != nil {
		return nil, _errors.RawSQLErrorWrap(err)
	}
	if row.ID == 0 {
		return nil, _errors.RawSQLErrorWrap(gorm.ErrRecordNotFound)
	}
	return this.attachItems(db, &row)
}

func (this *OrderRepositoryImpl) GetByIdempotencyKeyTx(ctx context.Context, tx unitofwork.Tx, userID uint, idempotencyKey string) (*entity.OrderEntity, error) {
	db := this.tx(ctx, tx)
	var row mapper.OrderRow
	err := db.Raw(
		`SELECT id, user_id, idempotency_key, total_price, status, expires_at, created_at
		 FROM orders WHERE user_id = ? AND idempotency_key = ? AND deleted_at IS NULL`,
		userID, idempotencyKey,
	).Scan(&row).Error
	if err != nil {
		return nil, _errors.RawSQLErrorWrap(err)
	}
	if row.ID == 0 {
		return nil, _errors.RawSQLErrorWrap(gorm.ErrRecordNotFound)
	}
	return this.attachItems(db, &row)
}

func (this *OrderRepositoryImpl) GetByID(ctx context.Context, id uint) (*entity.OrderEntity, error) {
	db := this.db(ctx)
	var row mapper.OrderRow
	err := db.Raw(
		`SELECT id, user_id, idempotency_key, total_price, status, expires_at, created_at
		 FROM orders WHERE id = ? AND deleted_at IS NULL`,
		id,
	).Scan(&row).Error
	if err != nil {
		return nil, _errors.RawSQLErrorWrap(err)
	}
	if row.ID == 0 {
		return nil, _errors.RawSQLErrorWrap(gorm.ErrRecordNotFound)
	}
	return this.attachItems(db, &row)
}

func (this *OrderRepositoryImpl) GetForUpdateTx(ctx context.Context, tx unitofwork.Tx, id uint) (*entity.OrderEntity, error) {
	db := this.tx(ctx, tx)
	var row mapper.OrderRow
	err := db.Raw(
		`SELECT id, user_id, idempotency_key, total_price, status, expires_at, created_at
		 FROM orders WHERE id = ? AND deleted_at IS NULL FOR UPDATE`,
		id,
	).Scan(&row).Error
	if err != nil {
		return nil, _errors.RawSQLErrorWrap(err)
	}
	if row.ID == 0 {
		return nil, _errors.RawSQLErrorWrap(gorm.ErrRecordNotFound)
	}
	return this.attachItems(db, &row)
}

func (this *OrderRepositoryImpl) UpdateStatusTx(ctx context.Context, tx unitofwork.Tx, id uint, status enum.OrderStatus, cancelReason *enum.OrderCancelReason) error {
	result := this.tx(ctx, tx).Exec(
		`UPDATE orders SET status = ?, cancel_reason = ?,
		 completed_at = CASE WHEN ? = 'confirmed' THEN now() ELSE completed_at END,
		 cancelled_at = CASE WHEN ? = 'cancelled' THEN now() ELSE cancelled_at END,
		 updated_at = now()
		 WHERE id = ?`,
		status, cancelReason, status, status, id,
	)
	return _errors.RawSQLErrorWrap(result.Error)
}

func (this *OrderRepositoryImpl) List(ctx context.Context, filter query.OrderListFilter) (*entity.PagingEntity[*entity.OrderEntity], error) {
	db := this.db(ctx)

	var total int64
	if err := db.Raw(
		`SELECT count(*) FROM orders WHERE deleted_at IS NULL`,
	).Scan(&total).Error; err != nil {
		return nil, _errors.RawSQLErrorWrap(err)
	}

	var rows []mapper.OrderRow
	if err := db.Raw(
		`SELECT id, user_id, idempotency_key, total_price, status, expires_at, created_at
		 FROM orders
		 WHERE deleted_at IS NULL
		 ORDER BY created_at DESC
		 LIMIT ? OFFSET ?`,
		filter.Limit, filter.Offset(),
	).Scan(&rows).Error; err != nil {
		return nil, _errors.RawSQLErrorWrap(err)
	}

	orderIDs := make([]int64, len(rows))
	for i, row := range rows {
		orderIDs[i] = int64(row.ID)
	}

	var itemRows []struct {
		OrderID uint
		mapper.OrderItemRow
	}
	if len(orderIDs) > 0 {
		if err := db.Raw(
			`SELECT order_id, id, product_id, quantity, unit_price
			 FROM order_items
			 WHERE order_id = ANY(?::bigint[]) AND deleted_at IS NULL`,
			int64ArrayLiteral(orderIDs),
		).Scan(&itemRows).Error; err != nil {
			return nil, _errors.RawSQLErrorWrap(err)
		}
	}

	itemsByOrder := make(map[uint][]*entity.OrderItemEntity, len(rows))
	for i := range itemRows {
		itemsByOrder[itemRows[i].OrderID] = append(itemsByOrder[itemRows[i].OrderID], mapper.OrderItemRowToEntity(&itemRows[i].OrderItemRow))
	}

	items := make([]*entity.OrderEntity, len(rows))
	for i := range rows {
		order := mapper.OrderRowToEntity(&rows[i])
		order.Items = itemsByOrder[order.ID]
		items[i] = order
	}

	return entity.NewPagingEntity(filter.Page, filter.Limit, total, items), nil
}

func (this *OrderRepositoryImpl) ListByUserID(ctx context.Context, userID uint, filter query.OrderListFilter) (*entity.PagingEntity[*entity.OrderEntity], error) {
	db := this.db(ctx)

	var total int64
	if err := db.Raw(
		`SELECT count(*) FROM orders WHERE user_id = ? AND deleted_at IS NULL`,
		userID,
	).Scan(&total).Error; err != nil {
		return nil, _errors.RawSQLErrorWrap(err)
	}

	var rows []mapper.OrderRow
	if err := db.Raw(
		`SELECT id, user_id, idempotency_key, total_price, status, expires_at, created_at
		 FROM orders
		 WHERE user_id = ? AND deleted_at IS NULL
		 ORDER BY created_at DESC
		 LIMIT ? OFFSET ?`,
		userID, filter.Limit, filter.Offset(),
	).Scan(&rows).Error; err != nil {
		return nil, _errors.RawSQLErrorWrap(err)
	}

	orderIDs := make([]int64, len(rows))
	for i, row := range rows {
		orderIDs[i] = int64(row.ID)
	}

	var itemRows []struct {
		OrderID uint
		mapper.OrderItemRow
	}
	if len(orderIDs) > 0 {
		if err := db.Raw(
			`SELECT order_id, id, product_id, quantity, unit_price
			 FROM order_items
			 WHERE order_id = ANY(?::bigint[]) AND deleted_at IS NULL`,
			int64ArrayLiteral(orderIDs),
		).Scan(&itemRows).Error; err != nil {
			return nil, _errors.RawSQLErrorWrap(err)
		}
	}

	itemsByOrder := make(map[uint][]*entity.OrderItemEntity, len(rows))
	for i := range itemRows {
		itemsByOrder[itemRows[i].OrderID] = append(itemsByOrder[itemRows[i].OrderID], mapper.OrderItemRowToEntity(&itemRows[i].OrderItemRow))
	}

	items := make([]*entity.OrderEntity, len(rows))
	for i := range rows {
		order := mapper.OrderRowToEntity(&rows[i])
		order.Items = itemsByOrder[order.ID]
		items[i] = order
	}

	return entity.NewPagingEntity(filter.Page, filter.Limit, total, items), nil
}

func (this *OrderRepositoryImpl) attachItems(db *gorm.DB, row *mapper.OrderRow) (*entity.OrderEntity, error) {
	var itemRows []mapper.OrderItemRow
	err := db.Raw(
		`SELECT id, product_id, quantity, unit_price FROM order_items WHERE order_id = ? AND deleted_at IS NULL`,
		row.ID,
	).Scan(&itemRows).Error
	if err != nil {
		return nil, _errors.RawSQLErrorWrap(err)
	}

	result := mapper.OrderRowToEntity(row)
	result.Items = make([]*entity.OrderItemEntity, 0, len(itemRows))
	for i := range itemRows {
		result.Items = append(result.Items, mapper.OrderItemRowToEntity(&itemRows[i]))
	}
	return result, nil
}
