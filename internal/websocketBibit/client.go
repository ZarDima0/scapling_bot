package websocketBibit

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	conn             *websocket.Conn
	orderBookChannel chan OrderBookData
}

type OrderBookData struct {
	Topic string `json:"topic"`
	Type  string `json:"type"` // snapshot или delta
	Ts    int64  `json:"ts"`
	Data  struct {
		Symbol string     `json:"s"`
		Bids   [][]string `json:"b"`
		Asks   [][]string `json:"a"`
		U      int64      `json:"u"` // Update ID
		Seq    int64      `json:"seq"`
		Cts    int64      `json:"cts"`
	} `json:"data"`
}

func NewClient(url string) (*Client, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	orderBookChannel := make(chan OrderBookData)

	client := &Client{
		conn:             conn,
		orderBookChannel: orderBookChannel,
	}

	go client.listenWebSocket()

	return client, nil
}

func (c *Client) SubscribeToTicker(symbol string) error {
	subscription := map[string]interface{}{
		"op":   "subscribe",
		"args": []string{"orderbook.1." + symbol},
	}
	subscriptionJSON, err := json.Marshal(subscription)
	if err != nil {
		return err
	}
	err = c.conn.WriteMessage(websocket.TextMessage, subscriptionJSON)
	return err
}

func (c *Client) listenWebSocket() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("Failed to read message: %v", err)
			close(c.orderBookChannel)
			return
		}

		var orderBookData OrderBookData
		err = json.Unmarshal(message, &orderBookData)
		if err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		c.orderBookChannel <- orderBookData
	}
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetOrderBookChannel() chan OrderBookData {
	return c.orderBookChannel
}
