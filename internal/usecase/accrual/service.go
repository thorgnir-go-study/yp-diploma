package accrual

import (
	"context"
	"github.com/thorgnir-go-study/yp-diploma/internal/entity"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAccrual(ctx context.Context, orderNumber entity.OrderNumber) (*entity.AccrualOrder, error) {
	return s.repo.Get(ctx, orderNumber)
}
