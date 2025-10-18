package service

import (
	"L0-wb/config"
	"L0-wb/internal/cache"
	"L0-wb/internal/models"
	"L0-wb/internal/repo"
	"context"
	"fmt"
)

type UserService struct {
	UserRepo repo.Repository
	cache    cache.Cache
}

func NewService(ur repo.Repository) (Service, error) {
	maxSize := config.GetCacheStartupSize()
	s := &UserService{
		UserRepo: ur,
		cache:    cache.NewCache(maxSize),
	}

	if err := s.RestoreCache(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to restore cache: %w", err)
	}

	return s, nil
}

func (s *UserService) GetOrderByUID(ctx context.Context, orderUID string) (*models.Order, error) {
	if orderUID == "" {
		return nil, fmt.Errorf("order_uid cannot be empty")
	}

	// Check cache first
	if order, found := s.cache.Get(orderUID); found {
		return order, nil
	}

	// Get from DB if not in cache
	orderDB, err := s.UserRepo.GetOrder(ctx, orderUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Save to cache
	s.cache.Set(orderUID, &orderDB)

	return &orderDB, nil
}

func (s *UserService) GetOrderResponse(ctx context.Context, orderUID string) (*models.OrderResponse, error) {
	order, err := s.GetOrderByUID(ctx, orderUID)
	if err != nil {
		return nil, err
	}
	return order.ConvertToOrderResponse(), nil
}

func (s *UserService) CreateOrder(ctx context.Context, order *models.Order) error {
	if err := order.Validate(); err != nil {
		return fmt.Errorf("invalid order: %w", err)
	}

	if err := s.UserRepo.CreateOrder(ctx, *order); err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	s.cache.Set(order.OrderUID, order)
	return nil
}

func (s *UserService) SaveOrder(ctx context.Context, order *models.Order) error {
	return s.CreateOrder(ctx, order)
}

func (s *UserService) RestoreCache(ctx context.Context) error {
	lastOrderQuantity := config.GetLimitCache()
	orders, err := s.UserRepo.GetLastOrders(ctx, lastOrderQuantity)
	if err != nil {
		return fmt.Errorf("failed to restore cache: %w", err)
	}

	for _, order := range orders {
		orderCopy := order
		s.cache.Set(order.OrderUID, &orderCopy)
	}
	return nil
}

func (s *UserService) Close() error {
	if s.cache != nil {
		s.cache.Close()
	}
	return nil
}
