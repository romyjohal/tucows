package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

type paymentMsg struct {
	Order_id      string  `json:"order_id"`
	PaymentAmount float64 `json:"payment_amount"`
	PaymentInfo   string  `json:"payment_info"`
}

func processPayment(cardnumber string) bool {
	return cardnumber == "555"
}

func updateOrder(db *pgx.Conn, status bool, orderID string) error {
	statusText := "FAILED"
	if status {
		statusText = "SUCESS"
	}
	updateSQL := fmt.Sprintf(`
	UPDATE item_order
	SET payment_status = '%s'
	WHERE id = $1
	`, statusText)
	rows, err := db.Query(context.Background(), updateSQL, orderID)
	rows.Close()
	return err
}
