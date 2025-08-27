package service

import (
	"L0-wb/config"
	"L0-wb/internal/models"
	"L0-wb/internal/repo"
	"context"
	"sync"
)

type UserService struct {
	UserRepo *repo.PostgresRepo
	cache    map[string]*models.Order
	mu       sync.RWMutex
}

func NewService(ur *repo.PostgresRepo) (*UserService, error) {
	s := &UserService{
		UserRepo: ur,
		cache:    make(map[string]*models.Order),
	}
	if err := s.RestoreCache(context.Background()); err != nil {
		return nil, err
	}
	return s, nil
}

// загрузка кэша из БД
func (s *UserService) RestoreCache(ctx context.Context) error {
	last_order_quantity := config.GetLimitCache()
	orders, err := s.UserRepo.GetLastOrders(ctx, last_order_quantity)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, order := range orders {
		orderCopy := order
		s.cache[order.OrderUID] = &orderCopy
	}
	return nil
}

// GetOrder возвращает заказ по ID сначала из кэша, если в кэше нет то делает запрос к БД
func (s *UserService) GetOrder(ctx context.Context, orderUID string) (*models.Order, error) {
	s.mu.RLock()
	order, ok := s.cache[orderUID]
	s.mu.RUnlock()
	if ok {
		return order, nil
	}
	orderDB, err := s.UserRepo.GetOrder(ctx, orderUID)
	if err != nil {
		return nil, err
	}
	s.mu.Lock()
	orderCopy := orderDB
	s.cache[orderUID] = &orderCopy
	s.mu.Unlock()
	return &orderCopy, nil
}

// GetOrderResponse возвращает структуру заказа заказа для пользователя без лишних полей
func (s *UserService) GetOrderResponse(ctx context.Context, orderUID string) (*models.OrderResponse, error) {
	order, err := s.GetOrder(ctx, orderUID)
	if err != nil {
		return nil, err
	}
	return s.convertToOrderResponse(order), nil
}

// CreateOrder сохраняет заказ в БД и кэш
func (s *UserService) CreateOrder(ctx context.Context, order *models.Order) error {
	if err := s.UserRepo.CreateOrder(ctx, *order); err != nil {
		return err
	}
	s.mu.Lock()
	s.cache[order.OrderUID] = order
	s.mu.Unlock()
	return nil
}

// convertToOrderResponse конвертирует Order в OrderResponse
func (s *UserService) convertToOrderResponse(order *models.Order) *models.OrderResponse {
	itemsResponse := make(models.ItemsResponse, len(order.Items))
	for i, item := range order.Items {
		itemsResponse[i] = models.ItemResponse{
			TrackNumber: item.TrackNumber,
			Price:       item.Price,
			Name:        item.Name,
			Sale:        item.Sale,
			Size:        item.Size,
			TotalPrice:  item.TotalPrice,
			Brand:       item.Brand,
		}
	}

	return &models.OrderResponse{
		OrderUID:    order.OrderUID,
		TrackNumber: order.TrackNumber,
		Delivery:    order.Delivery,
		Payment: models.PaymentResponse{
			Currency:     order.Payment.Currency,
			Provider:     order.Payment.Provider,
			Amount:       order.Payment.Amount,
			PaymentDt:    order.Payment.PaymentDt,
			Bank:         order.Payment.Bank,
			DeliveryCost: order.Payment.DeliveryCost,
			GoodsTotal:   order.Payment.GoodsTotal,
			CustomFee:    order.Payment.CustomFee,
		},
		Items:           itemsResponse,
		Locale:          order.Locale,
		DeliveryService: order.DeliveryService,
		DateCreated:     order.DateCreated,
	}
}
