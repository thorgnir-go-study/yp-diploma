package withdrawal

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/thorgnir-go-study/yp-diploma/internal/entity"
	"github.com/thorgnir-go-study/yp-diploma/internal/usecase/withdrawal"
)

type PostgresWithdrawalRepository struct {
	dbpool *pgxpool.Pool
}

func NewPostgresWithdrawalRepository(dbpool *pgxpool.Pool) withdrawal.Repository {
	return &PostgresWithdrawalRepository{dbpool: dbpool}
}

func (p PostgresWithdrawalRepository) List(ctx context.Context, userID entity.ID) ([]*entity.Withdrawal, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresWithdrawalRepository) GetSum(ctx context.Context, userID entity.ID) (decimal.NullDecimal, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresWithdrawalRepository) Create(ctx context.Context, withdrawal entity.Withdrawal) error {
	//TODO implement me
	panic("implement me")
}
