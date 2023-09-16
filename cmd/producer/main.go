package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/teatou/rabbitmq/internal"
)

func main() {
	conn, err := internal.ConnectRabbitMQ("xu", "senjoxu", "localhost:5672", "customers")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client, err := internal.NewRabbitMQClient(conn)
	if err != nil {
		panic("error creating rabbit client")
	}
	defer client.Close()

	consumeConn, err := internal.ConnectRabbitMQ("xu", "senjoxu", "localhost:5672", "customers")
	if err != nil {
		panic(err)
	}
	defer consumeConn.Close()

	consumeClient, err := internal.NewRabbitMQClient(consumeConn)
	if err != nil {
		panic("error creating rabbit client")
	}
	defer consumeClient.Close()

	queue, err := consumeClient.CreateQueue("", true, true)
	if err != nil {
		panic(err)
	}

	if err = consumeClient.CreateBinding(queue.Name, queue.Name, "customer_callbacks"); err != nil {
		panic(err)
	}

	messageBus, err := consumeClient.Consume(queue.Name, "customer-api", true)
	if err != nil {
		panic(err)
	}

	go func() {
		for message := range messageBus {
			log.Printf("Message callback %s\n", message.CorrelationId)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := 0; i < 10; i++ {
		if err = client.Send(ctx, "customer_events", "customers.created.us", amqp091.Publishing{
			ContentType:   "text/plain",
			DeliveryMode:  amqp091.Persistent,
			ReplyTo:       queue.Name,
			CorrelationId: fmt.Sprintf("customer_created_%d", i),
			Body:          []byte(`Message between services`),
		}); err != nil {
			panic(err)
		}
	}

	log.Println(client)
	var blockingChan chan struct{}
	<-blockingChan
}
