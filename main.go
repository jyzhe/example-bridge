package main

import "example/node"

func main() {
	// Start the node that handles the key generation and transaction forwarding.
	node := node.NewNode()
	node.Start()
}
