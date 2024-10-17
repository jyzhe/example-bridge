package components

import (
	"example/types"
	"fmt"
)

type Chain struct {
	ChainId   uint32
	ChainName string

	// Txns
	Txns []types.TransferTx // This should really be a map of block height to []TransferTx

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

func (c *Chain) Transfer(tx types.TransferTx) {
	source := tx.Source
	destination := tx.Destination
	amount := tx.Amount

	// Perform some basic validation.
	if tx.Amount == 0 {
		tx.TxStatus = types.STATUS_FAIL
		c.Txns = append(c.Txns, tx)
		return
	}

	// Commit the changes to application state.
	// First verify that the source has enough balance
	if balance, ok := c.Balances[source]; !ok || balance < amount {
		tx.TxStatus = types.STATUS_FAIL
		c.Txns = append(c.Txns, tx)
		return
	}

	tx.TxStatus = types.STATUS_OK
	c.Txns = append(c.Txns, tx)

	// Transfer amount from source to destination
	c.Balances[source] -= amount
	c.Balances[destination] += amount
	fmt.Printf("Transferred %d from %s to %s on %s\n", amount, source, destination, c.ChainName)
}

func (c *Chain) GetBalance(address string) uint32 {
	if balance, ok := c.Balances[address]; ok {
		return balance
	}
	return 0
}

func (c *Chain) GetNewTransactions() []types.TransferTx {
	result := c.Txns
	c.Txns = make([]types.TransferTx, 0)
	return result
}
