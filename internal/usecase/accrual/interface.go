package accrual

import (
	"context"
	"github.com/thorgnir-go-study/yp-diploma/internal/entity"
)

type Reader interface {
	Get(ctx context.Context, orderNumber entity.OrderNumber) (*entity.AccrualOrder, error)
}

type Repository interface {
	Reader
}

type UseCase interface {
	GetAccrual(ctx context.Context, orderNumber entity.OrderNumber) (*entity.AccrualOrder, error)
}
