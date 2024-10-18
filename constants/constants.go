package constants

// Node config
const (
	PORT = ":8080"
)

// Solana (source chain)
const (
	// Solana Devnet RPC endpoint
	WS_URL = "wss://api.devnet.solana.com"

	// Methods
	ACCOUNT_SUBSCRIPTION = "accountSubscribe"
	ACCOUNT_NOTIFICATION = "accountNotification"
)

// ETH Sepolia (destination chain)
const (
	// Eth node
	ETH_NODE_URL = "https://ethereum-sepolia-rpc.publicnode.com"

	// ERC-20 ABI
	ERC_20_ABI = `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"}]`

	// Target contract address on Ethereum Sepolia testnet (wrapped SOL in this case)
	TARGET_CONTRACT_ADDRESS = "0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238"
)
