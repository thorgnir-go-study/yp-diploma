package entity

import "github.com/ShiraazMoollatjie/goluhn"

type OrderNumber string

func StringToOrderNumber(raw string) (OrderNumber, error) {
	var on = OrderNumber(raw)
	err := on.Validate()
	if err != nil {
		return "", err
	}
	return on, nil
}

func (n OrderNumber) Validate() error {
	err := goluhn.Validate(string(n))
	if err != nil {
		return ErrInvalidOrderNumber
	}
	return nil
}

func (n OrderNumber) String() string {
	return string(n)
}
