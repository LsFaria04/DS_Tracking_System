package blockchain

import (
	"log"
	"github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/ethclient"
)

//Deploys a contract
func DeployContract() (common.Address, error) {
    client, err := NewClient()
    if err != nil {
        return common.Address{}, err
    }
	auth := client.Auth
	ethClient := client.EthClient
    address, tx, _, err := DeployBlockchain(auth, ethClient)
    if err != nil {
        return common.Address{}, err
    }

    log.Printf("Contract deployed at: %s", address.Hex())
    log.Printf("Transaction hash: %s", tx.Hash().Hex())
    return address, nil
}

// Gets a instance of a contract
func GetContractInstance(client *ethclient.Client, contractAddress string) (*Blockchain, error) {
    addr := common.HexToAddress(contractAddress)
    instance, err := NewBlockchain(addr, client)
    if err != nil {
        return nil, err
    }
    return instance, nil
}