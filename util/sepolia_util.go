package util

import (
	"context"
	"crypto/ecdsa"
	"example/constants"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// SpotSend sends a specified amount of SOL to a destination address.
//
// Note that this method does not mint / issue any new tokens (out of scope for this project)
// but sends the existing tokens from a pre-funded account.
func SpotSend(destination string, amount int) error {
	// Connect directly to the Ethereum node (JSON-RPC)
	client, err := ethclient.Dial(constants.ETH_NODE_URL)
	if err != nil {
		return fmt.Errorf("failed to connect to the Ethereum client: %v", err)
	}

	// Load asset pool's private key
	privateKey, err := crypto.HexToECDSA(constants.TEST_ACCOUNT_PRIVATE_KEY)
	if err != nil {
		return fmt.Errorf("error loading private key: %v", err)
	}

	// Get public key from the private key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("Sending SOL from: %s, to %s\n", fromAddress.Hex(), destination)

	// Get the current nonce for account
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %v", err)
	}

	// Load the contract ABI
	parsedABI, err := abi.JSON(strings.NewReader(constants.ERC_20_ABI))
	if err != nil {
		return fmt.Errorf("failed to parse ERC-20 ABI: %v", err)
	}

	// Address of the recipient
	toAddress := common.HexToAddress(destination)

	// Amount of SOL to send (in SOL's smallest unit, 6 decimals)
	bigAmount := new(big.Int).SetInt64(int64(amount))

	// Create an instance of the SOL contract
	contractAddress := common.HexToAddress(constants.TARGET_CONTRACT_ADDRESS)
	callData, err := parsedABI.Pack("transfer", toAddress, bigAmount)
	if err != nil {
		return fmt.Errorf("failed to pack the transaction data: %v", err)
	}

	// Suggest gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get suggested gas price: %v", err)
	}

	// Create the transaction
	tx := types.NewTransaction(
		nonce,
		contractAddress,
		big.NewInt(0),
		uint64(300_000),
		gasPrice,
		callData,
	)

	// Get the chain ID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %v", err)
	}

	// Sign the transaction with the private key
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Broadcast the signed transaction to the Ethereum network
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %v", err)
	}

	// Print the transaction hash
	fmt.Printf("Transaction sent! TX Hash: %s\n", signedTx.Hash().Hex())
	return nil
}
