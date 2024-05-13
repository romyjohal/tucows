package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
)

type orderReq struct {
	Customer_id string `json:"customer_id"`
	Item_id     string `json:"item_id"`
	Quantity    int    `json:"quantity"`
	PaymentInfo string `json:"payment_info"`
}

type paymentMsg struct {
	Order_id      string  `json:"order_id"`
	PaymentAmount float64 `json:"payment_amount"`
	PaymentInfo   string  `json:"payment_info"`
}

func (b *BaseHandler) Order() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		body, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer req.Body.Close()

		var o orderReq
		err = json.Unmarshal(body, &o)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		query := `
		SELECT
		  PRICE
		FROM
		  item
		WHERE
		  item.id = $1
		;`

		row := b.db.QueryRow(ctx, query, o.Item_id)
		price := 0.00
		err = row.Scan(&price)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}

		orderRecordInsert := `
		INSERT INTO item_order (customer_id, item_id, quantity, unit_price)
		  VALUES ($1, $2, $3, $4)
		  RETURNING id;
		`

		row = b.db.QueryRow(ctx, orderRecordInsert, o.Customer_id, o.Item_id, o.Quantity, price)
		var orderID int
		err = row.Scan(&orderID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		p := paymentMsg{
			Order_id:      strconv.Itoa(orderID),
			PaymentAmount: math.Round((price*float64(o.Quantity))*100) / 100,
			PaymentInfo:   o.PaymentInfo,
		}

		pJSON, err := json.Marshal(p)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = b.msgQueue.Publish(
			"",
			"payments",
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         pJSON,
			},
		)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write([]byte(strconv.Itoa(orderID)))
	}
}
