package brute

import (
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

// Define result codes as constants for clarity and maintainability.
const (
	resultError   = 0 // Indicates an error occurred during processing.
	resultSuccess = 1 // Indicates successful processing.
	resultInvalid = 2 // Indicates an invalid result or input.
)

// Publish sends an item to the RabbitMQ publish queue associated with the Brute instance.
func (b *Brute) Publish(item string) error {
	// Ignore empty items to avoid publishing unnecessary messages.
	if item == "" {
		return nil
	}

	// Log the item being queued if verbose logging is enabled.
	if b.Verbose {
		b.logger.Info("Queueing it up", zap.String("item", item))
	}

	// Attempt to publish the item to the designated RabbitMQ queue.
	err := b.rabbit.Channel.Publish(
		"",                         // Exchange: Default exchange is used, routing directly to queues.
		b.rabbit.PublishQueue.Name, // Routing key: The name of the queue to publish to.
		false,                      // Mandatory: Do not return an error if no queues match the routing key.
		false,                      // Immediate: Do not return an error if no consumers are available.
		amqp.Publishing{
			ContentType: "text/plain", // Set the content type to text/plain.
			Body:        []byte(item), // Convert the item to a byte slice for the message body.
		})

	// Return any error encountered during the publish operation.
	return err
}
