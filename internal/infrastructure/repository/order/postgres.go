package order

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
	"github.com/thorgnir-go-study/yp-diploma/internal/usecase/order"
	"time"
)

type PostgresOrderRepository struct {
	dbpool *pgxpool.Pool
}

type dbEntity struct {
	ID         uuid.UUID           `db:"id"`
	Number     string              `db:"order_number"`
	UserID     uuid.UUID           `db:"user_id"`
	StatusID   int                 `db:"status_id"`
	Accrual    decimal.NullDecimal `db:"accrual"`
	UploadedAt time.Time           `db:"uploaded_at"`
	UpdatedAt  time.Time           `db:"updated_at"`
}

func NewPostgresOrderRepository(dbpool *pgxpool.Pool) order.Repository {
	return &PostgresOrderRepository{dbpool: dbpool}
}

func (p PostgresOrderRepository) List(ctx context.Context, userID entity.ID) ([]*entity.Order, error) {
	var orders []*dbEntity
	log.Info().Str("currentUser", userID.String()).Msg("CurrentUser")
	if err := pgxscan.Select(ctx, p.dbpool, &orders, `
SELECT id, order_number, user_id, status_id, accrual, uploaded_at, updated_at 
FROM gophermart.order 
WHERE user_id = $1 
ORDER BY uploaded_at ASC
`, userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		log.Error().Err(err).Msg("Error while getting user orders")
		return nil, err
	}

	result := make([]*entity.Order, len(orders))
	for i := range orders {
		result[i] = mapOrder(orders[i])
	}

	return result, nil

}

func (p PostgresOrderRepository) GetAccrualsSum(ctx context.Context, userID entity.ID) (decimal.NullDecimal, error) {
	var accrual decimal.NullDecimal
	if err := p.dbpool.QueryRow(ctx, `SELECT SUM(accrual) FROM gophermart.order WHERE user_id = $1`, userID).Scan(&accrual); err != nil {
		return accrual, err
	}

	return accrual, nil
}

func (p PostgresOrderRepository) GetNewOrders(ctx context.Context) ([]*entity.Order, error) {
	var orders []*dbEntity
	if err := pgxscan.Select(ctx, p.dbpool, &orders, `
SELECT id, order_number, user_id, status_id, accrual, uploaded_at, updated_at 
FROM gophermart.order 
WHERE status_id = $1
`, entity.OrderStatusNew); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		log.Error().Err(err).Msg("Error while getting new orders")
		return nil, err
	}

	result := make([]*entity.Order, len(orders))
	for i := range orders {
		result[i] = mapOrder(orders[i])
	}

	return result, nil
}

func (p PostgresOrderRepository) Create(ctx context.Context, order entity.Order) (entity.ID, error) {
	var insertedOrderID uuid.UUID
	var newOrderUserID uuid.UUID
	var inserted bool

	if err := p.dbpool.QueryRow(ctx, `
WITH new_order AS (
    INSERT INTO gophermart."order"(
	id, order_number, user_id, status_id, accrual, uploaded_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
    ON CONFLICT(order_number) DO NOTHING
    RETURNING id, user_id, true as inserted
) SELECT * from new_order 
UNION 
SELECT id, user_id, false as inserted FROM gophermart."order" 
WHERE order_number = $2
`, order.ID, order.Number.String(), order.UserID, order.Status, order.Accrual, order.UploadedAt, order.UpdatedAt).
		Scan(&insertedOrderID, &newOrderUserID, &inserted); err != nil {
		return entity.NilID, err
	}
	if !inserted {
		if order.UserID == newOrderUserID {
			return entity.NilID, entity.ErrOrderAlreadyRegistered
		}
		return entity.NilID, entity.ErrOrderRegisteredByAnotherUser
	}
	return insertedOrderID, nil
}

func (p PostgresOrderRepository) SetOrderAccrualAndStatus(ctx context.Context, orderID entity.ID, accrual decimal.NullDecimal, status entity.OrderStatus) error {
	if _, err := p.dbpool.Exec(ctx,
		`UPDATE gophermart."order" SET accrual = $1, status_id = $2, updated_at = $3 WHERE id = $4`,
		accrual, status, time.Now(), orderID); err != nil {
		log.Error().Err(err).
			Str("orderId", orderID.String()).
			Str("status", status.String()).
			Str("accrual", accrual.Decimal.String()).
			Msg("Error while setting accrual and status")
		return err
	}

	return nil
}

func (p PostgresOrderRepository) SetOrderStatus(ctx context.Context, orderID entity.ID, status entity.OrderStatus) error {
	if _, err := p.dbpool.Exec(ctx,
		`UPDATE gophermart."order" SET status_id = $1, updated_at = $2 WHERE id = $3`,
		status, time.Now(), orderID); err != nil {
		log.Error().Err(err).
			Str("orderId", orderID.String()).
			Str("status", status.String()).
			Msg("Error while setting order status")
		return err
	}

	return nil
}

func mapOrder(dbOrder *dbEntity) *entity.Order {
	orderNumber, _ := entity.StringToOrderNumber(dbOrder.Number)
	orderStatus := entity.OrderStatus(dbOrder.StatusID)
	return &entity.Order{
		ID:         dbOrder.ID,
		UserID:     dbOrder.UserID,
		Number:     orderNumber,
		Status:     orderStatus,
		Accrual:    dbOrder.Accrual,
		UploadedAt: dbOrder.UploadedAt,
		UpdatedAt:  dbOrder.UpdatedAt,
	}
}
