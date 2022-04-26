package index

import (
	"bufio"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type AaveTransferEvent struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}

func (e AaveTransferEvent) String() string {
	return fmt.Sprintf("from: %s, to: %s, amount: %s", e.From, e.To, e.Value)
}

func (e AaveTransferEvent) Type() string {
	return "AaveTransfer"
}

type AaveLiquidationCallEvent struct {
	_collateral                 common.Address
	_reserve                    common.Address
	_user                       common.Address
	_purchaseAmount             *big.Int
	_liquidatedCollateralAmount *big.Int
	_accruedBorrowInterest      *big.Int
	_liquidator                 common.Address
	_receiveAToken              bool
	_timestamp                  *big.Int
}

func (e AaveLiquidationCallEvent) String() string {
	return fmt.Sprintf("liquidator: %s, borrower: %s, tokenCollateral: %s, repayAmount: %s, seizeTokens: %s", e._liquidator, e._user, e._collateral, e._purchaseAmount.String(), e._liquidatedCollateralAmount.String())
}

func (e AaveLiquidationCallEvent) Type() string {
	return "AaveLiquidationCall"
}

type AaveBorrowEvent struct {
	_reserve               common.Address
	_user                  common.Address
	_amount                *big.Int
	_borrowRateMode        *big.Int
	_borrowRate            *big.Int
	_originationFee        *big.Int
	_borrowBalanceIncrease *big.Int
	_referral              *big.Int
	_timestamp             *big.Int
}

func (e AaveBorrowEvent) String() string {
	return fmt.Sprintf("borrower: %s, borrowAmount: %s", e._user, e._amount.String())
}

func (e AaveBorrowEvent) Type() string {
	return "AaveBorrow"
}

type AaveRepay struct {
	_reserve               common.Address
	_user                  common.Address
	_repayer               common.Address
	_amountMinusFees       *big.Int
	_fees                  *big.Int
	_borrowBalanceIncrease *big.Int
	_timestamp             *big.Int
}

func (e AaveRepay) String() string {
	return fmt.Sprintf("payer: %s, repayAmount: %s, borrower: %s", e._repayer.Hex(), e._amountMinusFees.String(), e._user.Hex())
}

func (e AaveRepay) Type() string {
	return "AaveRepay"
}

func unpackAaveEvent(abi abi.ABI, le LoanEvent, data []byte) (Typable, error) {
	if le.TopicName == "Borrow" {
		event := AaveBorrowEvent{}
		err := abi.UnpackIntoInterface(&event, le.TopicName, data)
		return event, err
	} else if le.TopicName == "Repay" {
		event := AaveRepay{}
		err := abi.UnpackIntoInterface(&event, le.TopicName, data)
		return event, err
	} else if le.TopicName == "LiquidationCall" {
		event := AaveLiquidationCallEvent{}
		err := abi.UnpackIntoInterface(&event, le.TopicName, data)
		return event, err
	} else if le.TopicName == "Transfer" {
		event := AaveTransferEvent{}
		err := abi.UnpackIntoInterface(&event, le.TopicName, data)
		return event, err
	}
	return nil, errors.New(fmt.Sprintf("%s topic name is not supported", le.TopicName))
}

func GetAEthDefn(abiPath string) (Protocol, error) {
	borrowEventSignature := []byte("Borrow(address,address,uint256,uint256,uint256,uint256,uint256,uint16,uint256)")
	borrowTopicHash := crypto.Keccak256Hash(borrowEventSignature)

	repayEventSignature := []byte("Repay(address,address,address,uint256,uint256,uint256,uint256)")
	repayTopicHash := crypto.Keccak256Hash(repayEventSignature)

	liquidationCallEventSignature := []byte("LiquidationCall(address,address,address,uint256,uint256,uint256,address,bool,uint256)")
	liquidateCallTopicHash := crypto.Keccak256Hash(liquidationCallEventSignature)

	transferEventSignature := []byte("Transfer(address,address,uint256)")
	transferTopicHash := crypto.Keccak256Hash(transferEventSignature)

	f, err := os.Open(path.Join(abiPath, "aETH.abi"))
	if err != nil {
		return Protocol{}, err
	}
	contractAbi, err := abi.JSON(bufio.NewReader(f))
	if err != nil {
		return Protocol{}, err
	}

	return Protocol{
		"AEth",
		common.HexToAddress("0x3a3A65aAb0dd2A17E3F1947bA16138cd37d08c04"),
		[]LoanEvent{
			{"Borrow", borrowTopicHash, borrowTopicHash.Hex()},
			{"Repay", repayTopicHash, repayTopicHash.Hex()},
			{"LiquidationCall", liquidateCallTopicHash, liquidateCallTopicHash.Hex()},
			{"Transfer", transferTopicHash, transferTopicHash.Hex()},
		},
		[]common.Hash{
			borrowTopicHash, repayTopicHash, liquidateCallTopicHash, transferTopicHash,
		},
		contractAbi,
		unpackAaveEvent,
	}, nil
}
