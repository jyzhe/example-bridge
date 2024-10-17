package components

import (
	"example/types"
	"fmt"
)

const ASSET_POOL = "ASSETPOOL"

// The intermediate node is responsible for generating key pairs and listens for
// updates from the relayer. Once the message is received form the relayer,
// the intermediate node will send the corresponding transaction to the
type IntermediateNode struct {
	// Maintain a two way mapping between desosit addresses and destination addresses.
	DespositToDestination map[string]string
	DesinationToDeposit   map[string]string

	// Maintain a two way mapping between source addresses and deposit addresses.
	SourceToDeposit map[string]string
	DepositToSource map[string]string

	SourceChain      *Chain
	DestinationChain *Chain
}

func NewIntermediateNode(sourceChain *Chain, destinationChain *Chain) *IntermediateNode {
	return &IntermediateNode{
		DespositToDestination: make(map[string]string),
		DesinationToDeposit:   make(map[string]string),
		SourceToDeposit:       make(map[string]string),
		DepositToSource:       make(map[string]string),

		SourceChain:      sourceChain,
		DestinationChain: destinationChain,
	}
}

func (i *IntermediateNode) GenerateDepositAddress(destination string) string {
	// Generate a destination address for the source address.
	// This is the address that users can deposit to.
	//
	// Public / private key pairs are also generated here, omitted for brevity.
	depositAddress := RandString(8)
	i.DespositToDestination[depositAddress] = destination
	i.DesinationToDeposit[destination] = depositAddress
	fmt.Printf("Generated deposit address: %s for destination address: %s\n", depositAddress, destination)
	return depositAddress
}

func (i *IntermediateNode) ConfirmDepositOnSource(txn types.CommittedTx) {
	if txn.TxStatus != types.STATUS_OK {
		fmt.Println("Error: Transaction failed on source chain")
		return
	}

	// See if the destination address is one of our intermediate deposit addresses.
	depositAddress := txn.Tx.Destination
	if destination, ok := i.DespositToDestination[depositAddress]; ok {
		// Send a transaction to the destination chain.
		i.DestinationChain.Transfer(
			types.TransferTx{
				Source:      ASSET_POOL,
				Destination: destination,
				Amount:      txn.Tx.Amount,
			},
		)
		fmt.Printf("Sending %d to destination address: %s\n", txn.Tx.Amount, destination)

		// Also maintain the mapping between the source and deposit addresses.
		i.SourceToDeposit[txn.Tx.Source] = depositAddress
		i.DepositToSource[depositAddress] = txn.Tx.Source
	}
}

func (i *IntermediateNode) ConfirmDepositOnDestination(txn types.CommittedTx) {
	if txn.TxStatus != types.STATUS_OK {
		fmt.Println("Error: Transaction failed on destination chain")
		// The transaction failed on the destination chain.
		// Now we need to release the funds back to the source address.
		depositAddress := i.DesinationToDeposit[txn.Tx.Destination]
		sourceAddress := i.DepositToSource[depositAddress]

		// Send a transaction back to the source chain.
		i.SourceChain.Transfer(
			types.TransferTx{
				Source:      depositAddress,
				Destination: sourceAddress,
				Amount:      txn.Tx.Amount,
			},
		)
		fmt.Printf("Releasing %d back to source address: %s\n", txn.Tx.Amount, sourceAddress)
		// TODO: do a little bit of clean up here.
	} else {
		fmt.Println("Transaction succeeded on destination chain")
	}
}
