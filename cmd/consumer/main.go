package main

import (
	"context"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/teatou/rabbitmq/internal"
	"golang.org/x/sync/errgroup"
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

	publishConn, err := internal.ConnectRabbitMQ("xu", "senjoxu", "localhost:5672", "customers")
	if err != nil {
		panic(err)
	}
	defer publishConn.Close()

	publishClient, err := internal.NewRabbitMQClient(publishConn)
	if err != nil {
		panic("error creating rabbit client")
	}
	defer publishClient.Close()

	queue, err := client.CreateQueue("", true, true)
	if err != nil {
		panic(err)
	}

	if err := client.CreateBinding(queue.Name, "", "customer_events"); err != nil {
		panic(err)
	}

	messageBus, err := client.Consume(queue.Name, "email_service", false)
	if err != nil {
		panic(err)
	}

	var blockingChan chan struct{}

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	if err = client.ApplyQos(10, 0, true); err != nil {
		panic(err)
	}

	g.SetLimit(10)
	go func() {
		for message := range messageBus {
			msg := message
			g.Go(func() error {
				log.Printf("New message: %v\n", msg)
				time.Sleep(10 * time.Second)
				if err := msg.Ack(false); err != nil {
					log.Println("Ack message failed")
					return err
				}

				if err := publishClient.Send(ctx, "customer_callbacks", msg.ReplyTo, amqp091.Publishing{
					ContentType:   "text/plain",
					DeliveryMode:  amqp091.Persistent,
					Body:          []byte("Rpc complete"),
					CorrelationId: msg.CorrelationId,
				}); err != nil {
					panic(err)
				}

				log.Printf("Acknowleged message %s\n", msg.MessageId)
				return nil
			})
		}
	}()

	log.Println("Consuming, do CTRL+C")
	<-blockingChan
}
