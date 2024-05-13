package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type orderStatusReq struct {
	Customer_id string `json:"customer_id"`
}

type orderStatus struct {
	ID            int       `json:"id"`
	CustomerID    int       `json:"customer_id"`
	ItemID        int       `json:"item_id"`
	Quantity      int       `json:"quantity"`
	UnitPrice     float32   `json:"unit_price"`
	OrderDate     time.Time `json:"order_date"`
	PaymentStatus string    `json:"payment_status"`
}

func (b *BaseHandler) OrderStatus() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		body, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}

		defer req.Body.Close()

		var o orderStatusReq
		err = json.Unmarshal(body, &o)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}

		query := `
		SELECT
			id,
			customer_id,
			item_id,
			quantity,
			unit_price,
			order_date,
			payment_status
		FROM
			item_order
		WHERE
			item_order.customer_id = $1
		;`

		rows, err := b.db.Query(ctx, query, o.Customer_id)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}

		defer rows.Close()

		orderStatuses := make([]orderStatus, 0)
		for rows.Next() {
			o := orderStatus{}
			err := rows.Scan(
				&o.ID,
				&o.CustomerID,
				&o.ItemID,
				&o.Quantity,
				&o.UnitPrice,
				&o.OrderDate,
				&o.PaymentStatus,
			)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				w.WriteHeader(http.StatusInternalServerError)
			}
			orderStatuses = append(orderStatuses, o)
		}

		orderStatusesJSON, err := json.Marshal(orderStatuses)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Write(orderStatusesJSON)
	}
}
