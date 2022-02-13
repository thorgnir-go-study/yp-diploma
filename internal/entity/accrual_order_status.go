package entity

import "errors"

type AccrualOrderStatus int

const (
	AccrualOrderStatusUnknown AccrualOrderStatus = iota
	AccrualOrderStatusRegistered
	AccrualOrderStatusInvalid
	AccrualOrderStatusProcessing
	AccrualOrderStatusProcessed
)

var accrualOrderStatusStringValues = [...]string{"REGISTERED", "INVALID", "PROCESSING", "PROCESSED"}

func StringToAccrualOrderStatus(raw string) (AccrualOrderStatus, error) {
	found := false
	var st AccrualOrderStatus
	for i := range accrualOrderStatusStringValues {
		if accrualOrderStatusStringValues[i] == raw {
			found = true
			st = AccrualOrderStatus(i + 1)
			break
		}
	}
	if !found {
		return AccrualOrderStatusUnknown, errors.New("invalid accrual order status")
	}
	err := st.Validate()
	if err != nil {
		return AccrualOrderStatusUnknown, err
	}
	return st, nil
}

func (s *AccrualOrderStatus) String() string {
	if s == nil {
		return "nil"
	}
	return accrualOrderStatusStringValues[*s]
}

func (s AccrualOrderStatus) Validate() error {
	switch s {
	case AccrualOrderStatusRegistered, AccrualOrderStatusInvalid, AccrualOrderStatusProcessing, AccrualOrderStatusProcessed:
		return nil
	}
	return errors.New("invalid accrual order status")
}
