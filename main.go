package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	PingInterval = time.Second * 30
	Url          = "wss://ws.blockchain.info:443/inv"
	Usage        = `
	  Usage:
		bithook blocks -- Subscribe to new blocks.
		bithook unconfirmed -- Subscribe to new unconfirmed transactions.
		bithook address <address> -- Subscribe to address.
		bithook test -- Receives latest transaction. Use for testing.
		bithook help -- This help menu.
		bithook version -- This version.
	`
	Version = "0.0.1"
)

var webhookFlag string

func init() {
	args := flag.NewFlagSet("", flag.ExitOnError)
	args.StringVar(&webhookFlag, "webhook", "http://requestb.in/140q6so1", "Webhook URL.")
	args.Parse(os.Args[2:])
	flag.Parse()
}

type connection struct {
	ws   *websocket.Conn
	conn *websocket.Dialer
}

type response struct {
	Data map[string]interface{} `json:"data"`
}

// Send messages wrapper.
func (c *connection) send(messageType int, payload []byte) {
	fmt.Printf("Send: %s\n", payload)
	c.ws.WriteMessage(messageType, payload)

	/*
		// Blockchain websocket accepts message strings
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
	ticker := time.NewTicker(PingInterval)

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

// Sends POST request along with json data
func webHook(data []byte) {
	fmt.Println("Sending Request: ", webhookFlag)
	url := webhookFlag

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	ua := fmt.Sprintf("bithook-client-%s", Version)
	req.Header.Set("User-Agent", ua)
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

// Creates and returns connection
func connect() *connection {
	fmt.Println("Starting connection...")

	dialer := websocket.Dialer{}
	conn, _, _ := dialer.Dial(Url, nil)

	c := &connection{ws: conn}

	// Set Pinger keep-alive
	c.setPinger()

	return c
}

// Parses and validates cli arguments
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

	case "version":
		fmt.Printf("Version: %s\n", Version)

	default:
		fmt.Println("Not a valid command.")
		fmt.Println(Usage)

	}
}

func main() {
	if len(os.Args[1:]) < 1 {
		fmt.Println("Please enter a command.")
		os.Exit(1)
	}

	args := os.Args[1:]
	parseArgs(args)
}
