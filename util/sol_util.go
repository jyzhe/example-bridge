package util

import (
	"context"
	"encoding/json"
	"example/constants"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type subscriptionRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// Struct to handle account notifications from Solana
type accountUpdate struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		Result struct {
			Value struct {
				Lamports int `json:"lamports"`
			} `json:"value"`
		} `json:"result"`
	} `json:"params"`
}

// GetSOLBalance fetches the SOL balance for a given token account
func GetSOLBalance(ctx context.Context, addr string) (int, error) {
	// Create a new WebSocket connection to Solana's Devnet
	conn, _, err := websocket.DefaultDialer.Dial(constants.WS_URL, http.Header{})
	if err != nil {
		return 0, fmt.Errorf("error connecting to WebSocket: %v", err)
	}
	defer conn.Close()

	// Create the subscription request
	subscriptionReq := subscriptionRequest{
		Jsonrpc: "2.0",
		ID:      1,
		Method:  constants.ACCOUNT_SUBSCRIPTION,
		Params: []interface{}{
			addr,
			map[string]interface{}{
				"encoding":   "jsonParsed",
				"commitment": "finalized",
			},
		},
	}

	// Submit the subscription request to the WebSocket
	err = conn.WriteJSON(subscriptionReq)
	if err != nil {
		return 0, fmt.Errorf("error sending subscription request: %v", err)
	}

	fmt.Printf("Subscribed to account: %s\n", addr)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Exiting subscription loop")
			return 0, fmt.Errorf("context cancelled")
		default:
			// Wait for messages from the WebSocket
			_, message, err := conn.ReadMessage()
			if err != nil {
				return 0, fmt.Errorf("error reading from WebSocket: %w", err)
			}

			// Parse the incoming message to detect balance changes
			var update accountUpdate
			if err := json.Unmarshal(message, &update); err != nil {
				log.Printf("Error parsing update: %v", err)
				continue
			}

			// Check if the message contains a balance change
			if update.Method == "accountNotification" {
				lamports := update.Params.Result.Value.Lamports

				solana := float64(lamports) / 1000000000 // Convert lamports to SOL
				fmt.Printf("Balance updated: %.9f SOL\n", solana)

				return lamports, nil
			}
		}
	}
}
