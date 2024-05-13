package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := pgx.Connect(context.Background(), "postgres://user:user@0.0.0.0:5432/user")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Println("connected to database")
	}
	defer conn.Close(context.Background())

	msgQueueConn, err := amqp.Dial("amqp://guest:guest@0.0.0.0:5672/")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to rabbitMQ: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Println("connected to msg queue")
	}
	defer msgQueueConn.Close()

	channel, err := msgQueueConn.Channel()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open channel on rabbitMQ: %v\n", err)
		os.Exit(1)
	}
	defer channel.Close()

	msgs, err := channel.Consume(
		"payments",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create channel on rabbitMQ: %v\n", err)
		os.Exit(1)
	}

	_, err = channel.QueueDeclare(
		"payments",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create queue on rabbitMQ: %v\n", err)
		os.Exit(1)
	}

	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			var p paymentMsg
			json.Unmarshal(msg.Body, &p)
			status := processPayment(p.PaymentInfo)
			err := updateOrder(conn, status, p.Order_id)
			if err != nil {
				fmt.Printf("Error updating order_id: <%s> Error:<%s>", p.Order_id, err.Error())
				continue
			}
			fmt.Printf("Processed Payment for order_id: <%s>\n", p.Order_id)
		}
	}()

	fmt.Println("Payment serivce Running")
	<-forever
}
