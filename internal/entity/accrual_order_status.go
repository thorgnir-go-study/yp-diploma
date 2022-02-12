package entity

import "errors"

type AccrualOrderStatus string

const (
	AccrualOrderStatusRegistered AccrualOrderStatus = "REGISTERED"
	AccrualOrderStatusInvalid    AccrualOrderStatus = "INVALID"
	AccrualOrderStatusProcessing AccrualOrderStatus = "PROCESSING"
	AccrualOrderStatussProcessed AccrualOrderStatus = "PROCESSED"
)

func StringToAccrualOrderStatus(raw string) (AccrualOrderStatus, error) {
	s := AccrualOrderStatus(raw)
	err := s.Validate()
	if err != nil {
		return "", err
	}
	return s, nil
}

func (s AccrualOrderStatus) Validate() error {
	switch s {
	case AccrualOrderStatusRegistered, AccrualOrderStatusInvalid, AccrualOrderStatusProcessing, AccrualOrderStatussProcessed:
		return nil
	}
	return errors.New("invalid accrual order status")
}
