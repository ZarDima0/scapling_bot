package webSocketBibit

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

const (
	TestBibitWSUrlSpot = "wss://stream.bybit.com/v5/public/spot"
)

type WebSocketClient struct {
	conn             *websocket.Conn
	orderBookChannel chan OrderBookData
	tickerChannel    chan TickerData
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

type TickerData struct {
	Topic string `json:"topic"`
	Type  string `json:"type"`
	Ts    int64  `json:"ts"`
	Data  struct {
		Symbol       string `json:"symbol"`
		LastPrice    string `json:"lastPrice"`
		HighPrice24h string `json:"highPrice24h"`
		LowPrice24h  string `json:"lowPrice24h"`
		Volume24h    string `json:"volume24h"`
		Turnover24h  string `json:"turnover24h"`
	} `json:"data"`
}

func (c *WebSocketClient) OrderBookChannel() chan OrderBookData {
	return c.orderBookChannel
}

func (c *WebSocketClient) TickerChannel() chan TickerData {
	return c.tickerChannel
}

func StartWebSocketClient(symbol string) *WebSocketClient {
	conn, _, err := websocket.DefaultDialer.Dial(TestBibitWSUrlSpot, nil)
	if err != nil {
		log.Fatalf("Error connecting to WebSocket: %v", err)
	}
	orderBookChannel := make(chan OrderBookData, 10)
	tickerChannel := make(chan TickerData, 10)

	subOrderBook := map[string]interface{}{"op": "subscribe", "args": []string{"orderbook.1." + symbol}}
	subOrderBookJSON, err := json.Marshal(subOrderBook)
	if err != nil {
		log.Fatalf("Error marshaling order book subscription: %v", err)
	}
	err = conn.WriteMessage(websocket.TextMessage, subOrderBookJSON)
	if err != nil {
		log.Fatalf("Error subscribing to order book: %v", err)
	}

	subTicker := map[string]interface{}{"op": "subscribe", "args": []string{"tickers." + symbol}}
	subTickerJSON, err := json.Marshal(subTicker)
	if err != nil {
		log.Fatalf("Error marshaling ticker subscription: %v", err)
	}
	err = conn.WriteMessage(websocket.TextMessage, subTickerJSON)
	if err != nil {
		log.Fatalf("Error subscribing to ticker: %v", err)
	}

	client := WebSocketClient{
		conn:             conn,
		orderBookChannel: orderBookChannel,
		tickerChannel:    tickerChannel,
	}

	go listenWebSocketMessages(&client)

	return &client
}

func listenWebSocketMessages(c *WebSocketClient) {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("Failed to read message: %v", err)
			close(c.orderBookChannel)
			close(c.tickerChannel)
			return
		}
		var tickerData TickerData
		if err := json.Unmarshal(message, &tickerData); err == nil {
			c.tickerChannel <- tickerData
			continue
		}

		log.Printf("Received an unknown message: %v", string(message))
	}
}
