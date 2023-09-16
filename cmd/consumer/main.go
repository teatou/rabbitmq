package main

import (
	"log"

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

	messageBus, err := client.Consume("customers_created", "email_service", false)
	if err != nil {
		panic(err)
	}

	var blockingChan chan struct{}

	go func() {
		for message := range messageBus {
			log.Printf("New message: %v\n", message)

			if !message.Redelivered {
				message.Nack(false, true)
				continue
			}

			if err := message.Ack(false); err != nil {
				log.Panicln("Ack message failed")
				continue
			}
			log.Printf("Ack message %s\n", message.MessageId)
		}
	}()

	log.Println("Consuming, do CTRL + C")
	<-blockingChan
}
