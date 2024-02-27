package brute

import (
	"sync/atomic"
)

// AddError records one or more errors encountered during the brute force operation.
// This method is safe to be called from multiple goroutines concurrently.
func (b *Brute) AddError(value []string) {
	b.stats.mx.Lock() // Ensure exclusive access to the errors slice
	defer b.stats.mx.Unlock()
	b.stats.errors = append(b.stats.errors, value...)
}

// GetErrors returns a copy of all recorded errors.
// This method provides a snapshot of errors, ensuring the original slice is not modified.
func (b *Brute) GetErrors() []string {
	b.stats.mx.Lock() // Ensure exclusive access to the errors slice
	defer b.stats.mx.Unlock()

	// Create a copy of the errors slice to prevent external modifications
	temp := make([]string, len(b.stats.errors))
	copy(temp, b.stats.errors)
	return temp
}

// IncrementTotal atomically increments the total count of operations or attempts.
// This method can be called safely from multiple goroutines without using mutexes.
func (b *Brute) IncrementTotal() {
	atomic.AddInt32(&b.stats.total, 1)
}

// GetTotal returns the current total count of operations or attempts.
// Since the total is managed atomically, this method is safe to call concurrently
// without the need for locking.
func (b *Brute) GetTotal() int32 {
	return atomic.LoadInt32(&b.stats.total)
}
