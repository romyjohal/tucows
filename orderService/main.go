package main

import (
	"context"
	"fmt"
	"net/http"
	"orderService/handlers"
	"os"

	"github.com/jackc/pgx/v4"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// connectecting to database - move to seperate function/package later
	conn, err := pgx.Connect(context.Background(), "postgres://user:user@0.0.0.0:5432/user")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Println("connected to database")
	}
	defer conn.Close(context.Background())

	// run database migration to set up tables
	migrations := make([]string, 0)
	migrations = append(migrations, `
	CREATE TABLE IF NOT EXISTS customer(
		id serial primary key, 
		name VARCHAR(40) NOT NULL
	)`)
	migrations = append(migrations, `
	CREATE TABLE IF NOT EXISTS item_order(
		id SERIAL primary key,
		customer_id INT NOT NULL,
		item_id INT NOT NULL,
		quantity INT NOT NULL,
		unit_price REAL NOT NULL,
		order_date DATE NOT NULL DEFAULT NOW(),
		payment_status VARCHAR(40) NOT NULL DEFAULT 'PENDING'
	)`)
	migrations = append(migrations, `
	CREATE TABLE IF NOT EXISTS item(
		id SERIAL NOT NULL,
		name VARCHAR(40) NOT NULL,
		price REAL NOT NULL
	)`)
	migrations = append(migrations, `
	INSERT INTO customer (id, name)
		VALUES(1, 'Greg Johnson')
		ON CONFLICT DO NOTHING
	`)
	migrations = append(migrations, `
	INSERT INTO item (id, name, price)
		VALUES(1, 'Bar of soap', 2.50)
		ON CONFLICT DO NOTHING
	`)

	for _, m := range migrations {
		_, err = conn.Exec(context.Background(), m)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}

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

	baseHandler := handlers.NewBaseHandler(conn, channel)

	http.HandleFunc("/list", baseHandler.List())
	http.HandleFunc("/order", baseHandler.Order())
	http.HandleFunc("/orderstatus", baseHandler.OrderStatus())

	fmt.Println("Order Service Running")
	http.ListenAndServe(":8090", nil)
}
