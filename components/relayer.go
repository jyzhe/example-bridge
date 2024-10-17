package components

import "fmt"

// POC relayer implementation.
type SourceRelayer struct {
	SourceChain      *Chain
	IntermediateNode *IntermediateNode
}

func NewSourceRelayer(
	sourceChain *Chain,
	intermediateNode *IntermediateNode,
) *SourceRelayer {
	return &SourceRelayer{
		SourceChain:      sourceChain,
		IntermediateNode: intermediateNode,
	}
}

func (r *SourceRelayer) Relay() {
	// This is where the relayer would relay the transaction to the other chain.
	fmt.Println("Relaying transaction from source chain")
	txns := r.SourceChain.GetNewTransactions()
	for _, txn := range txns {
		r.IntermediateNode.ConfirmDepositOnSource(txn)
	}
}

type DestinationRelayer struct {
	DestinationChain *Chain
	IntermediateNode *IntermediateNode
}

func NewDestinationRelayer(
	destinationChain *Chain,
	intermediateNode *IntermediateNode,
) *DestinationRelayer {
	return &DestinationRelayer{
		DestinationChain: destinationChain,
		IntermediateNode: intermediateNode,
	}
}

func (r *DestinationRelayer) Relay() {
	// This is where the relayer would relay the transaction to the other chain.
	fmt.Println("Relaying transaction from destination chain")
	txns := r.DestinationChain.GetNewTransactions()
	for _, txn := range txns {
		r.IntermediateNode.ConfirmDepositOnDestination(txn)
	}
}
