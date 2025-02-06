package repository

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"
)

type Ticket struct {
	ID          int
	Symbol      string
	LastPrice   string
	HighPrice24 string
	LowPrice24  string
	Volume24    string
	Turnover24  string
	RawTicket   json.RawMessage
	CreatedAt   time.Time
}

func CreateTicket(db *sql.DB, ticket Ticket) {
	query := `
		INSERT INTO tickets (symbol, last_price, high_price_24, low_price_24, volume_24, turnover_24, raw_ticket, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := db.Exec(
		query,
		ticket.Symbol,
		ticket.LastPrice,
		ticket.HighPrice24,
		ticket.LowPrice24,
		ticket.Volume24,
		ticket.Turnover24,
		ticket.RawTicket,
		ticket.CreatedAt,
	)

	if err != nil {
		log.Fatalf("Error creating ticket: %s", err)
	}
}
