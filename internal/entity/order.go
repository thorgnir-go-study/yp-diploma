package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

type Order struct {
	ID         ID
	UserID     ID
	Number     OrderNumber
	Status     OrderStatus
	Accrual    decimal.NullDecimal
	UploadedAt time.Time
	UpdatedAt  time.Time
}

func NewOrder(userID ID, orderNumber OrderNumber) (*Order, error) {

	id, err := NewID()
	if err != nil {
		return nil, err
	}
	o := &Order{
		ID:         id,
		UserID:     userID,
		Number:     orderNumber,
		Status:     OrderStatusNew,
		UploadedAt: time.Now(),
		UpdatedAt:  time.Now(),
	}
	err = o.Validate()
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (o *Order) Validate() error {
	if o.ID.IsNil() || o.Number.Validate() != nil || o.Status.Validate() != nil {
		return ErrInvalidEntity
	}
	return nil
}
