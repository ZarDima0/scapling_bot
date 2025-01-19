package main

import (
	"fmt"
	"github.com/ZarDima0/scapling_bot/internal/config"
	"github.com/ZarDima0/scapling_bot/internal/webSocketBibit"
	"github.com/hirokisan/bybit/v2"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	coins := []bybit.Coin{
		bybit.CoinUSDT,
		bybit.Coin("TON"),
	}
	client := bybit.NewClient().WithAuth(cfg.TestKey, cfg.TestSecret).WithBaseURL(cfg.TestApiBybit)
	wallet, err := client.V5().Account().GetWalletBalance(bybit.AccountTypeV5UNIFIED, coins)
	websocket := webSocketBibit.StartWebSocketClient("TONUSDT")
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	for ticket := range websocket.TickerChannel() {
		fmt.Println(ticket.Data.LastPrice)
	}

	if err != nil {
		log.Fatal("Error getting wallet balance: ", err)
	}
	fmt.Printf("Wallet Balance: %+v\n", wallet.RetCode)
}
