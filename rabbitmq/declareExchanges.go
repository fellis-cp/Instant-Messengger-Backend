package rabbitmq

import "fmt"

// DeclareExchange creates an exchange on the RabbitMQ server.
func DeclareExchange(exchangeName string, exchangeType string) error {
	ch, err := rabbitMQ.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName, // Name of the exchange
		exchangeType, // Type of the exchange (e.g., "fanout", "direct", "topic", etc.)
		true,         // Durable (survives server restarts)
		false,        // Auto-delete (deleted when no longer in use)
		false,        // Internal (not visible to other connections)
		false,        // No-wait (do not wait for a response)
		nil,          // Arguments (optional)
	)
	if err != nil {
		return err
	}
	fmt.Printf("Exchange '%s' declared successfully!\n", exchangeName)
	return nil
}
