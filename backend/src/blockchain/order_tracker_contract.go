// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package blockchain

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// BlockchainMetaData contains all meta data concerning the Blockchain contract.
var BlockchainMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"orderId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"OrderUpdateHashStored\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"orderId\",\"type\":\"uint256\"}],\"name\":\"getUpdateHash\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"orderId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"storeUpdateHash\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"updateHashes\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b5061043c8061001c5f395ff3fe608060405234801561000f575f5ffd5b506004361061003f575f3560e01c806325b15ec1146100435780637ba061131461005f578063f716fc251461008f575b5f5ffd5b61005d6004803603810190610058919061022a565b6100bf565b005b61007960048036038101906100749190610268565b610130565b604051610086919061034a565b60405180910390f35b6100a960048036038101906100a4919061036a565b610196565b6040516100b691906103b7565b60405180910390f35b5f5f8381526020019081526020015f2081908060018154018082558091505060019003905f5260205f20015f90919091909150557fdf9d1da8115cdeb16d6e58a22276651268fc42fc4b955be13b1134686f0498bd82826040516101249291906103df565b60405180910390a15050565b60605f5f8381526020019081526020015f2080548060200260200160405190810160405280929190818152602001828054801561018a57602002820191905f5260205f20905b815481526020019060010190808311610176575b50505050509050919050565b5f602052815f5260405f2081815481106101ae575f80fd5b905f5260205f20015f91509150505481565b5f5ffd5b5f819050919050565b6101d6816101c4565b81146101e0575f5ffd5b50565b5f813590506101f1816101cd565b92915050565b5f819050919050565b610209816101f7565b8114610213575f5ffd5b50565b5f8135905061022481610200565b92915050565b5f5f604083850312156102405761023f6101c0565b5b5f61024d858286016101e3565b925050602061025e85828601610216565b9150509250929050565b5f6020828403121561027d5761027c6101c0565b5b5f61028a848285016101e3565b91505092915050565b5f81519050919050565b5f82825260208201905092915050565b5f819050602082019050919050565b6102c5816101f7565b82525050565b5f6102d683836102bc565b60208301905092915050565b5f602082019050919050565b5f6102f882610293565b610302818561029d565b935061030d836102ad565b805f5b8381101561033d57815161032488826102cb565b975061032f836102e2565b925050600181019050610310565b5085935050505092915050565b5f6020820190508181035f83015261036281846102ee565b905092915050565b5f5f604083850312156103805761037f6101c0565b5b5f61038d858286016101e3565b925050602061039e858286016101e3565b9150509250929050565b6103b1816101f7565b82525050565b5f6020820190506103ca5f8301846103a8565b92915050565b6103d9816101c4565b82525050565b5f6040820190506103f25f8301856103d0565b6103ff60208301846103a8565b939250505056fea2646970667358221220f5dda527ceaf962fe153872cd0d1c0dc4534bb8db0b613cb4f02abce1ba2599c64736f6c634300081e0033",
}

// BlockchainABI is the input ABI used to generate the binding from.
// Deprecated: Use BlockchainMetaData.ABI instead.
var BlockchainABI = BlockchainMetaData.ABI

// BlockchainBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BlockchainMetaData.Bin instead.
var BlockchainBin = BlockchainMetaData.Bin

// DeployBlockchain deploys a new Ethereum contract, binding an instance of Blockchain to it.
func DeployBlockchain(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Blockchain, error) {
	parsed, err := BlockchainMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BlockchainBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Blockchain{BlockchainCaller: BlockchainCaller{contract: contract}, BlockchainTransactor: BlockchainTransactor{contract: contract}, BlockchainFilterer: BlockchainFilterer{contract: contract}}, nil
}

// Blockchain is an auto generated Go binding around an Ethereum contract.
type Blockchain struct {
	BlockchainCaller     // Read-only binding to the contract
	BlockchainTransactor // Write-only binding to the contract
	BlockchainFilterer   // Log filterer for contract events
}

