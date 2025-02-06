package main

import (
	"database/sql"
	"encoding/json"
	"github.com/ZarDima0/scapling_bot/internal/config"
	"github.com/ZarDima0/scapling_bot/internal/database"
	"github.com/ZarDima0/scapling_bot/internal/repository"
	"github.com/ZarDima0/scapling_bot/internal/webSocketBibit"
	"github.com/hirokisan/bybit/v2"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	cfg := config.MustLoad()
	//coins := []bybit.Coin{
	//	bybit.CoinUSDT,
	//	bybit.Coin("TON"),
	//}
	//wallet, err := client.V5().Account().GetWalletBalance(bybit.AccountTypeV5UNIFIED, coins)
	//if err != nil {
	//	log.Fatal("Error getting wallet balance: ", err)
	//}
	db := database.NewDB(cfg)

	client := bybit.NewClient().WithAuth(cfg.TestKey, cfg.TestSecret).WithBaseURL(cfg.TestApiBybit)
	err := client.SyncServerTime()
	if err != nil {
		log.Fatalf("Error syncing server time: %v", err)
	}
	websocket := webSocketBibit.StartWebSocketClient("TONUSDT")
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	for ticket := range websocket.TickerChannel() {

		if ticket.Data.LastPrice == "" {
			continue
		}

		orderParam := bybit.V5CreateOrderParam{
			Category:  bybit.CategoryV5Spot,
			Symbol:    "TONUSDT",
			Side:      bybit.SideBuy,
			OrderType: bybit.OrderTypeMarket,
			Qty:       "1.0",
		}
		orderBybit, err := client.V5().Order().CreateOrder(orderParam)
		if err != nil {
			log.Fatalf("Error creating order: %s", err)
		}

		rawRequest, err := json.Marshal(orderParam)
		if err != nil {
			log.Fatalf("raw request error : %s", err)
		}
		rawResponse, err := json.Marshal(orderBybit)
		if err != nil {
			log.Fatalf("raw response error : %v", err)
		}
		order := repository.Order{
			ExternalOrderId: orderBybit.Result.OrderID,
			Status:          orderBybit.RetMsg,
			Price:           ticket.Data.LastPrice,
			Quantity:        "1",
			Side:            string(bybit.SideBuy),
			RawRequest:      rawRequest,
			RawAnswer:       rawResponse,
		}
		repository.CreateOrder(db, order)
		makeTicket(db, &ticket)
	}
}

func makeTicket(db *sql.DB, ticket *webSocketBibit.TickerData) {
	jsonData, err := json.Marshal(ticket)
	if err != nil {
		log.Fatal("Error serializing ticket:", err)
	}

	repository.CreateTicket(db, repository.Ticket{
		Symbol:      ticket.Topic,
		LastPrice:   ticket.Data.LastPrice,
		HighPrice24: ticket.Data.HighPrice24h,
		LowPrice24:  ticket.Data.LowPrice24h,
		Volume24:    ticket.Data.Volume24h,
		Turnover24:  ticket.Data.Turnover24h,
		RawTicket:   jsonData,
		CreatedAt:   time.Now(),
	})
}

func calculatePriceWithPercentage(priceStr string, percentage float64) (float64, error) {
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		panic("Error parsing price")
	}

	return price * (1 + percentage/100), nil
}

//func processOrder(bybitClient bybit.Client, db *sql.DB, ticket *webSocketBibit.TickerData) {
//	orderParam := bybit.V5CreateOrderParam{
//		Category:  bybit.CategoryV5Spot,
//		Symbol:    "TONUSDT",
//		Side:      bybit.SideBuy,
//		OrderType: bybit.OrderTypeMarket,
//		Qty:       "1",
//	}
//	orderBybit, err := bybitClient.V5().Order().CreateOrder(orderParam)
//	if err != nil {
//		log.Fatalf("Error creating order: %s", err)
//	}
//
//	rawRequest, err := json.Marshal(orderParam)
//	if err != nil {
//		log.Fatalf("raw request error : %s", err)
//	}
//	rawResponse, err := json.Marshal(orderBybit)
//	if err != nil {
//		log.Fatalf("raw response error : %v", err)
//	}
//	order := repository.Order{
//		ExternalOrderId: orderBybit.Result.OrderID,
//		Status:          orderBybit.RetMsg,
//		Price:           ticket.Data.LastPrice,
//		Quantity:        "1",
//		Side:            string(bybit.SideBuy),
//		RawRequest:      rawRequest,
//		RawAnswer:       rawResponse,
//	}
//	repository.CreateOrder(db, order)
//}
