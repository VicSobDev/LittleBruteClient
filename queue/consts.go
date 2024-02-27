package queue

import "errors"

const (
	// Capacity for the channel used to wait for the next element
	WaitForNextElementChanCapacity = 1000
)

var (
	ErrQueueIsEmpty = errors.New("queue is empty")
)
