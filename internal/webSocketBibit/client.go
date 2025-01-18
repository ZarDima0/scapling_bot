package webSocketBibit

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

const (
	TestBibitWSUrlSpot = "wss://stream-testnet.bybit.com/v5/public/spot"
)

type Client struct {
	conn             *websocket.Conn
	orderBookChannel chan OrderBookData
}

func (c *Client) OrderBookChannel() chan OrderBookData {
	return c.orderBookChannel
}

type OrderBookData struct {
	Topic string `json:"topic"`
	Type  string `json:"type"`
	Ts    int64  `json:"ts"`
	Data  struct {
		Symbol string     `json:"s"`
		Bids   [][]string `json:"b"`
		Asks   [][]string `json:"a"`
		U      int64      `json:"u"`
		Seq    int64      `json:"seq"`
		Cts    int64      `json:"cts"`
	} `json:"data"`
}

func StartOrderBookWebSocket(symbol string) *Client {
	conn, _, err := websocket.DefaultDialer.Dial(TestBibitWSUrlSpot, nil)

	if err != nil {
		log.Fatalf("Error connect WS: %s", err)
	}
	OrderBookChannel := make(chan OrderBookData)

	subscription := map[string]interface{}{
		"op":   "subscribe",
		"args": []string{"orderbook.1." + symbol},
	}
	subscriptionJSON, err := json.Marshal(subscription)
	if err != nil {
		log.Fatalf("Error json.Marshal %v", err)
	}
	err = conn.WriteMessage(websocket.TextMessage, subscriptionJSON)
	if err != nil {
		log.Fatalf("Error WriteMessage: %v", err)
	}
	client := Client{
		conn:             conn,
		orderBookChannel: OrderBookChannel,
	}
	go listenWebSocketOrderBook(&client)

	return &client
}

func listenWebSocketOrderBook(c *Client) {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("Failed to read message: %v", err)
			close(c.orderBookChannel)
			return
		}
		fmt.Println("Raw message received:", string(message))
		var orderBookData OrderBookData
		err = json.Unmarshal(message, &orderBookData)
		if err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		c.orderBookChannel <- orderBookData
	}
}
