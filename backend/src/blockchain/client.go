package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client represents a connection to the Ethereum blockchain
type Client struct {
	EthClient       *ethclient.Client
	Auth            *bind.TransactOpts
	ChainID         *big.Int
	WalletAddress   common.Address
	ContractAddress common.Address
}

// NewClient creates a new blockchain client connected to Sepolia testnet
func NewClient() (*Client, error) {
	// Get RPC URL from environment
	rpcURL := os.Getenv("BLOCKCHAIN_RPC_URL")
	if rpcURL == "" {
		return nil, fmt.Errorf("BLOCKCHAIN_RPC_URL environment variable not set")
	}

	// Connect to Ethereum node (Sepolia via Infura)
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}

	// Get private key from environment
	privateKeyHex := os.Getenv("BLOCKCHAIN_PRIVATE_KEY")
	if privateKeyHex == "" {
		return nil, fmt.Errorf("BLOCKCHAIN_PRIVATE_KEY environment variable not set")
	}

	// Remove 0x prefix if present
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")

	// Parse private key
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Get public key and wallet address
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to get public key")
	}
	walletAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Sepolia Chain ID is 11155111
	chainID := big.NewInt(11155111)

	// Create authenticated transactor for signing transactions
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	// Parse contract address if provided
	var contractAddress common.Address
	contractAddressStr := os.Getenv("BLOCKCHAIN_CONTRACT_ADDRESS")
	if contractAddressStr != "" {
		contractAddress = common.HexToAddress(contractAddressStr)
	}

	log.Printf("Connected to Sepolia testnet via %s", rpcURL)
	log.Printf("Wallet address: %s", walletAddress.Hex())
	if contractAddressStr != "" {
		log.Printf("Contract address: %s", contractAddress.Hex())
	}

	return &Client{
		EthClient:       client,
		Auth:            auth,
		ChainID:         chainID,
		WalletAddress:   walletAddress,
		ContractAddress: contractAddress,
	}, nil
}

// Close closes the connection to the Ethereum node
func (c *Client) Close() {
	if c.EthClient != nil {
		c.EthClient.Close()
		log.Println("Disconnected from Ethereum node")
	}
}

// GetBalance returns the ETH balance of a given address
func (c *Client) GetBalance(address string) (*big.Int, error) {
	account := common.HexToAddress(address)
	balance, err := c.EthClient.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}
	return balance, nil
}

// GetWalletBalance returns the ETH balance of the configured wallet
func (c *Client) GetWalletBalance() (*big.Int, error) {
	return c.GetBalance(c.WalletAddress.Hex())
}

// GetBlockNumber returns the current block number
func (c *Client) GetBlockNumber() (uint64, error) {
	blockNumber, err := c.EthClient.BlockNumber(context.Background())
	if err != nil {
		return 0, fmt.Errorf("failed to get block number: %w", err)
	}
	return blockNumber, nil
}

// GetNetworkID returns the network/chain ID
func (c *Client) GetNetworkID() (*big.Int, error) {
	networkID, err := c.EthClient.NetworkID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get network ID: %w", err)
	}
	return networkID, nil
}

// FormatBalance converts Wei to ETH (divides by 10^18)
func FormatBalance(balanceWei *big.Int) string {
	// Convert Wei to ETH (1 ETH = 10^18 Wei)
	fBalance := new(big.Float)
	fBalance.SetString(balanceWei.String())
	ethValue := new(big.Float).Quo(fBalance, big.NewFloat(1e18))
	
	return fmt.Sprintf("%.6f ETH", ethValue)
}

