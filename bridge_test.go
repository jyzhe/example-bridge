package test

import (
	"testing"

	"example/components"
	"example/types"
)

const (
	USER_SOURCE_ADDRESS      = "ABC"
	USER_DESTINATION_ADDRESS = "XYZ"
)

func TestBridge_Success(t *testing.T) {
	// Initialize the source and destination chains.
	source := components.NewChain(1, "Source")
	destination := components.NewChain(2, "Destination")

	// Initialize the intermediate node.
	intermediateNode := components.NewIntermediateNode(source, destination)

	// Initialize the relayers.
	sourceRelayer := components.NewSourceRelayer(source, intermediateNode)
	destinationRelayer := components.NewDestinationRelayer(destination, intermediateNode)

	// For this example, let's set the balance of the alloyed asset pool to 1000 on the destination chain.
	source.Balances[USER_SOURCE_ADDRESS] = 1_000_000
	destination.Balances[components.ASSET_POOL] = 1_000

	// User flow
	// 1. User wants to deposit 1000 tokens to the address "XYZ" on the destination chain.
	// 2. FE/client gets a deposit address from the intermediate node.
	despositAddress := intermediateNode.GenerateDepositAddress(USER_DESTINATION_ADDRESS)

	// 3. User deposits 1000 tokens to the deposit address.
	source.Transfer(types.TransferTx{
		Source:      USER_SOURCE_ADDRESS,
		Destination: despositAddress,
		Amount:      1_000,
	})
	// 4. Source chain relayer listens for transactions, and confirms that transaction was indeed
	// included onchain and waits for finality. then relays the transaction to the intermediate node.
	// 5. Intermediate node sends the spotSend on the destination chain.
	sourceRelayer.Relay()

	// 6. Destination chain relayer listens for transactions, and confirms that transaction was successful.
	destinationRelayer.Relay()

	if source.GetBalance(USER_SOURCE_ADDRESS) != 999_000 {
		t.Errorf("Expected source balance to be 999_000, got %d", source.GetBalance(USER_SOURCE_ADDRESS))
	}

	if destination.GetBalance(USER_DESTINATION_ADDRESS) != 1_000 {
		t.Errorf("Expected destination balance to be 1_000, got %d", destination.GetBalance(USER_DESTINATION_ADDRESS))
	}
}

func TestBridge_Failure(t *testing.T) {
	// Initialize the source and destination chains.
	source := components.NewChain(1, "Source")
	destination := components.NewChain(2, "Destination")

	// Initialize the intermediate node.
	intermediateNode := components.NewIntermediateNode(source, destination)

	// Initialize the relayers.
	sourceRelayer := components.NewSourceRelayer(source, intermediateNode)
	destinationRelayer := components.NewDestinationRelayer(destination, intermediateNode)

	// For this example, let's set the balance of the alloyed asset pool to 1000 on the destination chain.
	source.Balances[USER_SOURCE_ADDRESS] = 1_000_000
	destination.Balances[components.ASSET_POOL] = 1_000

	// User flow
	// 1. User wants to deposit 1000 tokens to the address "XYZ" on the destination chain.
	// 2. FE/client gets a deposit address from the intermediate node.
	despositAddress := intermediateNode.GenerateDepositAddress(USER_DESTINATION_ADDRESS)

	// 3. User deposits 2000 tokens to the deposit address.
	source.Transfer(types.TransferTx{
		Source:      USER_SOURCE_ADDRESS,
		Destination: despositAddress,
		Amount:      2_000, // More than the balance of the alloyed asset pool.
	})
	// 4. Source chain relayer listens for transactions, and confirms that transaction was indeed
	// included onchain and waits for finality. then relays the transaction to the intermediate node.
	// 5. Intermediate node sends the spotSend on the destination chain.
	sourceRelayer.Relay()

	// 6. Destination chain relayer listens for transactions,
	// intermediate nodes releases funds from the deposit address when txn fails.
	destinationRelayer.Relay()

	if source.GetBalance(USER_SOURCE_ADDRESS) != 1_000_000 {
		t.Errorf("Expected source balance to be 1_000_000, got %d", source.GetBalance(USER_SOURCE_ADDRESS))
	}

	if destination.GetBalance(USER_DESTINATION_ADDRESS) != 0 {
		t.Errorf("Expected destination balance to be 0, got %d", destination.GetBalance(USER_DESTINATION_ADDRESS))
	}
}
