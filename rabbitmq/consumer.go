package rabbitmq

import (
	"context"
	"fmt"
	"instant-messenger-backend/controllers"
)

func ConsumeMessages(ctx context.Context) error {
	ch, err := rabbitMQ.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Declare the variables to receive the channel and error
	msgs, err := ch.Consume(
		"DatabaseQueue",    // Queue name
		"DatabaseConsumer", // Consumer tag (for identification)
		true,               // Auto-ack (acknowledges messages automatically)
		false,              // Exclusive (only accessible to this consumer)
		false,              // No-local (don't receive messages published by itself)
		false,              // No-wait (don't wait for a server response)
		nil,                // Arguments (optional)
	)
	if err != nil {
		return err
	}

	fmt.Printf("DatabaseConsumer successfully Running!\n")

	// Process messages in a loop
	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				// Channel closed
				break
			}

			// Send message to database
			_, err := controllers.SaveMessageToDatabase(msg.Body)
			if err != nil {
				fmt.Printf("Error sending message to database: %v\n", err)
				// Handle the error appropriately
			}

		case <-ctx.Done():
			// Context canceled, stop consuming
			return nil
		}
	}
}
