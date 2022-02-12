package entity

import "github.com/shopspring/decimal"

type AccrualOrder struct {
	OrderNumber OrderNumber
	Status      AccrualOrderStatus
	Accrual     *decimal.Decimal
}
