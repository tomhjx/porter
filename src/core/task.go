package core

type Task struct {
	ID    int64
	State int8
	Order Order
}

const (
	TaskReadyState      int8 = 1
	TaskProcessingState int8 = 2
	TaskSucceedState    int8 = 3
	TaskFailState       int8 = 4
)
