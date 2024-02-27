package queue

import "context"

// Helper functions below are used to abstract and encapsulate specific behaviors for clarity and maintainability.

// signalConsumerIfWaiting signals a waiting consumer that a new item is available.
func (q *Queue) signalConsumerIfWaiting() {
	select {
	case q.signalNewItemOrUnlockConsumerChan <- struct{}{}:
	default:
	}
}

// trySendToWaitingConsumer attempts to send an item directly to a waiting consumer.
func (q *Queue) trySendToWaitingConsumer(item interface{}) bool {
	select {
	case listener := <-q.newItemAvailableChan: // Try to get a waiting consumer's channel.
		select {
		case listener <- item: // Send item to the waiting consumer.
			return true
		default:
		}
	default:
	}
	return false // Return false if no waiting consumer was found.
}

// lockAndEnqueueItem locks the queue and enqueues the item.
func (q *Queue) lockAndEnqueueItem(item interface{}) {
	q.rwmutex.Lock()                // Lock the queue for safe access.
	q.items = append(q.items, item) // Add the item to the queue.
	q.rwmutex.Unlock()              // Unlock the queue.
}

// dequeueIfAvailable locks the queue and attempts to dequeue an item if available.
func (q *Queue) dequeueIfAvailable() (interface{}, error) {
	q.rwmutex.Lock() // Lock the queue for safe access.
	defer q.rwmutex.Unlock()

	if len(q.items) == 0 { // Check if the queue is empty.
		return nil, ErrQueueIsEmpty
	}

	item := q.items[0]    // Retrieve the first item.
	q.items = q.items[1:] // Remove the first item from the queue.
	return item, nil
}

// waitForNewItemOrContextCancellation waits for a new item to be enqueued or for the context to be cancelled.
func (q *Queue) waitForNewItemOrContextCancellation(ctx context.Context) (interface{}, error) {
	waitChan := make(chan interface{}) // Create a channel to wait for a new item.
	select {
	case q.newItemAvailableChan <- waitChan: // Register the wait channel for a new item.
		return q.handleWaitForNewItemOrContextCancellation(ctx, waitChan)
	default:
		return nil, ErrQueueIsEmpty // Return error if the queue is currently empty and no waiting is set up.
	}
}

// handleWaitForNewItemOrContextCancellation handles the logic of waiting for a new item or context cancellation.
func (q *Queue) handleWaitForNewItemOrContextCancellation(ctx context.Context, waitChan chan interface{}) (interface{}, error) {
	select {
	case item := <-waitChan: // Wait for a new item to be sent.
		return item, nil
	case <-ctx.Done(): // Return if the context is cancelled.
		return nil, ctx.Err()
	case <-q.signalNewItemOrUnlockConsumerChan: // Check again if a new item is available after being signaled.
		return q.dequeueIfAvailable()
	}
}