// BlockchainCaller is an auto generated read-only Go binding around an Ethereum contract.
type BlockchainCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BlockchainTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BlockchainTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BlockchainFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BlockchainFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BlockchainSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BlockchainSession struct {
	Contract     *Blockchain       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BlockchainCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BlockchainCallerSession struct {
	Contract *BlockchainCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// BlockchainTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BlockchainTransactorSession struct {
	Contract     *BlockchainTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// BlockchainRaw is an auto generated low-level Go binding around an Ethereum contract.
type BlockchainRaw struct {
	Contract *Blockchain // Generic contract binding to access the raw methods on
}

// BlockchainCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BlockchainCallerRaw struct {
	Contract *BlockchainCaller // Generic read-only contract binding to access the raw methods on
}

// BlockchainTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BlockchainTransactorRaw struct {
	Contract *BlockchainTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBlockchain creates a new instance of Blockchain, bound to a specific deployed contract.
func NewBlockchain(address common.Address, backend bind.ContractBackend) (*Blockchain, error) {
	contract, err := bindBlockchain(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Blockchain{BlockchainCaller: BlockchainCaller{contract: contract}, BlockchainTransactor: BlockchainTransactor{contract: contract}, BlockchainFilterer: BlockchainFilterer{contract: contract}}, nil
}

// NewBlockchainCaller creates a new read-only instance of Blockchain, bound to a specific deployed contract.
func NewBlockchainCaller(address common.Address, caller bind.ContractCaller) (*BlockchainCaller, error) {
	contract, err := bindBlockchain(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BlockchainCaller{contract: contract}, nil
}

// NewBlockchainTransactor creates a new write-only instance of Blockchain, bound to a specific deployed contract.
func NewBlockchainTransactor(address common.Address, transactor bind.ContractTransactor) (*BlockchainTransactor, error) {
	contract, err := bindBlockchain(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BlockchainTransactor{contract: contract}, nil
}

// NewBlockchainFilterer creates a new log filterer instance of Blockchain, bound to a specific deployed contract.
func NewBlockchainFilterer(address common.Address, filterer bind.ContractFilterer) (*BlockchainFilterer, error) {
	contract, err := bindBlockchain(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BlockchainFilterer{contract: contract}, nil
}

// bindBlockchain binds a generic wrapper to an already deployed contract.
func bindBlockchain(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BlockchainMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Blockchain *BlockchainRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Blockchain.Contract.BlockchainCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Blockchain *BlockchainRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Blockchain.Contract.BlockchainTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Blockchain *BlockchainRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Blockchain.Contract.BlockchainTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Blockchain *BlockchainCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Blockchain.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Blockchain *BlockchainTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Blockchain.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Blockchain *BlockchainTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Blockchain.Contract.contract.Transact(opts, method, params...)
}

// GetUpdateHash is a free data retrieval call binding the contract method 0x7ba06113.
//
// Solidity: function getUpdateHash(uint256 orderId) view returns(bytes32[])
func (_Blockchain *BlockchainCaller) GetUpdateHash(opts *bind.CallOpts, orderId *big.Int) ([][32]byte, error) {
	var out []interface{}
	err := _Blockchain.contract.Call(opts, &out, "getUpdateHash", orderId)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetUpdateHash is a free data retrieval call binding the contract method 0x7ba06113.
//
// Solidity: function getUpdateHash(uint256 orderId) view returns(bytes32[])
func (_Blockchain *BlockchainSession) GetUpdateHash(orderId *big.Int) ([][32]byte, error) {
	return _Blockchain.Contract.GetUpdateHash(&_Blockchain.CallOpts, orderId)
}

// GetUpdateHash is a free data retrieval call binding the contract method 0x7ba06113.
//
// Solidity: function getUpdateHash(uint256 orderId) view returns(bytes32[])
func (_Blockchain *BlockchainCallerSession) GetUpdateHash(orderId *big.Int) ([][32]byte, error) {
	return _Blockchain.Contract.GetUpdateHash(&_Blockchain.CallOpts, orderId)
}

// UpdateHashes is a free data retrieval call binding the contract method 0xf716fc25.
//
// Solidity: function updateHashes(uint256 , uint256 ) view returns(bytes32)
func (_Blockchain *BlockchainCaller) UpdateHashes(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Blockchain.contract.Call(opts, &out, "updateHashes", arg0, arg1)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// UpdateHashes is a free data retrieval call binding the contract method 0xf716fc25.
//
// Solidity: function updateHashes(uint256 , uint256 ) view returns(bytes32)
func (_Blockchain *BlockchainSession) UpdateHashes(arg0 *big.Int, arg1 *big.Int) ([32]byte, error) {
	return _Blockchain.Contract.UpdateHashes(&_Blockchain.CallOpts, arg0, arg1)
}

// UpdateHashes is a free data retrieval call binding the contract method 0xf716fc25.
//
// Solidity: function updateHashes(uint256 , uint256 ) view returns(bytes32)
func (_Blockchain *BlockchainCallerSession) UpdateHashes(arg0 *big.Int, arg1 *big.Int) ([32]byte, error) {
	return _Blockchain.Contract.UpdateHashes(&_Blockchain.CallOpts, arg0, arg1)
}

// StoreUpdateHash is a paid mutator transaction binding the contract method 0x25b15ec1.
//
// Solidity: function storeUpdateHash(uint256 orderId, bytes32 hash) returns()
func (_Blockchain *BlockchainTransactor) StoreUpdateHash(opts *bind.TransactOpts, orderId *big.Int, hash [32]byte) (*types.Transaction, error) {
	return _Blockchain.contract.Transact(opts, "storeUpdateHash", orderId, hash)
}

// StoreUpdateHash is a paid mutator transaction binding the contract method 0x25b15ec1.
//
// Solidity: function storeUpdateHash(uint256 orderId, bytes32 hash) returns()
func (_Blockchain *BlockchainSession) StoreUpdateHash(orderId *big.Int, hash [32]byte) (*types.Transaction, error) {
	return _Blockchain.Contract.StoreUpdateHash(&_Blockchain.TransactOpts, orderId, hash)
}

// StoreUpdateHash is a paid mutator transaction binding the contract method 0x25b15ec1.
//
// Solidity: function storeUpdateHash(uint256 orderId, bytes32 hash) returns()
func (_Blockchain *BlockchainTransactorSession) StoreUpdateHash(orderId *big.Int, hash [32]byte) (*types.Transaction, error) {
	return _Blockchain.Contract.StoreUpdateHash(&_Blockchain.TransactOpts, orderId, hash)
}

// BlockchainOrderUpdateHashStoredIterator is returned from FilterOrderUpdateHashStored and is used to iterate over the raw logs and unpacked data for OrderUpdateHashStored events raised by the Blockchain contract.
type BlockchainOrderUpdateHashStoredIterator struct {
	Event *BlockchainOrderUpdateHashStored // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BlockchainOrderUpdateHashStoredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockchainOrderUpdateHashStored)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BlockchainOrderUpdateHashStored)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BlockchainOrderUpdateHashStoredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockchainOrderUpdateHashStoredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockchainOrderUpdateHashStored represents a OrderUpdateHashStored event raised by the Blockchain contract.
type BlockchainOrderUpdateHashStored struct {
	OrderId *big.Int
	Hash    [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterOrderUpdateHashStored is a free log retrieval operation binding the contract event 0xdf9d1da8115cdeb16d6e58a22276651268fc42fc4b955be13b1134686f0498bd.
//
// Solidity: event OrderUpdateHashStored(uint256 orderId, bytes32 hash)
func (_Blockchain *BlockchainFilterer) FilterOrderUpdateHashStored(opts *bind.FilterOpts) (*BlockchainOrderUpdateHashStoredIterator, error) {

	logs, sub, err := _Blockchain.contract.FilterLogs(opts, "OrderUpdateHashStored")
	if err != nil {
		return nil, err
	}
	return &BlockchainOrderUpdateHashStoredIterator{contract: _Blockchain.contract, event: "OrderUpdateHashStored", logs: logs, sub: sub}, nil
}

// WatchOrderUpdateHashStored is a free log subscription operation binding the contract event 0xdf9d1da8115cdeb16d6e58a22276651268fc42fc4b955be13b1134686f0498bd.
//
// Solidity: event OrderUpdateHashStored(uint256 orderId, bytes32 hash)
func (_Blockchain *BlockchainFilterer) WatchOrderUpdateHashStored(opts *bind.WatchOpts, sink chan<- *BlockchainOrderUpdateHashStored) (event.Subscription, error) {

	logs, sub, err := _Blockchain.contract.WatchLogs(opts, "OrderUpdateHashStored")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockchainOrderUpdateHashStored)
				if err := _Blockchain.contract.UnpackLog(event, "OrderUpdateHashStored", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOrderUpdateHashStored is a log parse operation binding the contract event 0xdf9d1da8115cdeb16d6e58a22276651268fc42fc4b955be13b1134686f0498bd.
//
// Solidity: event OrderUpdateHashStored(uint256 orderId, bytes32 hash)
func (_Blockchain *BlockchainFilterer) ParseOrderUpdateHashStored(log types.Log) (*BlockchainOrderUpdateHashStored, error) {
	event := new(BlockchainOrderUpdateHashStored)
	if err := _Blockchain.contract.UnpackLog(event, "OrderUpdateHashStored", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
