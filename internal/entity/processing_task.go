package entity

import "time"

type ProcessingTask struct {
	ID          ID
	OrderID     ID
	OrderNumber OrderNumber
	ToRunAt     time.Time
	Status      ProcessingTaskStatus
	UpdatedAt   time.Time
}

func NewProcessingTask(orderID ID, orderNumber OrderNumber, runAt time.Time) (*ProcessingTask, error) {
	id, err := NewID()
	if err != nil {
		return nil, err
	}
	task := &ProcessingTask{
		ID:          id,
		OrderID:     orderID,
		OrderNumber: orderNumber,
		ToRunAt:     runAt,
		Status:      ProcessingTaskStatusScheduled,
		UpdatedAt:   time.Now(),
	}

	return task, nil
}
