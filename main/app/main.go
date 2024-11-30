package main

import (
	"fmt"
	"github.com/ZarDima0/scapling_bot/internal/websocketBibit"
	"log"
)

func main() {
	client, err := websocketBibit.NewClient("wss://stream.bybit.com/v5/public/spot")
	if err != nil {
		log.Fatalf("Failed to create WebSocket client: %v", err)
	}

	err = client.SubscribeToTicker("BTCUSDT")
	if err != nil {
		log.Fatalf("Failed to subscribe to ticker: %v", err)
	}

	orderBookChannel := client.GetOrderBookChannel()

	for orderBookData := range orderBookChannel {
		fmt.Println("Received order book data:")
		fmt.Println("Bids:", orderBookData.Data.Bids)
		fmt.Println("Asks:", orderBookData.Data.Asks)
	}

	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("Failed to close WebSocket: %v", err)
		}
	}()
}
