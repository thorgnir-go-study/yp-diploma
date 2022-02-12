package entity

import "errors"

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

func StringToOrderStatus(raw string) (OrderStatus, error) {
	s := OrderStatus(raw)
	err := s.Validate()
	if err != nil {
		return "", err
	}
	return s, nil
}

func (s OrderStatus) Validate() error {
	switch s {
	case OrderStatusNew, OrderStatusProcessing, OrderStatusInvalid, OrderStatusProcessed:
		return nil
	}
	return errors.New("invalid order status")
}
