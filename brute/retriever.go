package brute

import (
	"encoding/json"

	"go.uber.org/zap"
)

// StartRetriever launches a goroutine to handle message retrieval from a RabbitMQ queue.
func (b *Brute) StartRetriever() {
	// Start the retrieve function in a new goroutine to enable asynchronous processing.
	go b.retrieve()
}

// retrieve listens for messages on the RabbitMQ queue, processes each message, and updates the brute force operation's state accordingly.
func (b *Brute) retrieve() error {
	b.LogDebug("Starting retriever")

	// Begin consuming messages from the specified queue.
	msgs, err := b.rabbit.Channel.Consume(
		b.rabbit.RetrieveQueue.Name, // The name of the queue to consume from.
		"",                          // Consumer tag (identifier).
		false,                       // Auto-acknowledge messages upon receipt.
		false,                       // Exclusive consumer access.
		false,                       // No local - do not receive messages published by this connection.
		false,                       // No wait - do not wait for server acknowledgement.
		nil,                         // Additional arguments.
	)
	if err != nil {
		// Log and return the error if unable to start consuming from the queue.
		b.LogError("Error while consuming queue", err)
		return err
	}

	// Loop indefinitely to process messages as they arrive.
	for d := range msgs {
		resp := response{} // Initialize a variable to hold the unmarshalled response.

		// Attempt to unmarshal the JSON-encoded message body into a response struct.
		err := json.Unmarshal(d.Body, &resp)
		if err != nil {
			// Log the error and skip processing this message if unmarshalling fails.
			b.LogError("Error while unmarshalling response", err)
			continue
		}

		// Log the received response for informational purposes.
		b.LogInfo("Response received", zap.String("item", resp.Item), zap.Int("status", resp.Status))

		// Process the response based on its status.
		switch resp.Status {
		case resultSuccess:
			// For successful responses, enqueue the item and increment the total counter.
			b.stats.hits.Enqueue(resp.Item)
			b.IncrementTotal()
		case resultInvalid:
			// For invalid responses, enqueue the item into the bad queue and increment the total counter.
			b.stats.bad.Enqueue(resp.Item)
			b.IncrementTotal()
		default:
			// For all other responses, record the errors.
			b.AddError(resp.Errors)
		}
	}
	return nil // Optionally, return nil to indicate the function completes without an error.
}
