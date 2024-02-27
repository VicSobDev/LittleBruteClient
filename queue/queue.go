package queue

import (
	"context"
)

// Enqueue adds an item to the queue, potentially unlocking a waiting consumer.
func (q *Queue) Enqueue(item interface{}) {
	q.signalConsumerIfWaiting()            // Signal waiting consumer if any.
	if !q.trySendToWaitingConsumer(item) { // Attempt to send item directly to a waiting consumer.
		q.lockAndEnqueueItem(item) // If no consumer is waiting, lock the queue and add the item.
	}
}

// Dequeue removes and returns the first item from the queue if available.
func (q *Queue) Dequeue() (interface{}, error) {
	q.rwmutex.Lock() // Lock the queue for safe access.
	defer q.rwmutex.Unlock()

	if len(q.items) == 0 { // Check if the queue is empty.
		return nil, ErrQueueIsEmpty
	}

	item := q.items[0]    // Retrieve the first item.
	q.items = q.items[1:] // Remove the first item from the queue.
	return item, nil
}

// DequeueOrWaitForNextElementContext attempts to dequeue an item or waits for an item to be enqueued, respecting the given context.
func (q *Queue) DequeueOrWaitForNextElementContext(ctx context.Context) (interface{}, error) {
	for {
		// Try to dequeue if items are available.
		if item, err := q.dequeueIfAvailable(); err == nil {
			return item, nil
		}

		// Wait for a new item to be enqueued or for the context to be cancelled.
		if item, err := q.waitForNewItemOrContextCancellation(ctx); err == nil {
			return item, nil
		} else if err != ErrQueueIsEmpty {
			return nil, err // Return error if it's not the empty queue error.
		}
	}
}

// GetLength returns the current number of items in the queue.
func (q *Queue) GetLength() int {
	q.rwmutex.RLock() // Acquire read lock for safe access.
	defer q.rwmutex.RUnlock()

	return len(q.items)
}
