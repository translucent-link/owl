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

type Transaction struct {
	TxHash string
}

func (e Transaction) Txn() string {
	return e.TxHash
}

type AaveTransferEvent struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Transaction
}

func (e AaveTransferEvent) String() string {
	return fmt.Sprintf("from: %s, to: %s, amount: %s", e.From, e.To, e.Value)
}

func (e AaveTransferEvent) Type() string {
	return "AaveTransfer"
}

type AaveLiquidationCallEvent struct {
	CollateralAsset            common.Address
	DebtAsset                  common.Address
	User                       common.Address
	DebtToCover                *big.Int
	LiquidatedCollateralAmount *big.Int
	Liquidator                 common.Address
	ReceiveAToken              bool
	Transaction
}

func (e AaveLiquidationCallEvent) String() string {
	return fmt.Sprintf("LiquidationCall liquidator: %s, borrower: %s, tokenCollateral: %s, repayAmount: %s, seizeTokens: %s, txn: %s", e.Liquidator, e.User, e.CollateralAsset, e.DebtToCover.String(), e.LiquidatedCollateralAmount.String(), e.TxHash)
}

func (e AaveLiquidationCallEvent) Type() string {
	return "AaveLiquidationCall"
}

func (e AaveLiquidationCallEvent) GetRepayAmount() *big.Int {
	return e.DebtToCover
}

func (e AaveLiquidationCallEvent) GetSeizeAmount() *big.Int {
	return e.LiquidatedCollateralAmount
}

func (e AaveLiquidationCallEvent) GetBorrower() common.Address {
	return e.User
}

func (e AaveLiquidationCallEvent) GetDebtToken() common.Address {
	return e.DebtAsset
}

func (e AaveLiquidationCallEvent) GetCollateralToken() common.Address {
	return e.CollateralAsset
}

func (e AaveLiquidationCallEvent) GetLiquidator() common.Address {
	return e.Liquidator
}

type AaveBorrowEvent struct {
	Reserve          common.Address
	User             common.Address
	OnBehalfOf       common.Address
	Amount           *big.Int
	InterestRateMode uint8
	BorrowRate       *big.Int
	ReferralCode     *big.Int
	Transaction
}

func (e AaveBorrowEvent) GetBorrowAmount() *big.Int {
	return e.Amount
}

func (e AaveBorrowEvent) GetBorrower() common.Address {
	return e.OnBehalfOf
}

func (e AaveBorrowEvent) GetBorrowToken() common.Address {
	return e.Reserve
}

func (e AaveBorrowEvent) String() string {
	return fmt.Sprintf("Borrow borrower: %s, borrowAmount: %s, reserve: %s, txn: %s", e.User, e.Amount.String(), e.Reserve, e.TxHash)
}

func (e AaveBorrowEvent) Type() string {
	return "AaveBorrow"
}

type AaveRepay struct {
	Reserve    common.Address
	User       common.Address
	Repayer    common.Address
	Amount     *big.Int
	UseATokens bool
	Transaction
}

func (e AaveRepay) GetRepayAmount() *big.Int {
	return e.Amount
}

func (e AaveRepay) GetBorrower() common.Address {
	return e.User
}

func (e AaveRepay) GetBorrowToken() common.Address {
	return e.Reserve
}

func (e AaveRepay) GetRepayer() common.Address {
	return e.Repayer
}

func (e AaveRepay) String() string {
	return fmt.Sprintf("Repay payer: %s, amount: %s, borrower: %s, txn: %s", e.Repayer.Hex(), e.Amount.String(), e.User.Hex(), e.TxHash)
}

func (e AaveRepay) Type() string {
	return "AaveRepay"
}

func UnpackAaveEvent(abi abi.ABI, protocolInstance *model.ProtocolInstance, eventDefn *model.EventDefn, log types.Log) (PossibleEvent, error) {
	if eventDefn.TopicName == "Borrow" {
		event := AaveBorrowEvent{}
		err := abi.UnpackIntoInterface(&event, eventDefn.TopicName, log.Data)
		event.TxHash = log.TxHash.String()
		event.Reserve = common.HexToAddress(log.Topics[1].Hex())
		event.OnBehalfOf = common.HexToAddress(log.Topics[2].Hex())
		return PossibleEvent{Borrowable: event}, err
	} else if eventDefn.TopicName == "Repay" {
		event := AaveRepay{}
		err := abi.UnpackIntoInterface(&event, eventDefn.TopicName, log.Data)
		event.TxHash = log.TxHash.String()
		event.Reserve = common.HexToAddress(log.Topics[1].Hex())
		event.User = common.HexToAddress(log.Topics[2].Hex())
		event.Repayer = common.HexToAddress(log.Topics[3].Hex())
		return PossibleEvent{Repayable: event}, err
	} else if eventDefn.TopicName == "LiquidationCall" {
		fmt.Printf("Topics: %d\n", len(log.Topics))
		event := AaveLiquidationCallEvent{}
		err := abi.UnpackIntoInterface(&event, eventDefn.TopicName, log.Data)
		event.TxHash = log.TxHash.String()
		return PossibleEvent{Liquidatable: event}, err
	}
	return PossibleEvent{}, errors.New(fmt.Sprintf("%s topic name is not supported", eventDefn.TopicName))
}
