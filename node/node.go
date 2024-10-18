package node

import (
	"example/constants"
	"fmt"
	"net/http"
)

type Node struct{}

func NewNode() *Node {
	return &Node{}
}

// Start the intermedaite node that handles the key generation and transaction forwarding.
func (n *Node) Start() {
	// Define a handler for the "/deposit" route with query parameters
	http.HandleFunc("/deposit", n.GenerateDepositAddress)

	// Set up the server to listen on port 8080
	fmt.Printf("Starting server at http://localhost%s\n", constants.PORT)
	if err := http.ListenAndServe(constants.PORT, nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
