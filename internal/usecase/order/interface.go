package order

import (
	"context"
	"github.com/shopspring/decimal"
	"github.com/thorgnir-go-study/yp-diploma/internal/entity"
)

type Reader interface {
	Get(ctx context.Context, orderID entity.ID) (*entity.Order, error)
	GetByOrderNumber(ctx context.Context, orderNumber entity.OrderNumber) (*entity.Order, error)
	List(ctx context.Context, userID entity.ID) ([]*entity.Order, error)
	GetAccrualsSum(ctx context.Context, userID entity.ID) (decimal.NullDecimal, error)
}
type Writer interface {
	Create(ctx context.Context, order entity.Order) (entity.ID, error)
	Update(ctx context.Context, order entity.Order) error
}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	CreateOrder(ctx context.Context, orderNumber entity.OrderNumber, userID entity.ID) (*entity.Order, error)
	GetUserOrders(ctx context.Context, userID entity.ID) ([]*entity.Order, error)
	SetOrderAccrualAndStatus(ctx context.Context, orderNumber entity.OrderNumber, accrual decimal.NullDecimal, status entity.OrderStatus) error
	GetAccrualsSum(ctx context.Context, userID entity.ID) (decimal.NullDecimal, error)
	SetOrderStatus(ctx context.Context, orderNumber entity.OrderNumber, status entity.OrderStatus) error
}
