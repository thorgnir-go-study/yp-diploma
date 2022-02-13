package withdrawal

import (
	"context"
	"errors"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"github.com/thorgnir-go-study/yp-diploma/internal/entity"
	"github.com/thorgnir-go-study/yp-diploma/internal/usecase/withdrawal"
	"time"
)

type PostgresWithdrawalRepository struct {
	dbpool *pgxpool.Pool
}

type dbEntity struct {
	ID          uuid.UUID       `db:"id"`
	UserID      uuid.UUID       `db:"user_id"`
	OrderNumber string          `db:"order_number"`
	Sum         decimal.Decimal `db:"sum"`
	ProcessedAt time.Time       `db:"processed_at"`
}

func NewPostgresWithdrawalRepository(dbpool *pgxpool.Pool) withdrawal.Repository {
	return &PostgresWithdrawalRepository{dbpool: dbpool}
}

func (p PostgresWithdrawalRepository) List(ctx context.Context, userID entity.ID) ([]*entity.Withdrawal, error) {
	var withdrawals []*dbEntity
	if err := pgxscan.Select(ctx, p.dbpool, &withdrawals, `
SELECT id, user_id, order_number, "sum", processed_at
FROM gophermart."withdrawal"
WHERE user_id = $1
`, userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		log.Error().Err(err).Msg("Error while getting withdrawals")
		return nil, err
	}
	result := make([]*entity.Withdrawal, len(withdrawals))
	for i := range withdrawals {
		orderNumber, _ := entity.StringToOrderNumber(withdrawals[i].OrderNumber)
		result[i] = &entity.Withdrawal{
			ID:          withdrawals[i].ID,
			UserID:      withdrawals[i].UserID,
			OrderNumber: orderNumber,
			Sum:         withdrawals[i].Sum,
			ProcessedAt: withdrawals[i].ProcessedAt,
		}
	}
	return result, nil
}

func (p PostgresWithdrawalRepository) GetSum(ctx context.Context, userID entity.ID) (decimal.NullDecimal, error) {
	var sum decimal.NullDecimal
	if err := p.dbpool.QueryRow(ctx, `SELECT SUM("sum") FROM gophermart."withdrawal" WHERE user_id = $1`, userID).Scan(&sum); err != nil {
		return sum, err
	}
	return sum, nil
}

func (p PostgresWithdrawalRepository) Create(ctx context.Context, withdrawal entity.Withdrawal) error {
	if _, err := p.dbpool.Exec(ctx, `insert into gophermart."withdrawal" 
(id, order_number, "sum", processed_at, user_id)
values ($1, $2, $3, $4, $5)`, withdrawal.ID, withdrawal.OrderNumber, withdrawal.Sum, withdrawal.ProcessedAt, withdrawal.UserID); err != nil {
		return err
	}
	return nil
}
