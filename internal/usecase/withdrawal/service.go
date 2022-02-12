package withdrawal

import (
	"context"
	"github.com/shopspring/decimal"
	"github.com/thorgnir-go-study/yp-diploma/internal/entity"
	"github.com/thorgnir-go-study/yp-diploma/internal/usecase/order"
	"sync"
)

type Service struct {
	repo    Repository
	orderUC order.UseCase
	userMX  sync.Map
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateWithdrawal(ctx context.Context, userID entity.ID, orderNumber entity.OrderNumber, sum decimal.Decimal) error {
	userLock := s.getUserMutex(userID)
	userLock.Lock()
	defer userLock.Unlock()

	dbWithdrawalsSum, err := s.repo.GetSum(ctx, userID)
	if err != nil {
		return err
	}

	withdrawalSum := decimal.NewFromInt(0)
	if dbWithdrawalsSum.Valid {
		withdrawalSum = dbWithdrawalsSum.Decimal
	}

	accrualsSum, err := s.orderUC.GetAccrualsSum(ctx, userID)
	if err != nil {
		return err
	}

	if !accrualsSum.Valid {
		return entity.ErrInsufficientBalance
	}

	balance := accrualsSum.Decimal.Sub(withdrawalSum)
	if sum.GreaterThan(balance) {
		return entity.ErrInsufficientBalance
	}

	withdrawal, err := entity.NewWithdrawal(orderNumber, sum)
	if err != nil {
		return err
	}
	return s.repo.Create(ctx, *withdrawal)
}

func (s *Service) GetWithdrawals(ctx context.Context, userID entity.ID) ([]*entity.Withdrawal, error) {
	userLock := s.getUserMutex(userID)
	userLock.RLock()
	defer userLock.RUnlock()

	return s.repo.List(ctx, userID)
}

func (s *Service) GetWithdrawalsSum(ctx context.Context, userID entity.ID) (decimal.NullDecimal, error) {
	userLock := s.getUserMutex(userID)
	userLock.RLock()
	defer userLock.RUnlock()

	return s.repo.GetSum(ctx, userID)
}

func (s *Service) getUserMutex(userID entity.ID) *sync.RWMutex {
	mx, _ := s.userMX.LoadOrStore(userID, &sync.RWMutex{})
	return mx.(*sync.RWMutex)
}
