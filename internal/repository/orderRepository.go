package repository

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"
)

type Order struct {
	ExternalOrderId string
	Status          string
	Price           string
	Quantity        string
	Side            string
	RawRequest      json.RawMessage
	RawAnswer       json.RawMessage
	CreatedAt       time.Time
}

func CreateOrder(db *sql.DB, order Order) {
	query := `
		INSERT INTO orders (external_order_id ,status, price, quantity, side, raw_request, raw_answer, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	err, _ := db.Exec(
		query,
		order.ExternalOrderId,
		order.Status,
		order.Price,
		order.Quantity,
		order.Side,
		order.RawRequest,
		order.RawAnswer,
		order.CreatedAt,
	)
	if err != nil {
		log.Fatalf("Error while inserting order: %v", err)
	}
}
