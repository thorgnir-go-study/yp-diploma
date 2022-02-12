package entity

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
	// TODO: вставить проверку по Луну
	return nil
}
