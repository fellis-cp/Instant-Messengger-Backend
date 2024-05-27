package rabbitmq

import "fmt"

// DeclareQueue creates a queue on the RabbitMQ server.
func DeclareQueue(queueName string) error {
	ch, err := rabbitMQ.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // Name of the queue
		true,      // Durable (survives server restarts)
		false,     // Auto-delete (deleted when no longer in use)
		false,     // Exclusive (only accessible to the current connection)
		false,     // No-wait (don't wait for a server response)
		nil,       // Arguments (optional)
	)
	if err != nil {
		return err
	}
	fmt.Printf("Queue '%s' declared successfully!\n", q.Name)
	return nil
}
