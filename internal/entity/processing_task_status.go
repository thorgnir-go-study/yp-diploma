package entity

type ProcessingTaskStatus int

const (
	ProcessingTaskStatusScheduled ProcessingTaskStatus = iota
	ProcessingTaskStatusProcessing
	ProcessingTaskStatusProcessed
)

var processingTaskStatusStringValues = [...]string{"SCHEDULED", "PROCESSING", "PROCESSED"}

func (p ProcessingTaskStatus) String() string {
	return processingTaskStatusStringValues[p]
}
