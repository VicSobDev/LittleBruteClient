package queue

import "sync"

// Queue represents a thread-safe FIFO (First In, First Out) queue.
type Queue struct {
	items                             []interface{}         // Slice to hold the queue's items.
	rwmutex                           sync.RWMutex          // Read/Write mutex to protect concurrent access to the queue.
	newItemAvailableChan              chan chan interface{} // Channel of channels to notify waiting consumers of new items.
	signalNewItemOrUnlockConsumerChan chan struct{}         // Channel used to signal that a new item is available or to unlock a waiting consumer.
}

// NewQueue initializes and returns a new Queue instance with default settings.
func NewQueue() *Queue {
	return &Queue{
		items:                             make([]interface{}, 0),
		newItemAvailableChan:              make(chan chan interface{}, WaitForNextElementChanCapacity),
		signalNewItemOrUnlockConsumerChan: make(chan struct{}, WaitForNextElementChanCapacity),
	}
}
