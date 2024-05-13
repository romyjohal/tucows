package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type item struct {
	Id    int     `json:"id"`
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}

func (b *BaseHandler) List() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		query := `
		SELECT
			ID,
			NAME,
			PRICE
		FROM
			item
		;`

		rows, err := b.db.Query(ctx, query)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}

		defer rows.Close()

		items := make([]item, 0)
		for rows.Next() {
			i := item{}
			err := rows.Scan(
				&i.Id,
				&i.Name,
				&i.Price,
			)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				w.WriteHeader(http.StatusInternalServerError)
			}
			items = append(items, i)
		}

		itemsJSON, err := json.Marshal(items)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Write(itemsJSON)
	}
}
