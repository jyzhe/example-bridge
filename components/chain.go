package components

import (
	"example/types"
	"fmt"
)

// Simulate a running blockchain chain.
type Chain struct {
	ChainId   uint32
	ChainName string

	// Txns
	Txns []types.CommittedTx

	// Application state.
	Balances map[string]uint32
}

func NewChain(chainId uint32, chainName string) *Chain {
	return &Chain{
		ChainId:   chainId,
		ChainName: chainName,
		Balances:  make(map[string]uint32),
	}
}

// Transfer simulates a transfer of tokens from one address to another.
func (c *Chain) Transfer(tx types.TransferTx) {
	source := tx.Source
	destination := tx.Destination
	amount := tx.Amount

	// Perform some basic validation.
	if tx.Amount == 0 {
		c.Txns = append(
			c.Txns,
			types.CommittedTx{
				Tx:       tx,
				TxStatus: types.STATUS_FAIL,
			},
		)
		return
	}

	// Commit the changes to application state.
	// First verify that the source has enough balance
	if balance, ok := c.Balances[source]; !ok || balance < amount {
		c.Txns = append(
			c.Txns,
			types.CommittedTx{
				Tx:       tx,
				TxStatus: types.STATUS_FAIL,
			},
		)
		return
	}

	c.Txns = append(
		c.Txns,
		types.CommittedTx{
			Tx:       tx,
			TxStatus: types.STATUS_OK,
		},
	)

	// Transfer amount from source to destination
	c.Balances[source] -= amount
	c.Balances[destination] += amount
	fmt.Printf("Transferred %d from %s to %s on %s\n", amount, source, destination, c.ChainName)
}

// GetBalance returns the balance of an address.
func (c *Chain) GetBalance(address string) uint32 {
	if balance, ok := c.Balances[address]; ok {
		return balance
	}
	return 0
}

// GetNewTransactions returns the new transactions that have been added to the chain.
func (c *Chain) GetNewTransactions() []types.CommittedTx {
	result := c.Txns
	c.Txns = make([]types.CommittedTx, 0)
	return result
}
