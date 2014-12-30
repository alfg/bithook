package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	c.read()
}

// Subscribe to new blocks data.
func (c *connection) subBlocks() {
	byt := []byte(`{"op": "blocks_sub"}`)
	c.send(websocket.TextMessage, byt)
	c.read()
}

// Subscribe to new unconfirmed transactions.
func (c *connection) subUnconfirmed() {
	byt := []byte(`{"op": "unconfirmed_sub"}`)
	c.send(websocket.TextMessage, byt)
	c.read()
}

// Debug: Test ping. Returns latest address transaction.
func (c *connection) debugPing() {
	byt := []byte(`{"op": "ping_tx"}`)
	c.send(websocket.TextMessage, byt)
	c.read()
}

// Debug: Test Ping block. Returns latest block transaction.
func (c *connection) debugPingBlock() {
	byt := []byte(`{"op": "ping_block"}`)
	c.send(websocket.TextMessage, byt)
	c.read()
}

type response struct {
	Data map[string]interface{} `json:"data"`
}

// Listen for new messages on websocket forloop.
func (c *connection) read() {
	for {
		var dat map[string]interface{}
		if err := c.ws.ReadJSON(&dat); err != nil {
			log.Fatal(err)
		}
		fmt.Println(dat)

		jsonResponse := &response{Data: dat}
		res, _ := json.Marshal(jsonResponse)
		webHook(res)
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

func webHook(data []byte) {
	url := "http://requestb.in/140q6so1"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func connect() *connection {
	fmt.Println("Starting connection...")

	dialer := websocket.Dialer{}
	conn, _, _ := dialer.Dial(url, nil)

	c := &connection{ws: conn}

	// Set Pinger keep-alive
	c.setPinger()

	return c
}

func parseArgs(args []string) {
	option := args[0]

	switch option {
	case "block", "blocks":
		fmt.Println("Subscribing to blocks")
		c := connect()
		c.subBlocks()

	case "unconfirmed":
		fmt.Println("Subscribing to unconfirmed addresses")
		c := connect()
		c.subUnconfirmed()

	case "address", "addr":
		if len(args) < 2 {
			fmt.Println("Please enter an address")
			os.Exit(1)
		}

		addr := args[1]
		fmt.Printf("Subscribing to address: %s\n", addr)
		c := connect()
		c.subAddress(addr)

	case "test":
		fmt.Println("Subscribing to test ping.")
		c := connect()
		c.debugPing()

	default:
		fmt.Println("Not a valid command.")

	}
}

func main() {
	args := os.Args[1:]
	parseArgs(args)
}
