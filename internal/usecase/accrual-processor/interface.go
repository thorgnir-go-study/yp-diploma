package accrual_processor

import (
	"context"
	"github.com/thorgnir-go-study/yp-diploma/internal/entity"
)

type UseCase interface {
	AddOrderToProcessingQueue(ctx context.Context, order *entity.Order) error
}
