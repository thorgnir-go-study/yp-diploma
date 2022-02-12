package order

import (
	"context"
	"github.com/shopspring/decimal"
	"github.com/thorgnir-go-study/yp-diploma/internal/entity"
	"time"
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

func (s *Service) SetOrderAccrualAndStatus(ctx context.Context, orderNumber entity.OrderNumber, accrual decimal.NullDecimal, status entity.OrderStatus) error {
	order, err := s.repo.GetByOrderNumber(ctx, orderNumber)
	if err != nil {
		return err
	}
	order.UpdatedAt = time.Now()
	order.Accrual = accrual
	order.Status = status
	err = order.Validate()
	if err != nil {
		return err
	}
	err = s.repo.Update(ctx, *order)
	return err
}

func (s *Service) GetAccrualsSum(ctx context.Context, userID entity.ID) (decimal.NullDecimal, error) {
	return s.repo.GetAccrualsSum(ctx, userID)
}

func (s *Service) SetOrderStatus(ctx context.Context, orderNumber entity.OrderNumber, status entity.OrderStatus) error {
	order, err := s.repo.GetByOrderNumber(ctx, orderNumber)
	if err != nil {
		return err
	}
	order.UpdatedAt = time.Now()
	order.Status = status
	err = order.Validate()
	if err != nil {
		return err
	}
	err = s.repo.Update(ctx, *order)
	return err
}
