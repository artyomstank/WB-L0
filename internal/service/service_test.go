package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"L0-wb/internal/mocks"
	"L0-wb/internal/models"
)

func TestService_GetOrderByUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockCache := mocks.NewMockCache(ctrl)

	svc := &UserService{
		UserRepo: mockRepo,
		cache:    mockCache,
	}

	testOrder := &models.Order{
		OrderUID:    "test-123",
		DateCreated: time.Now(),
	}

	tests := []struct {
		name    string
		uid     string
		setup   func()
		want    *models.Order
		wantErr error
	}{{
		name:    "get from cache",
		uid:     "test-123",
		setup:   func() { mockCache.EXPECT().Get("test-123").Return(testOrder, true) },
		want:    testOrder,
		wantErr: nil,
	}, {
		name: "get from db",
		uid:  "test-123",
		setup: func() {
			mockCache.EXPECT().Get("test-123").Return(nil, false)
			mockRepo.EXPECT().GetOrder(gomock.Any(), "test-123").Return(*testOrder, nil)
			mockCache.EXPECT().Set("test-123", testOrder)
		},
		want:    testOrder,
		wantErr: nil,
	}, {
		name: "not found",
		uid:  "not-exists",
		setup: func() {
			mockCache.EXPECT().Get("not-exists").Return(nil, false)
			mockRepo.EXPECT().GetOrder(gomock.Any(), "not-exists").Return(models.Order{}, sql.ErrNoRows)
		},
		want:    nil,
		wantErr: ErrNotFound,
	}, {
		name:    "empty uid",
		uid:     "",
		setup:   func() {},
		want:    nil,
		wantErr: errors.New("order_uid cannot be empty"),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			got, err := svc.GetOrderByUID(context.Background(), tt.uid)
			if tt.wantErr != nil {
				assert.Error(t, err)
				if tt.wantErr == ErrNotFound {
					assert.ErrorIs(t, err, ErrNotFound)
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserService_CreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockCache := mocks.NewMockCache(ctrl)
	svc := &UserService{UserRepo: mockRepo, cache: mockCache}

	validOrder := &models.Order{
		OrderUID:    "test-123",
		TrackNumber: "WBIL123456789",
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:    "Test User",
			Phone:   "+79991234567",
			Zip:     "123456",
			City:    "Moscow",
			Address: "Test St, 1",
		},
		Payment: models.Payment{
			Transaction: "tx-123",
			Currency:    "USD",
			Provider:    "stripe",
			Amount:      1000,
		},
		Items: []models.Item{
			{
				ChrtID:      1,
				TrackNumber: "WBIL123456789",
				Price:       100,
				Name:        "Test Item",
				TotalPrice:  100,
				Status:      200,
			},
		},
	}

	t.Run("success create and cache", func(t *testing.T) {
		mockRepo.EXPECT().
			CreateOrder(gomock.Any(), *validOrder).
			Return(nil)
		mockCache.EXPECT().
			Set(validOrder.OrderUID, validOrder)

		err := svc.CreateOrder(context.Background(), validOrder)
		assert.NoError(t, err)
	})

	invalidOrder := &models.Order{} // Пустой заказ для теста валидации
	t.Run("validation error", func(t *testing.T) {
		err := svc.CreateOrder(context.Background(), invalidOrder)
		assert.Error(t, err)
	})
}

func TestUserService_RestoreCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockCache := mocks.NewMockCache(ctrl)
	svc := &UserService{UserRepo: mockRepo, cache: mockCache}

	orders := []models.Order{
		{OrderUID: "1"},
		{OrderUID: "2"},
	}

	t.Run("success restore", func(t *testing.T) {
		mockRepo.EXPECT().
			GetLastOrders(gomock.Any(), gomock.Any()).
			Return(orders, nil)

		for _, order := range orders {
			orderCopy := order
			mockCache.EXPECT().
				Set(order.OrderUID, &orderCopy)
		}

		err := svc.RestoreCache(context.Background())
		assert.NoError(t, err)
	})
}
