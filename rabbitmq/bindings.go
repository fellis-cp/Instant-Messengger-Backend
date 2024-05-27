package rabbitmq

import (
	"fmt"
)

func BindQueueToExchange(queueName string, exchangeName string, routingKey string) error {
	ch, err := rabbitMQ.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.QueueBind(
			queueName,
			routingKey,
			exchangeName,
			false,
			nil,
	)
	if err != nil {
			return err
	}
	fmt.Printf("Queue '%s' bound to exchange '%s' with routing key '%s'\n", queueName, exchangeName, routingKey)
	return nil
}