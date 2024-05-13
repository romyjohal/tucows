package handlers

import (
	"github.com/jackc/pgx/v4"
	amqp "github.com/rabbitmq/amqp091-go"
)

type BaseHandler struct {
	db       *pgx.Conn
	msgQueue *amqp.Channel
}

// NewBaseHandler returns a new BaseHandler
func NewBaseHandler(db *pgx.Conn, q *amqp.Channel) *BaseHandler {
	return &BaseHandler{
		db:       db,
		msgQueue: q,
	}
}
