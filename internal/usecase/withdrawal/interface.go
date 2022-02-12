package withdrawal

import (
	"context"
	"github.com/shopspring/decimal"
	"github.com/thorgnir-go-study/yp-diploma/internal/entity"
)

type Reader interface {
	List(ctx context.Context, userID entity.ID) ([]*entity.Withdrawal, error)
	GetSum(ctx context.Context, userID entity.ID) (decimal.NullDecimal, error)
}

type Writer interface {
	Create(ctx context.Context, withdrawal entity.Withdrawal) error
}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	CreateWithdrawal(ctx context.Context, userID entity.ID, orderNumber entity.OrderNumber, sum decimal.Decimal) error
	GetWithdrawals(ctx context.Context, userID entity.ID) ([]*entity.Withdrawal, error)
	GetWithdrawalsSum(ctx context.Context, userID entity.ID) (decimal.NullDecimal, error)
}
