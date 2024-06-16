package rabbitmq

import (
	"context"
	"fmt"
	"instant-messenger-backend/controllers"
	"instant-messenger-backend/websocket"
	"log"
)

func ConsumeToDatabase(ctx context.Context) error {
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

			messageJson := msg.Body

			if success, err := controllers.SaveMessageToDatabase(messageJson); success {
				fmt.Println("Message saved to database successfully!")
				// After saving the message to the database, publish it to the client queue
				PublishToClient(messageJson)
			} else if err != nil {
				fmt.Printf("Error saving message to database: %v\n", err)
				// Handle the error appropriately
			}

		case <-ctx.Done():
			// Context canceled, stop consuming
			return nil
		}
	}
}

func ConsumeToClient(ctx context.Context) error {
	ch, err := rabbitMQ.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Declare the variables to receive the channel and error
	msgs, err := ch.Consume(
		"ClientQueue",    // Queue name
		"ClientConsumer", // Consumer tag (for identification)
		true,             // Auto-ack (acknowledges messages automatically)
		false,            // Exclusive (only accessible to this consumer)
		false,            // No-local (don't receive messages published by itself)
		false,            // No-wait (don't wait for a server response)
		nil,              // Arguments (optional)
	)
	if err != nil {
		return err
	}

	fmt.Printf("ClientConsumer successfully Running!\n")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			if success, err := websocket.ClientServer.SendChatToClient(d.Body); success {
				fmt.Printf("message sudah berhasil diterima di consumeToClient\n")
			} else {
				fmt.Printf("Error sending message to client: %v\n", err)
				// Handle the error appropriately
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

	return nil

}
