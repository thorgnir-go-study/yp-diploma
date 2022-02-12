package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

type Withdrawal struct {
	ID          ID
	OrderNumber OrderNumber
	Sum         decimal.Decimal
	ProcessedAt time.Time
}

func NewWithdrawal(orderNumber OrderNumber, sum decimal.Decimal) (*Withdrawal, error) {
	id, err := NewID()
	if err != nil {
		return nil, err
	}
	w := &Withdrawal{
		ID:          id,
		OrderNumber: orderNumber,
		Sum:         sum,
		ProcessedAt: time.Now(),
	}
	return w, nil
}
