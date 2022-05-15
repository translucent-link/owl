package index

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/translucent-link/owl/graph/model"
)

type CompoundLiquidateBorrowEvent struct {
	SeizeTokens      *big.Int
	RepayAmount      *big.Int
	Liquidator       common.Address
	Borrower         common.Address
	CTokenCollateral common.Address
}

func (e CompoundLiquidateBorrowEvent) String() string {
	return fmt.Sprintf("liquidator: %s, borrower: %s, tokenCollateral: %s, repayAmount: %s, seizeTokens: %s", e.Liquidator, e.Borrower, e.CTokenCollateral, e.RepayAmount.String(), e.SeizeTokens.String())
}

func (e CompoundLiquidateBorrowEvent) Type() string {
	return "LiquidateBorrow"
}

type CompoundBorrowEvent struct {
	BorrowAmount   *big.Int
	AccountBorrows *big.Int
	TotalBorrows   *big.Int
	Borrower       common.Address
}

func (e CompoundBorrowEvent) String() string {
	return fmt.Sprintf("borrower: %s, borrowAmount: %s, accountBorrows: %s, totalBorrows: %s", e.Borrower, e.BorrowAmount.String(), e.AccountBorrows.String(), e.TotalBorrows.String())
}

func (e CompoundBorrowEvent) Type() string {
	return "Borrow"
}

type CompoundRepayBorrowEvent struct {
	RepayAmount    *big.Int
	AccountBorrows *big.Int
	TotalBorrows   *big.Int
	Borrower       common.Address
	Payer          common.Address
}

func (e CompoundRepayBorrowEvent) String() string {
	return fmt.Sprintf("payer: %s, repayAmount: %s, borrower: %s, accountBorrows: %s, totalBorrows: %s", e.Payer.Hex(), e.RepayAmount.String(), e.Borrower.Hex(), e.AccountBorrows.String(), e.TotalBorrows.String())
}

func (e CompoundRepayBorrowEvent) Type() string {
	return "RepayBorrow"
}

func UnpackCompoundEvent(abi abi.ABI, eventDefn *model.EventDefn, log types.Log) (Typable, error) {
	if eventDefn.TopicName == "Borrow" {
		event := CompoundBorrowEvent{}
		err := abi.UnpackIntoInterface(&event, eventDefn.TopicName, log.Data)
		return event, err
	} else if eventDefn.TopicName == "RepayBorrow" {
		event := CompoundRepayBorrowEvent{}
		err := abi.UnpackIntoInterface(&event, eventDefn.TopicName, log.Data)
		return event, err
	} else if eventDefn.TopicName == "LiquidateBorrow" {
		event := CompoundLiquidateBorrowEvent{}
		err := abi.UnpackIntoInterface(&event, eventDefn.TopicName, log.Data)
		return event, err
	}
	return nil, errors.New(fmt.Sprintf("%s topic name is not supported", eventDefn.TopicName))
}

// func GetCEthDefn(abiPath string) (Protocol, error) {
// 	borrowEventSignature := []byte("Borrow(address,uint256,uint256,uint256)")
// 	borrowTopicHash := crypto.Keccak256Hash(borrowEventSignature)

// 	repayBorrowEventSignature := []byte("RepayBorrow(address,address,uint256,uint256,uint256)")
// 	repayBorrowTopicHash := crypto.Keccak256Hash(repayBorrowEventSignature)

// 	liquidationBorrowEventSignature := []byte("LiquidateBorrow(address,address,uint256,address,uint256)")
// 	liquidateBorrowTopicHash := crypto.Keccak256Hash(liquidationBorrowEventSignature)

// 	f, err := os.Open(path.Join(abiPath, "cETH.abi"))
// 	if err != nil {
// 		return Protocol{}, err
// 	}
// 	contractAbi, err := abi.JSON(bufio.NewReader(f))
// 	if err != nil {
// 		return Protocol{}, err
// 	}

// 	return Protocol{
// 		"CEth",
// 		common.HexToAddress("0x4Ddc2D193948926D02f9B1fE9e1daa0718270ED5"),
// 		[]LoanEvent{
// 			{"Borrow", borrowTopicHash, borrowTopicHash.Hex()},
// 			{"RepayBorrow", repayBorrowTopicHash, repayBorrowTopicHash.Hex()},
// 			{"LiquidateBorrow", liquidateBorrowTopicHash, liquidateBorrowTopicHash.Hex()},
// 		},
// 		[]common.Hash{
// 			borrowTopicHash, repayBorrowTopicHash, liquidateBorrowTopicHash,
// 		},
// 		contractAbi,
// 		unpackCompoundEvent,
// 	}, nil
// }
