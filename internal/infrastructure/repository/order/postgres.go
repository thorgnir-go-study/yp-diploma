package order

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/thorgnir-go-study/yp-diploma/internal/entity"
	"github.com/thorgnir-go-study/yp-diploma/internal/usecase/order"
)

type PostgresOrderRepository struct {
	dbpool *pgxpool.Pool
}

func NewPostgresOrderRepository(dbpool *pgxpool.Pool) order.Repository {
	return &PostgresOrderRepository{dbpool: dbpool}
}

func (p PostgresOrderRepository) Get(ctx context.Context, orderID entity.ID) (*entity.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresOrderRepository) GetByOrderNumber(ctx context.Context, orderNumber entity.OrderNumber) (*entity.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresOrderRepository) List(ctx context.Context, userID entity.ID) ([]*entity.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresOrderRepository) GetAccrualsSum(ctx context.Context, userID entity.ID) (decimal.NullDecimal, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresOrderRepository) GetNewOrders(ctx context.Context) ([]*entity.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresOrderRepository) Create(ctx context.Context, order entity.Order) (entity.ID, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresOrderRepository) SetOrderAccrualAndStatus(ctx context.Context, orderID entity.ID, accrual decimal.NullDecimal, status entity.OrderStatus) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresOrderRepository) SetOrderStatus(ctx context.Context, orderID entity.ID, status entity.OrderStatus) error {
	//TODO implement me
	panic("implement me")
}
