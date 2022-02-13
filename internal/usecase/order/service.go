package order

import (
	"context"
	"github.com/shopspring/decimal"
	"github.com/thorgnir-go-study/yp-diploma/internal/entity"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateOrder(ctx context.Context, orderNumber entity.OrderNumber, userID entity.ID) (*entity.Order, error) {
	order, err := entity.NewOrder(userID, orderNumber)
	if err != nil {
		return nil, err
	}
	_, err = s.repo.Create(ctx, *order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *Service) GetUserOrders(ctx context.Context, userID entity.ID) ([]*entity.Order, error) {
	return s.repo.List(ctx, userID)
}

func (s *Service) SetOrderAccrualAndStatus(ctx context.Context, orderID entity.ID, accrual decimal.NullDecimal, status entity.OrderStatus) error {
	return s.repo.SetOrderAccrualAndStatus(ctx, orderID, accrual, status)
}

func (s *Service) GetAccrualsSum(ctx context.Context, userID entity.ID) (decimal.NullDecimal, error) {
	return s.repo.GetAccrualsSum(ctx, userID)
}

func (s *Service) SetOrderStatus(ctx context.Context, orderID entity.ID, status entity.OrderStatus) error {
	return s.repo.SetOrderStatus(ctx, orderID, status)
}

func (s *Service) GetNewOrders(ctx context.Context) ([]*entity.Order, error) {
	return s.repo.GetNewOrders(ctx)
}
