package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

type Withdrawal struct {
	ID          ID
	UserID      ID
	OrderNumber OrderNumber
	Sum         decimal.Decimal
	ProcessedAt time.Time
}

func NewWithdrawal(userID ID, orderNumber OrderNumber, sum decimal.Decimal) (*Withdrawal, error) {
	id, err := NewID()
	if err != nil {
		return nil, err
	}
	w := &Withdrawal{
		ID:          id,
		UserID:      userID,
		OrderNumber: orderNumber,
		Sum:         sum,
		ProcessedAt: time.Now(),
	}
	return w, nil
}
