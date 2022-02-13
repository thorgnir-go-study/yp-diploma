package entity

import "errors"

type OrderStatus int

const (
	OrderStatusUnknown OrderStatus = iota
	OrderStatusNew
	OrderStatusProcessing
	OrderStatusInvalid
	OrderStatusProcessed
)

var orderStatusStringValues = [...]string{"NEW", "PROCESSING", "INVALID", "PROCESSED"}

func StringToOrderStatus(raw string) (OrderStatus, error) {
	found := false
	var st OrderStatus
	for i := range orderStatusStringValues {
		if orderStatusStringValues[i] == raw {
			found = true
			st = OrderStatus(i + 1)
			break
		}
	}
	if !found {
		return OrderStatusUnknown, errors.New("invalid order status")
	}
	err := st.Validate()
	if err != nil {
		return OrderStatusUnknown, err
	}
	return st, nil
}

func (s OrderStatus) String() string {
	return orderStatusStringValues[s]
}

func (s OrderStatus) Validate() error {
	switch s {
	case OrderStatusNew, OrderStatusProcessing, OrderStatusInvalid, OrderStatusProcessed:
		return nil
	}
	return errors.New("invalid order status")
}
