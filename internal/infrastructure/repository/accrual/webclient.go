package accrual

import (
	"context"
	"github.com/thorgnir-go-study/yp-diploma/internal/entity"
	"github.com/thorgnir-go-study/yp-diploma/internal/usecase/accrual"
)

type WebClient struct {
}

func NewWebClient() accrual.Repository {
	return &WebClient{}
}

func (w WebClient) Get(ctx context.Context, orderNumber entity.OrderNumber) (*entity.AccrualOrder, error) {
	//TODO implement me
	panic("implement me")
}
