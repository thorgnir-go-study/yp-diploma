package accrual_processor

import (
	"context"
	"github.com/thorgnir-go-study/yp-diploma/internal/entity"
	"github.com/thorgnir-go-study/yp-diploma/internal/usecase/accrual"
	"github.com/thorgnir-go-study/yp-diploma/internal/usecase/order"
)

type Service struct {
	accrualUC accrual.UseCase
	orderUC   order.UseCase
}

func NewService(accrualUC accrual.UseCase, orderUC order.UseCase) *Service {
	return &Service{accrualUC: accrualUC, orderUC: orderUC}
}

func (s *Service) AddOrderToProcessingQueue(ctx context.Context, order *entity.Order) error {
	//TODO implement me
	err := s.orderUC.SetOrderStatus(ctx, order.Number, entity.OrderStatusProcessing)
	if err != nil {
		return err
	}

	return nil
}
