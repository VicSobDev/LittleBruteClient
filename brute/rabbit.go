package brute

import (
	"fmt"

	"github.com/streadway/amqp"
	"github.com/vicsobdev/LittleBruteClient/utils"
	"go.uber.org/zap"
)

// InitializeQueues sets up the necessary queues in RabbitMQ for the brute force operation.
// It configures the connection and channels based on the provided RabbitConfig.
func (b *Brute) InitializeQueues(config RabbitConfig) error {
	// Log the initiation process with debug level verbosity.
	b.LogDebug("Initializing queues")

	// Store the provided RabbitMQ configuration in the Brute instance.
	b.rabbit.Config = config

	// Establish a connection to the RabbitMQ server.
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", config.Username, config.Password, config.Host, config.Port))
	if err != nil {
		// Log and return the error if the connection fails.
		b.LogError("Error while connecting to RabbitMQ server", err)
		return fmt.Errorf("error while connecting to RabbitMQ server: %v", err)
	}
	// Store the connection in the Brute instance.
	b.rabbit.Conn = conn

	b.LogDebug("Connected to RabbitMQ server")

	// Open a channel on the RabbitMQ connection.
	ch, err := conn.Channel()
	if err != nil {
		// Log and return the error if opening the channel fails.
		b.LogError("Error while creating a channel", err)
		return fmt.Errorf("error while creating a channel: %v", err)
	}
	// Store the channel in the Brute instance.
	b.rabbit.Channel = ch

	b.LogDebug("Channel created")

	// Generate unique names for the queues using a UUID prefix.
	namePrefix := utils.GenUUID()
	retrieveQueueName := fmt.Sprintf("%s-retrieve", namePrefix)
	// Create the retrieve queue.
	retrieveQueue, err := createQueue(retrieveQueueName, ch)
	if err != nil {
		b.LogError("Error while creating a retrieve queue", err)
		return err
	}
	b.rabbit.RetrieveQueue = retrieveQueue

	// Similarly, create the publish queue.
	publishQueueName := fmt.Sprintf("%s-publish", namePrefix)
	publishQueue, err := createQueue(publishQueueName, ch)
	if err != nil {
		b.LogError("Error while creating a publish queue", err)
		return err
	}
	b.rabbit.PublishQueue = publishQueue

	// Log the successful initialization of queues with their names.
	b.LogInfo("Queues initialized", zap.String("retrieve", retrieveQueue.Name), zap.String("publish", publishQueue.Name))
	return nil
}

// createQueue declares a new queue on the given AMQP channel with specified name.
func createQueue(name string, ch *amqp.Channel) (amqp.Queue, error) {
	// Attempt to declare a new queue with the given name and default properties.
	q, err := ch.QueueDeclare(
		name,  // Queue name.
		false, // Durable: Persists across server restarts.
		false, // Auto-delete when unused.
		false, // Exclusive access.
		false, // No wait for server reply.
		nil,   // Additional arguments.
	)
	if err != nil {
		// Return an error if queue declaration fails.
		return amqp.Queue{}, fmt.Errorf("error while declaring a queue: %v", err)
	}

	// Return the declared queue if successful.
	return q, nil
}

// GetQueueName returns the name of the retrieve queue associated with this Brute instance.
func (b *Brute) GetQueueName() string {
	return b.rabbit.RetrieveQueue.Name
}
