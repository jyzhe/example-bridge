package node

import (
	"context"
	"example/util"
	"fmt"
	"log"
	"net/http"

	"github.com/mr-tron/base58"
)

// TODO: in production, break this into two separate services. One for generating
// deposit addresses and another for listening to balance updates.
// This way they can scale independently.

// Generates a new deposit address for the given destination address.
func (n *Node) GenerateDepositAddress(w http.ResponseWriter, r *http.Request) {
	// TODO: perform some validation on the destination address.
	destination := r.URL.Query().Get("dest")
	if destination == "" {
		http.Error(w, "missing destination address", http.StatusBadRequest)
		return
	}

	// Generate a new public and private key pairs.
	newKey, err := util.GenerateNewEd25519Keys()
	if err != nil {
		http.Error(w, "error generating deposit address", http.StatusInternalServerError)
		return
	}

	// Encode the public key into Solana's base58 address format
	addr := base58.Encode(newKey.PublicKey)

	// Launch a goroutine to listen to balance updates.
	go func() {
		// Use background context for POC.
		//
		// TODO: in production, we need to properly handle timeouts, e.g. need to deal
		// with the case where user deposits to the address after timeout.
		n.Subscribe(context.Background(), destination, addr)
	}()

	fmt.Fprintf(w, "%s", addr)
}

// Subscribe listens to balance updates for the given deposit address and
// sends the balance to the destination address.
//
// TODO: properly handle errors.
func (n *Node) Subscribe(ctx context.Context, destination, deposit string) {
	balance, err := util.GetSOLBalance(ctx, deposit)
	if err != nil {
		log.Fatalf("Error getting SOL balance: %v", err)
	}

	// TODO: send the tokens back to the source address if deposits fail for
	// whatever reason on the destination chain.
	err = util.SpotSend(destination, balance)
	if err != nil {
		log.Fatalf("Error sending SOL: %v", err)
	}
}
