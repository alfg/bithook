package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	url          = "wss://ws.blockchain.info:443/inv"
	pingInterval = time.Second * 30
)

type connection struct {
	ws   *websocket.Conn
	conn *websocket.Dialer
}

// Send messages wrapper.
func (c *connection) send(messageType int, payload []byte) {
	fmt.Printf("Send: %s\n", payload)
	c.ws.WriteMessage(messageType, payload)

	/*
		c.ws.WriteJSON(payload)
		fmt.Printf("Send: %s\n", payload)
	*/
}

// Subscribe to new transactions for address.
func (c *connection) subAddress(addr string) {
	addrStr := fmt.Sprintf(`{"op": "addr_sub", "addr": "%s"}`, addr)
	byt := []byte(addrStr)
	c.send(websocket.TextMessage, byt)
}

// Subscribe to new blocks data.
func (c *connection) subBlocks() {
	byt := []byte(`{"op": "blocks_sub"}`)
	c.send(websocket.TextMessage, byt)
}

// Subscribe to new unconfirmed transactions.
func (c *connection) subUnconfirmed() {
	byt := []byte(`{"op": "unconfirmed_sub"}`)
	c.send(websocket.TextMessage, byt)
}

// Debug: Test ping. Returns latest address transaction.
func (c *connection) debugPing() {
	byt := []byte(`{"op": "ping_tx"}`)
	c.send(websocket.TextMessage, byt)
}

// Debug: Test Ping block. Returns latest block transaction.
func (c *connection) debugPingBlock() {
	byt := []byte(`{"op": "ping_block"}`)
	c.send(websocket.TextMessage, byt)
}

// Listen for new messages on websocket forloop.
func (c *connection) read() {
	for {
		var data map[string]interface{}
		if err := c.ws.ReadJSON(&data); err != nil {
			log.Fatal(err)
		}
		fmt.Println(data)
	}
}

// Sends ping every n seconds to keep connection alive
func (c *connection) setPinger() {
	ticker := time.NewTicker(pingInterval)

	go func() {
		for {
			select {
			case <-ticker.C:
				byt := []byte(`"op": "ping"`)
				c.send(websocket.PingMessage, byt)
			}

		}
	}()
}

func main() {
	fmt.Println("Starting connection...")

	dialer := websocket.Dialer{}
	conn, _, _ := dialer.Dial(url, nil)

	c := &connection{ws: conn}

	// Set Pinger keep-alive
	c.setPinger()

	// Subscriptions
	c.debugPing()
	//c.subAddress("1GZ2TY8PT3yNtjLBkoUxwrULE5w77WH7rU")

	// Read messages
	c.read()

}
