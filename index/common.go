package index

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/translucent-link/owl/graph/model"
)

type Unpacker func(abi abi.ABI, protocolInstance *model.ProtocolInstance, eventDefn *model.EventDefn, log types.Log) (PossibleEvent, error)

type Chain struct {
	RpcURL         string
	BlockFetchSize int
}

type Depositable interface {
	GetDepositAmount() *big.Int
	GetDepositor() common.Address
	GetDepositToken() common.Address
	Typable
}

type Borrowable interface {
	GetBorrowAmount() *big.Int
	GetBorrower() common.Address
	GetBorrowToken() common.Address
	Typable
}

type Repayable interface {
	GetRepayAmount() *big.Int
	GetBorrower() common.Address
	GetBorrowToken() common.Address
	GetRepayer() common.Address
	Typable
}

type Liquidatable interface {
	GetSeizeAmount() *big.Int
	GetRepayAmount() *big.Int
	GetBorrower() common.Address
	GetLiquidator() common.Address
	GetDebtToken() common.Address
	GetCollateralToken() common.Address
	Typable
}

type PossibleEvent struct {
	Borrowable   Borrowable
	Repayable    Repayable
	Liquidatable Liquidatable
	Depositable  Depositable
}

type Typable interface {
	Type() string
}

type Transactional interface {
	Txn() string
}

func DebugPrint(msg string, value interface{}) {
	fmt.Printf("%s %#v\n", msg, value)
}

func handleLogEvent(vLog types.Log) {
	DebugPrint("Log Event", vLog)
}

func GetClient(ethURL string) (*ethclient.Client, error) {
	return ethclient.Dial(ethURL)
}

func GetBlockHeader(client *ethclient.Client, blockNumber *big.Int) (*types.Header, error) {
	return client.HeaderByNumber(context.Background(), blockNumber)
}

func GetBlockTimestamp(client *ethclient.Client, blockNumber *big.Int) (time.Time, error) {
	header, err := GetBlockHeader(client, blockNumber)
	if err != nil {
		return time.Time{}, errors.Wrapf(err, "Unable to deteremine timestamp for block %d", blockNumber.Int64())
	}
	return time.Unix(int64(header.Time), 0), nil
}

func indexOf(values []string, tgt string) int {
	for i, v := range values {
		if v == tgt {
			return i
		}
	}
	return -1
}
