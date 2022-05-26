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
	DebtToken        common.Address
}

func (e CompoundLiquidateBorrowEvent) String() string {
	return fmt.Sprintf("liquidator: %s, borrower: %s, tokenCollateral: %s, repayAmount: %s, seizeTokens: %s", e.Liquidator, e.Borrower, e.CTokenCollateral, e.RepayAmount.String(), e.SeizeTokens.String())
}

func (e CompoundLiquidateBorrowEvent) Type() string {
	return "LiquidateBorrow"
}

func (e CompoundLiquidateBorrowEvent) GetRepayAmount() *big.Int {
	return e.RepayAmount
}

func (e CompoundLiquidateBorrowEvent) GetSeizeAmount() *big.Int {
	return e.SeizeTokens
}

func (e CompoundLiquidateBorrowEvent) GetBorrower() common.Address {
	return e.Borrower
}

func (e CompoundLiquidateBorrowEvent) GetDebtToken() common.Address {
	return e.DebtToken
}

func (e CompoundLiquidateBorrowEvent) GetCollateralToken() common.Address {
	return e.CTokenCollateral
}

func (e CompoundLiquidateBorrowEvent) GetLiquidator() common.Address {
	return e.Liquidator
}

type CompoundBorrowEvent struct {
	BorrowAmount   *big.Int
	AccountBorrows *big.Int
	TotalBorrows   *big.Int
	Borrower       common.Address
	BorrowToken    common.Address
}

func (e CompoundBorrowEvent) GetBorrowAmount() *big.Int {
	return e.BorrowAmount
}

func (e CompoundBorrowEvent) GetBorrower() common.Address {
	return e.Borrower
}

func (e CompoundBorrowEvent) GetBorrowToken() common.Address {
	return e.BorrowToken
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
	BorrowToken    common.Address
}

func (e CompoundRepayBorrowEvent) String() string {
	return fmt.Sprintf("payer: %s, repayAmount: %s, borrower: %s, accountBorrows: %s, totalBorrows: %s", e.Payer.Hex(), e.RepayAmount.String(), e.Borrower.Hex(), e.AccountBorrows.String(), e.TotalBorrows.String())
}

func (e CompoundRepayBorrowEvent) Type() string {
	return "RepayBorrow"
}

func (e CompoundRepayBorrowEvent) GetRepayAmount() *big.Int {
	return e.RepayAmount
}

func (e CompoundRepayBorrowEvent) GetBorrower() common.Address {
	return e.Borrower
}

func (e CompoundRepayBorrowEvent) GetBorrowToken() common.Address {
	return e.BorrowToken
}

func (e CompoundRepayBorrowEvent) GetRepayer() common.Address {
	return e.Payer
}

type CompoundMintEvent struct {
	MintAmount   *big.Int
	MintTokens   *big.Int
	Minter       common.Address
	DepositToken common.Address
}

func (e CompoundMintEvent) GetDepositAmount() *big.Int {
	return e.MintAmount
}

func (e CompoundMintEvent) GetDepositor() common.Address {
	return e.Minter
}

func (e CompoundMintEvent) GetDepositToken() common.Address {
	return e.DepositToken
}

func (e CompoundMintEvent) String() string {
	return fmt.Sprintf("depositor: %s, despositAmount: %s, token %s", e.Minter, e.MintAmount.String(), e.DepositToken.Hex())
}

func (e CompoundMintEvent) Type() string {
	return "Mint"
}

func UnpackCompoundEvent(abi abi.ABI, protocolInstance *model.ProtocolInstance, eventDefn *model.EventDefn, log types.Log) (PossibleEvent, error) {
	if eventDefn.TopicName == "Borrow" {
		event := CompoundBorrowEvent{}
		event.BorrowToken = common.HexToAddress(protocolInstance.ContractAddress)
		err := abi.UnpackIntoInterface(&event, eventDefn.TopicName, log.Data)
		return PossibleEvent{Borrowable: event}, err
	} else if eventDefn.TopicName == "RepayBorrow" {
		event := CompoundRepayBorrowEvent{}
		event.BorrowToken = common.HexToAddress(protocolInstance.ContractAddress)
		err := abi.UnpackIntoInterface(&event, eventDefn.TopicName, log.Data)
		return PossibleEvent{Repayable: event}, err
	} else if eventDefn.TopicName == "LiquidateBorrow" {
		event := CompoundLiquidateBorrowEvent{}
		event.DebtToken = common.HexToAddress(protocolInstance.ContractAddress)
		err := abi.UnpackIntoInterface(&event, eventDefn.TopicName, log.Data)
		return PossibleEvent{Liquidatable: event}, err
	} else if eventDefn.TopicName == "Mint" {
		event := CompoundMintEvent{}
		err := abi.UnpackIntoInterface(&event, eventDefn.TopicName, log.Data)
		event.DepositToken = common.HexToAddress(protocolInstance.ContractAddress)
		return PossibleEvent{Depositable: event}, err
	}
	return PossibleEvent{}, errors.New(fmt.Sprintf("%s topic name is not supported", eventDefn.TopicName))
}
