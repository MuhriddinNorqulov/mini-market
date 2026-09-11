package orderusecases_test

import (
	"context"
	"testing"

	"mini-market/src/core/application/response"
	"mini-market/src/core/application/usecases/orderusecases"
	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/entity/enum"
)

func TestGetOrderUseCase_OwnerCanViewOwnOrder(t *testing.T) {
	orderRepo := newFakeOrderRepository()
	orderRepo.orders[1] = &entity.OrderEntity{ID: 1, UserID: 1, Status: enum.OrderStatusPending}

	uc := orderusecases.NewGetOrderUseCase(orderRepo)

	order, err := uc.Invoke(context.Background(), &entity.UserEntity{ID: 1, Role: enum.RoleUser}, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if order.ID != 1 {
		t.Fatalf("expected order 1, got %d", order.ID)
	}
}

func TestGetOrderUseCase_NonOwnerNonAdminIsForbidden(t *testing.T) {
	orderRepo := newFakeOrderRepository()
	orderRepo.orders[1] = &entity.OrderEntity{ID: 1, UserID: 1, Status: enum.OrderStatusPending}

	uc := orderusecases.NewGetOrderUseCase(orderRepo)

	_, err := uc.Invoke(context.Background(), &entity.UserEntity{ID: 2, Role: enum.RoleUser}, 1)
	if !response.IsErrorCode(err, response.CodeForbidden) {
		t.Fatalf("expected CodeForbidden, got %v", err)
	}
}

func TestGetOrderUseCase_AdminCanViewAnyOrder(t *testing.T) {
	orderRepo := newFakeOrderRepository()
	orderRepo.orders[1] = &entity.OrderEntity{ID: 1, UserID: 1, Status: enum.OrderStatusPending}

	uc := orderusecases.NewGetOrderUseCase(orderRepo)

	order, err := uc.Invoke(context.Background(), &entity.UserEntity{ID: 2, Role: enum.RoleAdmin}, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if order.ID != 1 {
		t.Fatalf("expected order 1, got %d", order.ID)
	}
}
